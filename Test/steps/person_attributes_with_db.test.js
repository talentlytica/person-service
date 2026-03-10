/**
 * Person Attributes with Database Verification
 * 
 * This test:
 * 1. Creates a test person directly in database
 * 2. Creates an attribute via API
 * 3. Verifies the attribute exists in database (ENCRYPTED)
 * 4. Gets attribute via API and verifies decryption
 * 5. Updates the attribute via API
 * 6. Verifies update in database
 * 7. Deletes the attribute
 * 8. Verifies deletion in database
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import axios from 'axios';
import pg from 'pg';
const { Client } = pg;
import { config } from 'dotenv';

config();

describe('Person Attributes - API + Database Verification', () => {
  let apiClient;
  let dbClient;
  let testPersonId;
  let createdAttributeId;
  const encryptionKey = process.env.ENCRYPTION_KEY_1 || 'default-key-for-dev';
  
  beforeAll(async () => {
    console.log('\n🚀 Starting Person Attributes Test with DB Verification...\n');
    
    // Setup API client
    apiClient = axios.create({
      baseURL: process.env.BASE_URL,
      timeout: 10000,
      headers: {
        'x-api-key': process.env.AUTH_TOKEN,
        'Content-Type': 'application/json'
      }
    });
    
    // Setup DB client
    dbClient = new Client({
      host: process.env.DB_HOST,
      port: process.env.DB_PORT,
      database: process.env.DB_NAME,
      user: process.env.DB_USER,
      password: process.env.DB_PASSWORD
    });
    
    await dbClient.connect();
    console.log('✅ Database connected\n');
    
    // Create a test person in database
    console.log('🔧 Setup: Creating test person in database...');
    const result = await dbClient.query(`
      INSERT INTO person (id, client_id, created_at, updated_at)
      VALUES (gen_random_uuid(), $1, NOW(), NOW())
      RETURNING id
    `, [`test-client-${Date.now()}`]);
    
    testPersonId = result.rows[0].id;
    console.log(`✅ Test person created with ID: ${testPersonId}\n`);
  });
  
  afterAll(async () => {
    // Cleanup: Delete test data
    if (testPersonId) {
      try {
        await dbClient.query('DELETE FROM person_attributes WHERE person_id = $1', [testPersonId]);
        await dbClient.query('DELETE FROM person WHERE id = $1', [testPersonId]);
        console.log('\n🧹 Cleanup: Test data deleted from database');
      } catch (error) {
        console.log('⚠️  Cleanup warning:', error.message);
      }
    }
    
    await dbClient.end();
    console.log('✅ Database disconnected\n');
  });
  
  test('1. CREATE attribute via API and verify in database (ENCRYPTED)', async () => {
    console.log('📝 Test 1: CREATE Attribute\n');
    
    // Step 1: Create attribute via API
    console.log('🔵 Step 1: Sending POST request to create attribute...');
    const createData = {
      key: 'email',
      value: 'test@example.com',
      meta: {
        caller: 'test-suite',
        reason: 'automated-testing',
        traceId: `test-${Date.now()}`
      }
    };
    
    let apiResponse;
    try {
      apiResponse = await apiClient.post(`/persons/${testPersonId}/attributes`, createData);
    } catch (error) {
      if (error.response) {
        console.log('❌ API Error:', error.response.status, error.response.data);
      }
      throw error;
    }
    
    console.log('📊 API Response Status:', apiResponse.status);
    console.log('📊 API Response Data:', JSON.stringify(apiResponse.data, null, 2));
    
    // Step 2: Verify API response
    expect(apiResponse.status).toBe(201);
    expect(apiResponse.data).toHaveProperty('id');
    expect(apiResponse.data.key).toBe('email');
    expect(apiResponse.data.value).toBe('test@example.com'); // API should return decrypted
    
    createdAttributeId = apiResponse.data.id;
    console.log(`✅ Attribute created via API with ID: ${createdAttributeId}\n`);
    
    // Step 3: Query database directly
    console.log('🔍 Step 2: Querying database to verify attribute exists...');
    const dbResult = await dbClient.query(
      'SELECT * FROM person_attributes WHERE id = $1',
      [createdAttributeId]
    );
    
    console.log(`📊 Database query returned ${dbResult.rows.length} row(s)`);
    
    // Step 4: Verify data in database
    expect(dbResult.rows.length).toBe(1);
    const attributeInDb = dbResult.rows[0];
    
    console.log('📋 Attribute in Database:');
    console.log(`   ID: ${attributeInDb.id}`);
    console.log(`   Person ID: ${attributeInDb.person_id}`);
    console.log(`   Attribute Key: ${attributeInDb.attribute_key}`);
    const encryptedHex = Buffer.isBuffer(attributeInDb.encrypted_value)
      ? attributeInDb.encrypted_value.toString('hex').substring(0, 50)
      : String(attributeInDb.encrypted_value).substring(0, 50);
    console.log(`   Encrypted Value: ${encryptedHex}...`);
    console.log(`   Created At: ${attributeInDb.created_at}`);
    
    // Step 5: Verify encryption - encrypted_value is a Buffer (BYTEA)
    const encryptedStr = Buffer.isBuffer(attributeInDb.encrypted_value)
      ? attributeInDb.encrypted_value.toString('utf8')
      : String(attributeInDb.encrypted_value);
    expect(encryptedStr).not.toBe('test@example.com'); // Should be encrypted
    expect(Buffer.isBuffer(attributeInDb.encrypted_value)).toBe(true); // pgcrypto returns BYTEA
    
    // Step 6: Decrypt and verify
    console.log('\n🔓 Step 3: Decrypting value from database...');
    const decryptResult = await dbClient.query(`
      SELECT pgp_sym_decrypt(encrypted_value::bytea, $1) as decrypted_value
      FROM person_attributes
      WHERE id = $2
    `, [encryptionKey, createdAttributeId]);
    
    const decryptedValue = decryptResult.rows[0].decrypted_value.toString();
    console.log(`   Decrypted Value: ${decryptedValue}`);
    
    expect(decryptedValue).toBe('test@example.com');
    
    console.log('\n✅ Test 1 PASSED: Attribute created, encrypted in DB, and verified!\n');
  });
  
  test('2. GET attribute via API and verify decryption', async () => {
    console.log('📝 Test 2: GET Attribute\n');
    
    expect(createdAttributeId).toBeDefined();
    
    // Step 1: Get attribute via API
    console.log('🔵 Step 1: Sending GET request to retrieve attribute...');
    
    let apiResponse;
    try {
      apiResponse = await apiClient.get(`/persons/${testPersonId}/attributes/${createdAttributeId}`);
    } catch (error) {
      if (error.response) {
        console.log('❌ API Error:', error.response.status, error.response.data);
      }
      throw error;
    }
    
    console.log('📊 API Response Status:', apiResponse.status);
    console.log('📊 API Response Data:', JSON.stringify(apiResponse.data, null, 2));
    
    // Step 2: Verify API response (should be decrypted)
    expect(apiResponse.status).toBe(200);
    expect(apiResponse.data.key).toBe('email');
    expect(apiResponse.data.value).toBe('test@example.com'); // Decrypted by API
    
    console.log('✅ API correctly decrypted the value!');
    
    // Step 3: Verify in database it's still encrypted
    console.log('\n🔍 Step 2: Verifying database still has encrypted value...');
    const dbResult = await dbClient.query(
      'SELECT encrypted_value FROM person_attributes WHERE id = $1',
      [createdAttributeId]
    );
    
    const encryptedInDb = dbResult.rows[0].encrypted_value;
    expect(encryptedInDb).not.toBe('test@example.com'); // Still encrypted in DB
    
    console.log('✅ Value still encrypted in database (correct!)');
    console.log('\n✅ Test 2 PASSED: API decryption working correctly!\n');
  });
  
  test('3. UPDATE attribute via API and verify in database', async () => {
    console.log('📝 Test 3: UPDATE Attribute\n');
    
    expect(createdAttributeId).toBeDefined();
    
    // Step 1: Update attribute via API
    console.log('🔵 Step 1: Sending PUT request to update attribute...');
    const updateData = {
      key: 'email',
      value: 'updated@example.com',
      meta: {
        caller: 'test-suite',
        reason: 'automated-testing-update',
        traceId: `test-update-${Date.now()}`
      }
    };
    
    let apiResponse;
    try {
      apiResponse = await apiClient.put(`/persons/${testPersonId}/attributes/${createdAttributeId}`, updateData);
    } catch (error) {
      if (error.response) {
        console.log('❌ API Error:', error.response.status, error.response.data);
      }
      throw error;
    }
    
    console.log('📊 API Response Status:', apiResponse.status);
    console.log('📊 API Response Data:', JSON.stringify(apiResponse.data, null, 2));
    
    // Step 2: Verify API response
    expect(apiResponse.status).toBe(200);
    expect(apiResponse.data.value).toBe('updated@example.com');
    
    console.log('✅ Attribute updated via API\n');
    
    // Step 3: Decrypt and verify in database
    console.log('🔍 Step 2: Querying database to verify update...');
    const decryptResult = await dbClient.query(`
      SELECT pgp_sym_decrypt(encrypted_value::bytea, $1) as decrypted_value
      FROM person_attributes
      WHERE id = $2
    `, [encryptionKey, createdAttributeId]);
    
    const decryptedValue = decryptResult.rows[0].decrypted_value.toString();
    console.log(`   Decrypted Value in DB: ${decryptedValue}`);
    
    expect(decryptedValue).toBe('updated@example.com');
    
    console.log('\n✅ Test 3 PASSED: Attribute updated and verified in database!\n');
  });
  
  test('4. DELETE attribute via API and verify in database', async () => {
    console.log('📝 Test 4: DELETE Attribute\n');
    
    expect(createdAttributeId).toBeDefined();
    
    // Step 1: Delete attribute via API
    console.log('🔵 Step 1: Sending DELETE request to delete attribute...');
    
    let apiResponse;
    try {
      apiResponse = await apiClient.delete(`/persons/${testPersonId}/attributes/${createdAttributeId}`);
    } catch (error) {
      if (error.response) {
        console.log('❌ API Error:', error.response.status, error.response.data);
      }
      throw error;
    }
    
    console.log('📊 API Response Status:', apiResponse.status);
    
    // Step 2: Verify API response
    expect([200, 204]).toContain(apiResponse.status);
    
    console.log('✅ Attribute deleted via API\n');
    
    // Step 3: Verify deleted in database
    console.log('🔍 Step 2: Querying database to verify deletion...');
    const dbResult = await dbClient.query(
      'SELECT * FROM person_attributes WHERE id = $1',
      [createdAttributeId]
    );
    
    console.log(`📊 Database query returned ${dbResult.rows.length} row(s)`);
    
    // Verify attribute no longer exists (hard delete)
    expect(dbResult.rows.length).toBe(0);
    
    console.log('✅ Attribute successfully deleted from database');
    console.log('\n✅ Test 4 PASSED: Attribute deleted and verified!\n');
  });
  
  test('5. GET ALL attributes for person', async () => {
    console.log('📝 Test 5: GET ALL Attributes\n');
    
    // Create multiple attributes first
    console.log('🔧 Setup: Creating multiple attributes...');
    const metaData = {
      caller: 'test-suite',
      reason: 'automated-testing-multiple',
      traceId: `test-multi-${Date.now()}`
    };
    const attributes = [
      { key: 'phone', value: '+1234567890', meta: metaData },
      { key: 'address', value: '123 Main St', meta: metaData },
      { key: 'city', value: 'Jakarta', meta: metaData }
    ];
    
    for (const attr of attributes) {
      await apiClient.post(`/persons/${testPersonId}/attributes`, attr);
    }
    console.log(`✅ Created ${attributes.length} attributes\n`);
    
    // Get all attributes
    console.log('🔵 Step 1: Sending GET request to retrieve all attributes...');
    
    let apiResponse;
    try {
      apiResponse = await apiClient.get(`/persons/${testPersonId}/attributes`);
    } catch (error) {
      if (error.response) {
        console.log('❌ API Error:', error.response.status, error.response.data);
      }
      throw error;
    }
    
    console.log('📊 API Response Status:', apiResponse.status);
    console.log('📊 API Response Data:', JSON.stringify(apiResponse.data, null, 2));
    
    // Verify API response
    expect(apiResponse.status).toBe(200);
    expect(apiResponse.data).toBeInstanceOf(Array);
    expect(apiResponse.data.length).toBeGreaterThanOrEqual(3); // At least 3 attributes
    
    console.log(`✅ Retrieved ${apiResponse.data.length} attributes`);
    
    // Verify all values are decrypted
    apiResponse.data.forEach(attr => {
      console.log(`   - ${attr.key}: ${attr.value}`);
      expect(attr.value).not.toContain('\\x'); // Should be decrypted
    });
    
    // Verify in database all are encrypted
    console.log('\n🔍 Step 2: Verifying all are encrypted in database...');
    const dbResult = await dbClient.query(
      'SELECT attribute_key, encrypted_value FROM person_attributes WHERE person_id = $1',
      [testPersonId]
    );
    
    dbResult.rows.forEach(row => {
      expect(Buffer.isBuffer(row.encrypted_value)).toBe(true); // Should be encrypted BYTEA
    });
    
    console.log(`✅ All ${dbResult.rows.length} attributes are encrypted in database`);
    console.log('\n✅ Test 5 PASSED: Get all attributes working correctly!\n');
  });
});
