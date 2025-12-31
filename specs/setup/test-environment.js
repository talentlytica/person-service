import { Client } from 'pg';
import { GenericContainer, Network, Wait } from 'testcontainers';
import { PostgreSqlContainer } from '@testcontainers/postgresql';
import { exec } from 'child_process';
import { promisify } from 'util';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const execAsync = promisify(exec);

/**
 * TestEnvironment - Logic Layer
 *
 * Pure business logic for managing test infrastructure.
 * No direct Docker/HTTP calls here - those are delegated to dependencies.
 */
class TestEnvironment {
  constructor() {
    this.network = null;
    this.postgresContainer = null;
    this.serviceContainer = null;
    this.dbClient = null;
    this.serviceUrl = null;
    this.postgresPort = null;
  }

  /**
   * Initialize the entire test environment
   */
  async initialize() {
    try {
      // Step 1: Create network
      this.network = await new Network().start();
      console.log('✓ Docker network created');

      // Step 2: Start PostgreSQL
      await this._startPostgres();
      console.log('✓ PostgreSQL container started');

      // Step 3: Run migrations
      await this._runMigrations();
      console.log('✓ Database migrations completed');

      // Step 4: Start service
      await this._startService();
      console.log('✓ Service container started');

      console.log(`\n✓ Test environment ready!`);
      console.log(`  Service: ${this.serviceUrl}`);
      console.log(`  Database: postgresql://testuser:testpass@localhost:${this.postgresPort}/testdb`);
    } catch (error) {
      console.error('✗ Failed to initialize test environment:', error.message);
      await this.cleanup();
      throw error;
    }
  }

  /**
   * Start PostgreSQL container (Dependencies Layer interaction)
   */
  async _startPostgres() {
    // PostgreSQL always runs on 5432 inside the container.
    // Docker automatically maps this to a random available port on the host (e.g., 32769, 32770, etc.)
    // This avoids conflicts with your local PostgreSQL on port 5432.
    // For container-to-container communication within the network, the service uses postgres-db:5432
    this.postgresContainer = await new PostgreSqlContainer('postgres:18-alpine')
      .withDatabase('testdb')
      .withUsername('testuser')
      .withPassword('testpass')
      .withExposedPorts(5432)
      .withNetwork(this.network)
      .withNetworkAliases('postgres-db')
      .start();

    this.postgresPort = this.postgresContainer.getMappedPort(5432);

    // Initialize DB client for host communication (localhost with mapped port)
    this.dbClient = new Client({
      connectionString: this.postgresContainer.getConnectionUri()
    });

    await this.dbClient.connect();
  }

  /**
   * Run database migrations from schema.sql using psql command line
   */
  async _runMigrations() {
    const __filename = fileURLToPath(import.meta.url);
    const __dirname = dirname(__filename);
    const schemaPath = resolve(__dirname, '../../source/app/internal/db/schema.sql');

    // Execute schema.sql using psql command line tool
    const psqlCommand = `PGPASSWORD=testpass psql -h localhost -p ${this.postgresPort} -U testuser -d testdb -f "${schemaPath}"`;
    
    try {
      const { stdout, stderr } = await execAsync(psqlCommand);
      if (stdout) console.log('  Migration output:', stdout.trim());
      if (stderr && !stderr.includes('NOTICE')) {
        console.warn('  Migration warnings:', stderr.trim());
      }
    } catch (error) {
      console.error('  Migration failed:', error.message);
      throw error;
    }
  }

  /**
   * Start the Go service container (Dependencies Layer interaction)
   */
  async _startService() {
    const servicePort = 3000;
    
    // Workaround for environments without iptables: use host.docker.internal or gateway IP
    // Get the Docker bridge gateway IP to access host services
    let dbHost = 'host.docker.internal';  // Works on Docker Desktop
    let dbPort = this.postgresPort;        // Use the mapped host port
    
    // For Linux, we need to use the gateway IP instead
    try {
      const { stdout } = await execAsync(
        `docker network inspect ${this.network.getId()} -f '{{range .IPAM.Config}}{{.Gateway}}{{end}}'`
      );
      const gatewayIp = stdout.trim();
      if (gatewayIp && gatewayIp !== '') {
        dbHost = gatewayIp;
        console.log(`Using gateway IP for database: ${gatewayIp}`);
      }
    } catch (error) {
      console.log('Could not detect gateway IP, using host.docker.internal');
    }
    
    const databaseUrl = `postgresql://testuser:testpass@${dbHost}:${dbPort}/testdb?sslmode=disable`;

    console.log('Starting service container...');
    console.log(`Database connection: ${databaseUrl}`);

    let startedContainer = null;
    try {
      const container = new GenericContainer('source-person-service:latest')
        .withEnvironment({ DATABASE_URL: databaseUrl })
        .withEnvironment({ ENCRYPTION_KEY_1: 'test-encryption-key-12345' })
        .withExposedPorts(servicePort)
        .withDefaultLogDriver()
        .withNetwork(this.network)
        .withWaitStrategy(Wait.forLogMessage(/INFO: Server starting on port/))
        .withStartupTimeout(100000); // Increased to 2 minutes for database connection retries

      startedContainer = await container.start();
      this.serviceContainer = startedContainer;

      const mappedPort = this.serviceContainer.getMappedPort(servicePort);
      this.serviceUrl = `http://localhost:${mappedPort}`;
      console.log(`Service started on port ${mappedPort}`);
    } catch (error) {
      console.error('\n✗ Service container failed to start!');
      console.error(`  Error: ${error.message}\n`);

      // Try to get logs before any cleanup
      try {
        // Get the most recent service container
        const { stdout: containerId } = await execAsync(
          `docker ps -a --filter "ancestor=source-person-service:latest" --format "{{.ID}}" | head -1`
        );
        
        if (containerId?.trim()) {
          const { stdout: logs } = await execAsync(`docker logs ${containerId.trim()} 2>&1 || true`);
          if (logs?.trim()) {
            console.error('  Container logs:');
            logs.trim().split('\n').forEach(line => {
              if (line.trim()) console.error('  ' + line);
            });
          }
        }
      } catch (logError) {
        // Could not fetch logs
      }

      // Extract container ID from error message for further troubleshooting
      await this._logFailedContainer(error, databaseUrl);

      throw error;
    }
  }

  /**
   * Extract container ID from error and fetch logs from crashed container
   */
  async _logFailedContainer(error, databaseUrl) {
    // Try to extract container ID from error message
    let containerId = error.message?.match(/container ([a-f0-9]{64})/)?.[1];
    
    // If we have a service container reference but no ID in error, try to get it directly
    if (!containerId && this.serviceContainer) {
      try {
        containerId = this.serviceContainer.getId?.();
      } catch (e) {
        // Ignore
      }
    }

    if (containerId) {
      try {
        // Fetch logs from the dead container - this is the most reliable way to see why it failed
        const { stdout: logs } = await execAsync(`docker logs ${containerId} 2>&1 || true`);
        if (logs?.trim()) {
          console.error('\n  Container logs:');
          logs.trim().split('\n').forEach(line => {
            if (line.trim()) console.error('  ' + line);
          });
        }
      } catch (logError) {
        // Could not fetch logs
      }
    } else {
      // Try to get logs from all running containers with our image
      try {
        const { stdout: containers } = await execAsync(
          `docker ps -a --filter "ancestor=source-person-service:latest" --format "{{.ID}}" | tail -1`
        );
        const lastContainerId = containers?.trim();
        if (lastContainerId) {
          const { stdout: logs } = await execAsync(`docker logs ${lastContainerId} 2>&1 || true`);
          if (logs?.trim()) {
            console.error('\n  Container logs (from last container):');
            logs.trim().split('\n').forEach(line => {
              if (line.trim()) console.error('  ' + line);
            });
          }
        }
      } catch (e) {
        // Ignore
      }
    }

    // Common troubleshooting hints
    console.error('\n  Troubleshooting:');
    console.error('  1. Verify image exists: docker images | grep source-person-service');
    console.error(
      '  2. Test image manually: docker run --rm -e DATABASE_URL="' +
        databaseUrl +
        '" source-person-service:latest'
    );
    console.error('  3. Check service logs in troubleshooting steps');
  }

  /**
   * Truncate all tables between tests
   */
  async cleanupDatabase() {
    if (!this.dbClient) return;
  }

  /**
   * Clean up all containers and network
   */
  async cleanup() {
    // Close DB client first
    if (this.dbClient) {
      try {
        await this.dbClient.end();
      } catch (error) {
        // Ignore - connection may already be closed
      }
      this.dbClient = null;
    }

    // Stop service container (may already be stopped if it crashed)
    if (this.serviceContainer) {
      try {
        await this.serviceContainer.stop();
      } catch (error) {
        // Ignore 409 errors - container already stopped/paused
        if (!error.statusCode || error.statusCode !== 409) {
          console.warn('Warning: Failed to stop service container:', error.message);
        }
      }
      this.serviceContainer = null;
    }

    // Stop postgres container
    if (this.postgresContainer) {
      try {
        await this.postgresContainer.stop();
      } catch (error) {
        if (!error.statusCode || error.statusCode !== 409) {
          console.warn('Warning: Failed to stop postgres container:', error.message);
        }
      }
      this.postgresContainer = null;
    }

    // Stop network
    if (this.network) {
      try {
        await this.network.stop();
      } catch (error) {
        // Ignore network cleanup errors
      }
      this.network = null;
    }

    console.log('✓ Test environment cleaned up');
  }

  /**
   * Get database client for direct queries
   */
  getDbClient() {
    return this.dbClient;
  }

  /**
   * Get service base URL
   */
  getServiceUrl() {
    return this.serviceUrl;
  }
}

export default TestEnvironment;
