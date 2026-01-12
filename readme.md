# Installation Guide

This document explains how to set up all dependencies needed to run the person-service project.

## Prerequisites

The project requires the following tools to be installed:

- **Docker & Docker Compose** - For containerized deployment
- **Go 1.24+** - For building the service
- **Node.js & npm** - For the test suite (specs folder)
- **Python3** - For serving coverage reports

## Installation by OS

### macOS

1. **Install Docker Desktop** (includes Docker & Docker Compose)
   - Download from: https://docs.docker.com/get-docker/
   - Or via Homebrew: `brew install docker`

2. **Install Go**
   - Download from: https://golang.org/dl/
   - Or via Homebrew: `brew install go`

3. **Install Node.js & npm**
   - Download from: https://nodejs.org/
   - Or via Homebrew: `brew install node`

4. **Python3**
   - Usually pre-installed on macOS
   - If needed: `brew install python3`

5. **Run the install target**
   ```bash
   make install
   ```
   This will verify all dependencies and download project dependencies.

### Linux

1. **Install Docker Engine, Docker Compose, and dependencies**
   ```bash
   sudo sh get-docker.sh
   ```
   Or follow the [official Docker documentation](https://docs.docker.com/engine/install/) for your distribution.

2. **Install Go**
   ```bash
   # For Ubuntu/Debian
   sudo apt-get update
   sudo apt-get install golang-go
   
   # Or download from https://golang.org/dl/
   ```

3. **Install Node.js & npm**
   ```bash
   # For Ubuntu/Debian
   sudo apt-get install nodejs npm
   
   # Or download from https://nodejs.org/
   ```

4. **Install Python3**
   ```bash
   # Usually pre-installed, but if needed:
   sudo apt-get install python3
   ```

5. **Run the install target**
   ```bash
   make install
   ```
   This will verify all dependencies and download project dependencies.

### Windows

1. **Install Docker Desktop for Windows**
   - Download from: https://docs.docker.com/get-docker/
   - Use WSL 2 (Windows Subsystem for Linux 2) backend for best compatibility

2. **Install Go**
   - Download from: https://golang.org/dl/

3. **Install Node.js & npm**
   - Download from: https://nodejs.org/

4. **Install Python3**
   - Download from: https://www.python.org/downloads/
   - Or use Windows Package Manager: `winget install Python.Python.3.11`

5. **Open WSL terminal and run**
   ```bash
   make install
   ```


## Sanity Check
To check if all's good

```bash
make test
```

## Quick Start

Once all dependencies are installed:

```bash
# Install project dependencies (Go modules, npm packages)
make install

# Build the Docker image and prepare specs
make build

# Start the service
cd source && docker compose up

```

## Verify Installation

To verify all dependencies are correctly installed, run:

```bash
make install
```

This will check for all required tools and report which ones are missing with links to installation guides.

## Troubleshooting

### Docker not found
- Ensure Docker Desktop is installed and running
- Restart your terminal/IDE after installing Docker

### Go modules issues
```bash
cd source/app
go mod download
go mod tidy
```

### npm install issues
```bash
cd specs
npm install
# Or clean and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Permission denied errors on Linux
- Add your user to the docker group: `sudo usermod -aG docker $USER`
- Log out and back in, or run: `newgrp docker`

## Environment Configuration

The service uses a `.env` file in the `source/` directory for configuration. Make sure it exists before running `docker compose up`:

```
DATABASE_URL=postgres://<postgre username>:<postgre password>@host.docker.internal:5432/person_service?sslmode=disable
PORT=3000
PERSON_API_KEY_BLUE=person-service-key-<uuid>
PERSON_API_KEY_GREEN=person-service-key-<uuid>
```

You need to add .env manually and set with proper value

## Support

For issues related to:
- **Docker**: https://docs.docker.com/support/
- **Go**: https://golang.org/help
- **Node.js/npm**: https://nodejs.org/en/docs/
- **This project**: Check the source/README.md for project-specific documentation
