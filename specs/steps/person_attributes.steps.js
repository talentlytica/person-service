/**
 * Person Attributes Step Definitions
 * 
 * Step definitions for the Person Attributes feature tests.
 * Uses jest-cucumber to bind Gherkin scenarios to Jest tests.
 * 
 * IMPORTANT: This test file uses PostgreSQL's pgcrypto extension
 * for encrypting/decrypting attribute values stored in the database.
 * Encryption key is configurable via ENCRYPTION_KEY environment variable.
 */

import { loadFeature, defineFeature } from 'jest-cucumber';
import { expect, beforeEach } from '@jest/globals';
import {
  createTestContext,
  assertServiceRunning,
  sendPost,
  sendGet,
  assertResponseStatus
} from './common.steps.js';
import {
  parseJsonResponse,
  sendPutRequest,
  sendDeleteRequest
} from '../helpers/api.js';

const feature = loadFeature('../features/person_attributes.feature', {
  loadRelativePath: true
});

/**
 * Parse table data into an object
 * @param {Array} table - Cucumber table data
 * @returns {Object} Parsed object with key-value pairs
 */
function parseTableToObject(table) {
  // jest-cucumber passes table data as an array of objects
  // If table is already an object array, return the first object
  if (table && table.length > 0 && typeof table[0] === 'object' && !Array.isArray(table[0])) {
    return table[0];
  }

  // Fallback to original format (2D array)
  const headers = table[0];
  const values = table[1];
  const result = {};
  headers.forEach((header, index) => {
    result[header] = values[index];
  });
  return result;
}

/**
 * Parse table data into an array of objects
 * @param {Array} table - Cucumber table data
 * @returns {Array} Array of objects
 */
function parseTableToArray(table) {
  // jest-cucumber passes table data as an array of objects
  // If table is already an object array, return it as-is
  if (table && table.length > 0 && typeof table[0] === 'object' && !Array.isArray(table[0])) {
    return table;
  }

  // Fallback to original format (2D array)
  const headers = table[0];
  const result = [];
  for (let i = 1; i < table.length; i++) {
    const row = {};
    headers.forEach((header, index) => {
      row[header] = table[i][index];
    });
    result.push(row);
  }
  return result;
}

defineFeature(feature, (test) => {
  let ctx;
  let dbClient;
  const ENCRYPTION_KEY = 'test-encryption-key-12345';
  const KEY_VERSION = 1; // Integer representing encryption key version

  beforeEach(async () => {
    ctx = createTestContext();
    // Store test-specific data
    ctx.personId = null;
    ctx.attributeId = null;
    ctx.attributes = [];
    ctx.createdPerson = null;
    ctx.createdAttribute = null;
    ctx.meta = null;

    // Get database client from global test environment
    if (global.__TEST_ENV__) {
      dbClient = global.__TEST_ENV__.getDbClient();
    }
  });

  // Background step - runs before each scenario
  const setupBackground = ({ given, and }) => {
    given('the persons and attributes table is empty', async () => {
      if (!dbClient) {
        throw new Error('Database client not available');
      }
      // Clear attributes first (due to foreign key)
      await dbClient.query('DELETE FROM person_attributes');
      // Clear persons
      await dbClient.query('DELETE FROM person');
    });

    and('the service is running', () => {
      assertServiceRunning(expect);
    });
  };

  // Helper: Create a person in the database
  async function createPerson(personData) {
    const insertQuery = `
      INSERT INTO person (client_id, created_at, updated_at)
      VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
      RETURNING id, client_id, created_at, updated_at
    `;
    const result = await dbClient.query(insertQuery, [personData.clientId]);
    return result.rows[0];
  }

  // Helper: Create an attribute for a person (with encryption)
  async function createAttribute(personId, attributeData, meta = null) {
    const insertQuery = `
      INSERT INTO person_attributes (
        person_id, 
        attribute_key, 
        encrypted_value, 
        key_version,
        created_at, 
        updated_at
      )
      VALUES (
        $1, 
        $2, 
        pgp_sym_encrypt($3, $4), 
        $5,
        CURRENT_TIMESTAMP, 
        CURRENT_TIMESTAMP
      )
      RETURNING 
        id, 
        person_id, 
        attribute_key, 
        pgp_sym_decrypt(encrypted_value, $4) AS attribute_value,
        key_version,
        created_at, 
        updated_at
    `;
    const result = await dbClient.query(insertQuery, [
      personId,
      attributeData.key,
      attributeData.value || '',
      ENCRYPTION_KEY,
      KEY_VERSION
    ]);
    return result.rows[0];
  }

  // Helper: Get all attributes for a person (with decryption)
  async function getPersonAttributes(personId) {
    const query = `
      SELECT 
        id,
        person_id,
        attribute_key,
        pgp_sym_decrypt(encrypted_value, $2) AS attribute_value,
        key_version,
        created_at,
        updated_at
      FROM person_attributes 
      WHERE person_id = $1 
      ORDER BY attribute_key
    `;
    const result = await dbClient.query(query, [personId, ENCRYPTION_KEY]);
    return result.rows;
  }

  // Helper: Get raw (encrypted) attribute data from database
  async function getRawAttributeFromDb(personId, attributeKey) {
    const query = `
      SELECT 
        id,
        person_id,
        attribute_key,
        encrypted_value,
        key_version,
        created_at,
        updated_at
      FROM person_attributes 
      WHERE person_id = $1 AND attribute_key = $2
      LIMIT 1
    `;
    const result = await dbClient.query(query, [personId, attributeKey]);
    return result.rows[0];
  }

  // Scenario: Add a single attribute to a person
  test('Add a single attribute to a person', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.meta = metaData;
      ctx.requestBody.meta = metaData;

      // Send the request
      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
        if (ctx.responseData && ctx.responseData.id) {
          ctx.attributeId = ctx.responseData.id;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      expect(ctx.response).toBeDefined();
      expect(ctx.response.status).toBe(parseInt(statusCode));
    });

    and(/^the response should contain an attribute with:$/, async (table) => {
      const expectedData = parseTableToObject(table);
      expect(ctx.responseData).toBeDefined();
      expect(ctx.responseData.key).toBe(expectedData.key);
      expect(ctx.responseData.value).toBe(expectedData.value);
    });

    and(/^the attribute should have an "id"$/, () => {
      expect(ctx.responseData).toHaveProperty('id');
      expect(ctx.responseData.id).toBeDefined();
    });

    and(/^the attribute should have "createdAt" timestamp$/, () => {
      expect(ctx.responseData).toHaveProperty('createdAt');
      const createdAt = new Date(ctx.responseData.createdAt);
      expect(createdAt.toString()).not.toBe('Invalid Date');
    });

    and(/^the attribute should have "updatedAt" timestamp$/, () => {
      expect(ctx.responseData).toHaveProperty('updatedAt');
      const updatedAt = new Date(ctx.responseData.updatedAt);
      expect(updatedAt.toString()).not.toBe('Invalid Date');
    });
  });

  // Scenario: Add multiple attributes to a person
  test('Add multiple attributes to a person', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I add the following attributes to the person:$/, async (table) => {
      const attributesData = parseTableToArray(table);
      ctx.addedAttributes = [];

      const endpoint = `/persons/${ctx.personId}/attributes`;
      const meta = {
        caller: 'test',
        reason: 'add multiple attributes',
        traceId: 'test-trace-multiple'
      };

      for (const attrData of attributesData) {
        const requestBody = {
          key: attrData.key,
          value: attrData.value,
          meta: meta
        };

        const response = await sendPutRequest(endpoint, requestBody);
        ctx.addedAttributes.push(response);
      }
    });

    then('all attributes should be added successfully', () => {
      ctx.addedAttributes.forEach(response => {
        expect(response.status).toBe(201);
      });
    });

    and(/^the person should have (\d+) attributes?$/, async (count) => {
      const attributes = await getPersonAttributes(ctx.personId);
      expect(attributes.length).toBe(parseInt(count));
    });
  });

  // Scenario: Get all attributes for a person
  test('Get all attributes for a person', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and(/^the person has the following attributes:$/, async (table) => {
      const attributesData = parseTableToArray(table);
      ctx.createdAttributes = [];

      for (const attrData of attributesData) {
        const attribute = await createAttribute(ctx.personId, attrData);
        ctx.createdAttributes.push(attribute);
      }
    });

    when(/^I send a GET request to "\/persons\/\{personId\}\/attributes"$/, async () => {
      const endpoint = `/persons/${ctx.personId}/attributes`;
      await sendGet(ctx, endpoint);

      if (ctx.response && ctx.response.ok) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the response should contain (\d+) attributes?$/, (count) => {
      expect(ctx.responseData).toBeDefined();
      expect(Array.isArray(ctx.responseData)).toBe(true);
      expect(ctx.responseData.length).toBe(parseInt(count));
    });

    and(/^the attributes should include:$/, (table) => {
      const expectedAttributes = parseTableToArray(table);

      expectedAttributes.forEach(expected => {
        const found = ctx.responseData.find(attr =>
          attr.key === expected.key && attr.value === expected.value
        );
        expect(found).toBeDefined();
      });
    });
  });

  // Scenario: Update an existing attribute
  test('Update an existing attribute', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and(/^the person has an attribute:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.createdAttribute = await createAttribute(ctx.personId, attributeData);
      ctx.attributeId = ctx.createdAttribute.id;
      ctx.originalUpdatedAt = ctx.createdAttribute.updated_at;
    });

    when(/^I send a PUT request to "\/persons\/\{personId\}\/attributes\/\{attributeId\}" with:$/, async (table) => {
      const updateData = parseTableToObject(table);
      ctx.updateBody = {
        key: updateData.key,
        value: updateData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.updateBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes/${ctx.attributeId}`;
      ctx.response = await sendPutRequest(endpoint, ctx.updateBody);

      if (ctx.response && ctx.response.ok) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the attribute "([^"]*)" should have value "([^"]*)"$/, (key, value) => {
      expect(ctx.responseData).toBeDefined();
      expect(ctx.responseData.key).toBe(key);
      expect(ctx.responseData.value).toBe(value);
    });

    and(/^the attribute "updatedAt" timestamp should be updated$/, () => {
      expect(ctx.responseData.updatedAt).toBeDefined();
      const newUpdatedAt = new Date(ctx.responseData.updatedAt);
      const oldUpdatedAt = new Date(ctx.originalUpdatedAt);
      expect(newUpdatedAt.getTime()).toBeGreaterThanOrEqual(oldUpdatedAt.getTime());
    });
  });

  // Scenario: Delete an attribute from a person
  test('Delete an attribute from a person', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and(/^the person has the following attributes:$/, async (table) => {
      const attributesData = parseTableToArray(table);
      ctx.createdAttributes = [];

      for (const attrData of attributesData) {
        const attribute = await createAttribute(ctx.personId, attrData);
        ctx.createdAttributes.push(attribute);
      }
    });

    when(/^I send a DELETE request to "\/persons\/\{personId\}\/attributes\/\{attributeId\}" for attribute "([^"]*)"$/, async (attributeKey) => {
      const attrToDelete = ctx.createdAttributes.find(attr => attr.attribute_key === attributeKey);
      expect(attrToDelete).toBeDefined();
      ctx.attributeId = attrToDelete.id;
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);

      const endpoint = `/persons/${ctx.personId}/attributes/${ctx.attributeId}`;
      const requestBody = { meta: metaData };
      ctx.response = await sendDeleteRequest(endpoint, requestBody);

      if (ctx.response && ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = { success: true };
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and('the response should indicate success', () => {
      expect(ctx.responseData).toBeDefined();
    });

    and(/^the person should have (\d+) attributes? remaining$/, async (count) => {
      const attributes = await getPersonAttributes(ctx.personId);
      expect(attributes.length).toBe(parseInt(count));
    });

    and(/^the remaining attribute should be "([^"]*)"$/, async (key) => {
      const attributes = await getPersonAttributes(ctx.personId);
      expect(attributes.length).toBe(1);
      expect(attributes[0].attribute_key).toBe(key);
    });
  });

  // Scenario: Get attributes for a person with no attributes
  test('Get attributes for a person with no attributes', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and('the person has no attributes', async () => {
      const attributes = await getPersonAttributes(ctx.personId);
      expect(attributes.length).toBe(0);
    });

    when(/^I send a GET request to "\/persons\/\{personId\}\/attributes"$/, async () => {
      const endpoint = `/persons/${ctx.personId}/attributes`;
      await sendGet(ctx, endpoint);

      if (ctx.response && ctx.response.ok) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and('the response should contain an empty attributes array', () => {
      expect(ctx.responseData).toBeDefined();
      expect(Array.isArray(ctx.responseData)).toBe(true);
      expect(ctx.responseData.length).toBe(0);
    });
  });

  // Scenario: Attempt to add attribute to non-existent person
  test('Attempt to add attribute to non-existent person', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    when(/^I send a POST request to "\/persons\/99999\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      ctx.response = await sendPutRequest('/persons/99999/attributes', ctx.requestBody);

      if (ctx.response && !ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = null;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the error message should indicate "([^"]*)"$/, (errorMessage) => {
      expect(ctx.responseData).toBeDefined();
      const responseText = JSON.stringify(ctx.responseData).toLowerCase();
      expect(responseText).toContain(errorMessage.toLowerCase());
    });
  });

  // Scenario: Attempt to add attribute with missing required fields
  test('Attempt to add attribute with missing required fields', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with invalid data:$/, async (table) => {
      ctx.requestBody = { key: '' }; // Invalid/empty key
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && !ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = null;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and('the error should contain validation details', () => {
      expect(ctx.responseData).toBeDefined();
    });
  });

  // Scenario: Attempt to update non-existent attribute
  test('Attempt to update non-existent attribute', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a PUT request to "\/persons\/\{personId\}\/attributes\/99999" with:$/, async (table) => {
      const updateData = parseTableToObject(table);
      ctx.requestBody = {
        key: updateData.key,
        value: updateData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes/99999`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && !ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = null;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the error message should indicate "([^"]*)"$/, (errorMessage) => {
      expect(ctx.responseData).toBeDefined();
      const responseText = JSON.stringify(ctx.responseData).toLowerCase();
      expect(responseText).toContain(errorMessage.toLowerCase());
    });
  });

  // Scenario: Attempt to delete non-existent attribute
  test('Attempt to delete non-existent attribute', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a DELETE request to "\/persons\/\{personId\}\/attributes\/99999"$/, async () => {
      // Request body will be set in next step
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      const requestBody = { meta: metaData };

      const endpoint = `/persons/${ctx.personId}/attributes/99999`;
      ctx.response = await sendDeleteRequest(endpoint, requestBody);

      if (ctx.response && !ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = null;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the error message should indicate "([^"]*)"$/, (errorMessage) => {
      expect(ctx.responseData).toBeDefined();
      const responseText = JSON.stringify(ctx.responseData).toLowerCase();
      expect(responseText).toContain(errorMessage.toLowerCase());
    });
  });

  // Scenario: Get attributes for non-existent person
  test('Get attributes for non-existent person', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    when(/^I send a GET request to "\/persons\/99999\/attributes"$/, async () => {
      await sendGet(ctx, '/persons/99999/attributes');

      if (ctx.response && !ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = null;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the error message should indicate "([^"]*)"$/, (errorMessage) => {
      expect(ctx.responseData).toBeDefined();
      const responseText = JSON.stringify(ctx.responseData).toLowerCase();
      expect(responseText).toContain(errorMessage.toLowerCase());
    });
  });

  // Scenario: Attempt to add attribute without meta information
  test('Attempt to add attribute without meta information', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" without meta:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      const requestBody = {
        key: attributeData.key,
        value: attributeData.value
        // Intentionally omitting meta
      };

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, requestBody);

      if (ctx.response && !ctx.response.ok) {
        try {
          ctx.responseData = await parseJsonResponse(ctx.response);
        } catch (e) {
          ctx.responseData = null;
        }
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the error should indicate missing required field "([^"]*)"$/, (fieldName) => {
      expect(ctx.responseData).toBeDefined();
      const responseText = JSON.stringify(ctx.responseData).toLowerCase();
      expect(responseText).toContain(fieldName.toLowerCase());
    });
  });

  // Scenario: Update only the value of an attribute keeping the key same
  test('Update only the value of an attribute keeping the key same', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and(/^the person has an attribute:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.createdAttribute = await createAttribute(ctx.personId, attributeData);
      ctx.attributeId = ctx.createdAttribute.id;
    });

    when(/^I send a PUT request to "\/persons\/\{personId\}\/attributes\/\{attributeId\}" with:$/, async (table) => {
      const updateData = parseTableToObject(table);
      ctx.updateBody = {
        key: updateData.key,
        value: updateData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.updateBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes/${ctx.attributeId}`;
      ctx.response = await sendPutRequest(endpoint, ctx.updateBody);

      if (ctx.response && ctx.response.ok) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the attribute "([^"]*)" should have value "([^"]*)"$/, (key, value) => {
      expect(ctx.responseData).toBeDefined();
      expect(ctx.responseData.key).toBe(key);
      expect(ctx.responseData.value).toBe(value);
    });
  });

  // Scenario: Update attribute key and value
  test('Update attribute key and value', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and(/^the person has an attribute:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.createdAttribute = await createAttribute(ctx.personId, attributeData);
      ctx.attributeId = ctx.createdAttribute.id;
    });

    when(/^I send a PUT request to "\/persons\/\{personId\}\/attributes\/\{attributeId\}" with:$/, async (table) => {
      const updateData = parseTableToObject(table);
      ctx.updateBody = {
        key: updateData.key,
        value: updateData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.updateBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes/${ctx.attributeId}`;
      ctx.response = await sendPutRequest(endpoint, ctx.updateBody);

      if (ctx.response && ctx.response.ok) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the attribute should have key "([^"]*)" and value "([^"]*)"$/, (key, value) => {
      expect(ctx.responseData).toBeDefined();
      expect(ctx.responseData.key).toBe(key);
      expect(ctx.responseData.value).toBe(value);
    });
  });

  // Scenario: Complete attribute lifecycle - Create, Read, Update, Delete
  test('Complete attribute lifecycle - Create, Read, Update, Delete', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I add an attribute to the person:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      const meta = {
        caller: 'test',
        reason: 'lifecycle test',
        traceId: 'lifecycle-trace'
      };

      const requestBody = {
        key: attributeData.key,
        value: attributeData.value,
        meta: meta
      };

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
        ctx.attributeId = ctx.responseData.id;
        ctx.lifecycleAttribute = ctx.responseData;
      }
    });

    then('the attribute should be created successfully', () => {
      expect(ctx.response.status).toBe(201);
      expect(ctx.lifecycleAttribute).toBeDefined();
      expect(ctx.lifecycleAttribute.id).toBeDefined();
    });

    when('I retrieve all attributes for the person', async () => {
      const endpoint = `/persons/${ctx.personId}/attributes`;
      await sendGet(ctx, endpoint);

      if (ctx.response && ctx.response.ok) {
        ctx.allAttributes = await parseJsonResponse(ctx.response);
      }
    });

    then(/^I should see the "([^"]*)" attribute with value "([^"]*)"$/, (key, value) => {
      expect(ctx.allAttributes).toBeDefined();
      const found = ctx.allAttributes.find(attr => attr.key === key && attr.value === value);
      expect(found).toBeDefined();
    });

    when(/^I update the "([^"]*)" attribute to:$/, async (key, table) => {
      const updateData = parseTableToObject(table);
      const meta = {
        caller: 'test',
        reason: 'lifecycle update',
        traceId: 'lifecycle-update-trace'
      };

      const requestBody = {
        key: updateData.key,
        value: updateData.value,
        meta: meta
      };

      const endpoint = `/persons/${ctx.personId}/attributes/${ctx.attributeId}`;
      ctx.response = await sendPutRequest(endpoint, requestBody);

      if (ctx.response && ctx.response.ok) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then('the attribute should be updated successfully', () => {
      expect(ctx.response.status).toBe(200);
      expect(ctx.responseData).toBeDefined();
    });

    and(/^the "([^"]*)" attribute should have value "([^"]*)"$/, (key, value) => {
      expect(ctx.responseData.key).toBe(key);
      expect(ctx.responseData.value).toBe(value);
    });

    when(/^I delete the "([^"]*)" attribute$/, async (key) => {
      const meta = {
        caller: 'test',
        reason: 'lifecycle delete',
        traceId: 'lifecycle-delete-trace'
      };

      const endpoint = `/persons/${ctx.personId}/attributes/${ctx.attributeId}`;
      ctx.response = await sendDeleteRequest(endpoint, { meta });
    });

    then('the attribute should be deleted successfully', () => {
      expect(ctx.response.status).toBe(200);
    });

    and('the person should have no attributes', async () => {
      const attributes = await getPersonAttributes(ctx.personId);
      expect(attributes.length).toBe(0);
    });
  });

  // Scenario: Add attribute with special characters in value
  test('Add attribute with special characters in value', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.specialValue = attributeData.value;
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and('the attribute value should be stored correctly with special characters', () => {
      expect(ctx.responseData).toBeDefined();
      expect(ctx.responseData.value).toBe(ctx.specialValue);
    });
  });

  // Scenario: Add attribute with empty string value
  test('Add attribute with empty string value', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value || ''
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and('the attribute should be created with empty value', () => {
      expect(ctx.responseData).toBeDefined();
      expect(ctx.responseData.value).toBe('');
    });
  });

  // Scenario: Add multiple attributes with same key to same person
  // NOTE: Due to UNIQUE constraint on (person_id, attribute_key), 
  // adding an attribute with existing key will UPDATE it, not create a duplicate
  test('Add multiple attributes with same key to same person (upsert behavior)', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    and(/^the person has an attribute:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.firstAttribute = await createAttribute(ctx.personId, attributeData);
      ctx.attributeKey = attributeData.key;
      ctx.originalValue = attributeData.value;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
      ctx.newValue = attributeData.value;
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    // Due to UNIQUE constraint on (person_id, attribute_key), adding an attribute 
    // with existing key will UPDATE it (upsert), not create a duplicate
    and(/^the person should have (\d+) attribute with key "([^"]*)" and value "([^"]*)"$/, async (count, key, value) => {
      const attributes = await getPersonAttributes(ctx.personId);
      const matchingAttributes = attributes.filter(attr => attr.attribute_key === key);
      expect(matchingAttributes.length).toBe(parseInt(count));
      expect(matchingAttributes[0].attribute_value).toBe(value);
    });
  });

  // Scenario: Attribute value is stored encrypted in database
  test('Attribute value is stored encrypted in database', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.attributeKey = attributeData.key;
      ctx.attributeValue = attributeData.value;
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the raw database value for attribute "([^"]*)" should not equal "([^"]*)"$/, async (key, plainValue) => {
      const rawAttribute = await getRawAttributeFromDb(ctx.personId, key);
      expect(rawAttribute).toBeDefined();
      expect(rawAttribute.encrypted_value).toBeDefined();

      // Convert Buffer to string to verify it's not plain text
      const encryptedBytes = rawAttribute.encrypted_value;
      expect(encryptedBytes).not.toEqual(Buffer.from(plainValue));

      // Verify it's actually binary data (encrypted)
      expect(Buffer.isBuffer(encryptedBytes)).toBe(true);
    });

    and('the raw database value should be encrypted bytes', async () => {
      const rawAttribute = await getRawAttributeFromDb(ctx.personId, ctx.attributeKey);
      expect(rawAttribute).toBeDefined();
      expect(rawAttribute.encrypted_value).toBeDefined();
      expect(Buffer.isBuffer(rawAttribute.encrypted_value)).toBe(true);

      // Encrypted value should be different from plain text
      const encryptedString = rawAttribute.encrypted_value.toString('utf8');
      expect(encryptedString).not.toBe(ctx.attributeValue);
    });
  });

  // Scenario: Key version is stored correctly for each attribute
  test('Key version is stored correctly for each attribute', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.attributeKey = attributeData.key;
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the database should have key_version (\d+) for the attribute "([^"]*)"$/, async (expectedVersion, key) => {
      const rawAttribute = await getRawAttributeFromDb(ctx.personId, key);
      expect(rawAttribute).toBeDefined();
      expect(rawAttribute.key_version).toBeDefined();
      expect(rawAttribute.key_version).toBe(parseInt(expectedVersion));
    });
  });

  // Scenario: Decryption works correctly on retrieval
  test('Decryption works correctly on retrieval', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.attributeKey = attributeData.key;
      ctx.attributeValue = attributeData.value;
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);

      if (ctx.response && (ctx.response.ok || ctx.response.status === 201)) {
        ctx.responseData = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    when(/^I send a GET request to "\/persons\/\{personId\}\/attributes"$/, async () => {
      const endpoint = `/persons/${ctx.personId}/attributes`;
      await sendGet(ctx, endpoint);

      if (ctx.response && ctx.response.ok) {
        ctx.retrievedAttributes = await parseJsonResponse(ctx.response);
      }
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^the attribute "([^"]*)" should have value "([^"]*)"$/, (key, value) => {
      expect(ctx.retrievedAttributes).toBeDefined();
      const found = ctx.retrievedAttributes.find(attr => attr.key === key);
      expect(found).toBeDefined();
      expect(found.value).toBe(value);
    });

    and(/^the raw database value for attribute "([^"]*)" should not equal "([^"]*)"$/, async (key, plainValue) => {
      const rawAttribute = await getRawAttributeFromDb(ctx.personId, key);
      expect(rawAttribute).toBeDefined();
      expect(rawAttribute.encrypted_value).toBeDefined();

      // Convert Buffer to string to verify it's not plain text
      const encryptedBytes = rawAttribute.encrypted_value;
      expect(encryptedBytes).not.toEqual(Buffer.from(plainValue));

      // Verify it's actually binary data (encrypted)
      expect(Buffer.isBuffer(encryptedBytes)).toBe(true);
    });
  });

  // Scenario: Verify audit log creation
  test('Verify audit log creation', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;
      ctx.traceId = metaData.traceId;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);
      await parseJsonResponse(ctx.response);
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and(/^an audit record should be created for traceId "([^"]*)"$/, async (traceId) => {
      const query = 'SELECT * FROM request_log WHERE trace_id = $1';
      const result = await dbClient.query(query, [traceId]);
      expect(result.rows.length).toBeGreaterThan(0);
      ctx.auditRecord = result.rows[0];
    });

    and(/^the audit record should contain caller "([^"]*)" and reason "([^"]*)"$/, (caller, reason) => {
      expect(ctx.auditRecord.caller).toBe(caller);
      expect(ctx.auditRecord.reason).toBe(reason);
    });
  });

  // Scenario: Idempotency of request with same traceId
  test('Idempotency of request with same traceId', ({ given, when, then, and }) => {
    setupBackground({ given, and });

    given(/^a person exists with the following details:$/, async (table) => {
      const personData = parseTableToObject(table);
      ctx.createdPerson = await createPerson(personData);
      ctx.personId = ctx.createdPerson.id;
    });

    when(/^I send a POST request to "\/persons\/\{personId\}\/attributes" with:$/, async (table) => {
      const attributeData = parseTableToObject(table);
      ctx.requestBody = {
        key: attributeData.key,
        value: attributeData.value
      };
    });

    and(/^the request meta contains:$/, async (table) => {
      const metaData = parseTableToObject(table);
      ctx.requestBody.meta = metaData;
      ctx.traceId = metaData.traceId;

      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);
      await parseJsonResponse(ctx.response);
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    when(/^I send the same POST request again with traceId "([^"]*)"$/, async (traceId) => {
      const endpoint = `/persons/${ctx.personId}/attributes`;
      ctx.response = await sendPutRequest(endpoint, ctx.requestBody);
      ctx.responseData = await parseJsonResponse(ctx.response);
    });

    then(/^the response status should be (\d+)$/, (statusCode) => {
      assertResponseStatus(expect, ctx, parseInt(statusCode));
    });

    and('the attribute should be created only once', async () => {
      const attributes = await getPersonAttributes(ctx.personId);
      expect(attributes.length).toBe(1);
    });
  });
});
