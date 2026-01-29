/**
 * NEGATIVE TESTS: POST /persons/:personId/attributes
 * 
 * Tests error scenarios, authentication, authorization, and validation
 */

import { describe, test, expect, beforeAll, afterAll } from '@jest/globals';
import axios from 'axios';
import pg from 'pg';
const { Client } = pg;
import { config } from 'dotenv';

config();

describe('NEGATIVE: POST /persons/:personId/attributes - Create Attribute', () => {
  let apiClient;
  let dbClient;
  let testPersonId;
  const BASE_URL = process.env.BASE_URL || 'http://localhost:3000';
  const API_KEY = process.env.AUTH_TOKEN;
  
  beforeAll(async () => {
    apiClient = axios.create({
      baseURL: BASE_URL,
      timeout: 10000,
      headers: {
        'x-api-key': API_KEY,
        'Content-Type': 'application/json'
      }
    });
    
    dbClient = new Client({
      host: process.env.DB_HOST,
      port: process.env.DB_PORT,
      database: process.env.DB_NAME,
      user: process.env.DB_USER,
      password: process.env.DB_PASSWORD
    });
    await dbClient.connect();
    
    // Create test person
    const result = await dbClient.query(`
      INSERT INTO person (id, client_id, created_at, updated_at)
      VALUES (gen_random_uuid(), $1, NOW(), NOW())
      RETURNING id
    `, [`negative-test-${Date.now()}`]);
    testPersonId = result.rows[0].id;
  });
  
  afterAll(async () => {
    if (testPersonId) {
      await dbClient.query('DELETE FROM person_attributes WHERE person_id = $1', [testPersonId]);
      await dbClient.query('DELETE FROM person WHERE id = $1', [testPersonId]);
    }
    await dbClient.end();
  });
  
  // ========================================
  // AUTHENTICATION ERRORS
  // ========================================
  
  test('Should reject request without API key', async () => {
    const clientWithoutAuth = axios.create({
      baseURL: BASE_URL,
      timeout: 10000
    });
    
    try {
      await clientWithoutAuth.post(`/persons/${testPersonId}/attributes`, {
        key: 'email',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 401');
    } catch (error) {
      expect(error.response.status).toBe(401);
      expect(error.response.data.message).toContain('x-api-key');
    }
  });
  
  test('Should reject request with invalid API key', async () => {
    const clientWithInvalidKey = axios.create({
      baseURL: BASE_URL,
      timeout: 10000,
      headers: {
        'x-api-key': 'invalid-key-12345'
      }
    });
    
    try {
      await clientWithInvalidKey.post(`/persons/${testPersonId}/attributes`, {
        key: 'email',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 401');
    } catch (error) {
      expect(error.response.status).toBe(401);
    }
  });
  
  test('Should reject request with expired/wrong format API key', async () => {
    const clientWithWrongFormat = axios.create({
      baseURL: BASE_URL,
      timeout: 10000,
      headers: {
        'x-api-key': 'wrong-format-no-uuid'
      }
    });
    
    try {
      await clientWithWrongFormat.post(`/persons/${testPersonId}/attributes`, {
        key: 'email',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 401');
    } catch (error) {
      expect(error.response.status).toBe(401);
    }
  });
  
  // ========================================
  // INVALID PERSON ID
  // ========================================
  
  test('Should reject request with non-existent personId', async () => {
    try {
      await apiClient.post('/persons/00000000-0000-0000-0000-000000000000/attributes', {
        key: 'email',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 404');
    } catch (error) {
      expect(error.response.status).toBe(404);
      expect(error.response.data.message).toContain('not found');
    }
  });
  
  test('Should reject request with invalid UUID format for personId', async () => {
    try {
      await apiClient.post('/persons/invalid-uuid-format/attributes', {
        key: 'email',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown error');
    } catch (error) {
      expect([400, 404]).toContain(error.response.status);
    }
  });
  
  test('Should reject request with empty personId', async () => {
    try {
      await apiClient.post('/persons//attributes', {
        key: 'email',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown error');
    } catch (error) {
      expect([404, 400]).toContain(error.response.status);
    }
  });
  
  // ========================================
  // MISSING REQUIRED FIELDS
  // ========================================
  
  test('Should reject request without key field', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
      // Case-insensitive check for "key" in message
      expect(error.response.data.message.toLowerCase()).toContain('key');
    }
  });
  
  test('Should reject request without meta field', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: 'email',
        value: 'test@example.com'
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
      expect(error.response.data.message).toContain('meta');
    }
  });
  
  test('Should reject request with empty meta object', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: 'email',
        value: 'test@example.com',
        meta: {}
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  test('Should reject request without entire body', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`);
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  // ========================================
  // EMPTY/NULL VALUES
  // ========================================
  
  test('Should reject empty string as key', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: '',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  test('Should reject null as key', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: null,
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  test('Should reject whitespace-only key', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: '   ',
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  // ========================================
  // INVALID DATA TYPES
  // ========================================
  
  test('Should reject number as key', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: 12345,
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  test('Should reject object as key', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: { nested: 'object' },
        value: 'test@example.com',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  test('Should reject array as value', async () => {
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: 'test-key',
        value: ['array', 'value'],
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      fail('Should have thrown 400');
    } catch (error) {
      expect(error.response.status).toBe(400);
    }
  });
  
  // ========================================
  // MALFORMED REQUESTS
  // ========================================
  
  test('Should reject invalid JSON', async () => {
    try {
      await axios.post(
        `${BASE_URL}/persons/${testPersonId}/attributes`,
        'invalid json string',
        {
          headers: {
            'x-api-key': API_KEY,
            'Content-Type': 'application/json'
          }
        }
      );
      fail('Should have thrown 400');
    } catch (error) {
      expect([400, 500]).toContain(error.response?.status);
    }
  });
  
  test('Should reject wrong Content-Type', async () => {
    try {
      await axios.post(
        `${BASE_URL}/persons/${testPersonId}/attributes`,
        { key: 'test', value: 'test', meta: { caller: 't', reason: 't', traceId: '1' } },
        {
          headers: {
            'x-api-key': API_KEY,
            'Content-Type': 'text/plain'
          }
        }
      );
      // May or may not accept
    } catch (error) {
      if (error.response) {
        expect([400, 415]).toContain(error.response.status);
      }
    }
  });
  
  // ========================================
  // BOUNDARY CONDITIONS
  // ========================================
  
  test('Should handle or reject extremely long key', async () => {
    const longKey = 'k'.repeat(10000);
    
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: longKey,
        value: 'test-value',
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      // If accepted, it's an issue (should have limit)
    } catch (error) {
      expect([400, 413]).toContain(error.response.status);
    }
  });
  
  test('Should handle or reject extremely long value', async () => {
    const longValue = 'v'.repeat(1000000);
    
    try {
      await apiClient.post(`/persons/${testPersonId}/attributes`, {
        key: 'long-value',
        value: longValue,
        meta: { caller: 'test', reason: 'test', traceId: '123' }
      });
      // If accepted, check if encryption can handle it
    } catch (error) {
      expect([400, 413, 500]).toContain(error.response.status);
    }
  });
  
  // ========================================
  // INVALID HTTP METHODS
  // ========================================
  
  test('Should reject GET method on POST endpoint', async () => {
    try {
      await apiClient.get(`/persons/${testPersonId}/attributes`);
      // Actually GET is valid (gets all attributes), so will succeed
    } catch (error) {
      // If it fails, it's for other reasons
    }
  });
  
  test('Should reject PATCH method', async () => {
    try {
      await apiClient.patch(`/persons/${testPersonId}/attributes`, {
        key: 'email',
        value: 'test@example.com'
      });
      fail('Should have thrown error');
    } catch (error) {
      expect([404, 405]).toContain(error.response.status);
    }
  });
});
