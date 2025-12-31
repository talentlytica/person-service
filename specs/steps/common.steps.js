/**
 * Common Step Helpers
 * 
 * Reusable helper functions for HTTP API testing with JSON responses.
 * These helpers can be imported and used across multiple feature step definitions.
 */

import {
  sendGetRequest,
  sendPostRequest,
  sendPutRequest,
  sendDeleteRequest,
  sendGetRequestWithTimeout,
  sendConcurrentGetRequests,
  parseJsonResponse,
  getContentType,
  isValidJson
} from '../helpers/api.js';

/**
 * Create a shared test context for storing request/response state
 * @returns {Object} Test context object
 */
export function createTestContext() {
  return {
    response: null,
    responses: [],
    responseData: null,
    timedOut: false
  };
}

/**
 * Step helper: Assert service is running
 * @param {Function} expect - Jest expect function
 */
export function assertServiceRunning(expect) {
  expect(global.__TEST_ENV__).toBeDefined();
}

/**
 * Step helper: Send GET request and store in context
 * @param {Object} ctx - Test context
 * @param {string} path - API endpoint path
 */
export async function sendGet(ctx, path) {
  ctx.response = await sendGetRequest(path);
}

/**
 * Step helper: Send POST request and store in context
 * @param {Object} ctx - Test context
 * @param {string} path - API endpoint path
 * @param {Object} body - Request body
 */
export async function sendPost(ctx, path, body = {}) {
  ctx.response = await sendPostRequest(path, body);
  ctx.responseData = null; // Clear cached response data
}

/**
 * Step helper: Send PUT request and store in context
 * @param {Object} ctx - Test context
 * @param {string} path - API endpoint path
 * @param {Object} body - Request body
 */
export async function sendPut(ctx, path, body = {}) {
  ctx.response = await sendPutRequest(path, body);
  ctx.responseData = null; // Clear cached response data
}

/**
 * Step helper: Send DELETE request and store in context
 * @param {Object} ctx - Test context
 * @param {string} path - API endpoint path
 * @param {Object} body - Request body (optional)
 */
export async function sendDelete(ctx, path, body = null) {
  ctx.response = await sendDeleteRequest(path, body);
  ctx.responseData = null; // Clear cached response data
}

/**
 * Step helper: Send GET request with timeout and store in context
 * @param {Object} ctx - Test context
 * @param {string} path - API endpoint path
 * @param {number} timeoutMs - Timeout in milliseconds
 */
export async function sendGetWithTimeout(ctx, path, timeoutMs) {
  const result = await sendGetRequestWithTimeout(path, timeoutMs);
  ctx.response = result.response;
  ctx.timedOut = result.timedOut;
  ctx.responseData = null;
}

/**
 * Step helper: Send concurrent GET requests and store in context
 * @param {Object} ctx - Test context
 * @param {string} path - API endpoint path
 * @param {number} count - Number of concurrent requests
 */
export async function sendConcurrentGet(ctx, path, count) {
  ctx.responses = await sendConcurrentGetRequests(path, count);
}

/**
 * Step helper: Assert response status code
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 * @param {number} expectedStatus - Expected HTTP status code
 */
export function assertResponseStatus(expect, ctx, expectedStatus) {
  expect(ctx.response).not.toBeNull();
  expect(ctx.response.status).toBe(expectedStatus);
}

/**
 * Step helper: Assert response content type
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 * @param {string} expectedType - Expected content type pattern
 */
export function assertContentType(expect, ctx, expectedType) {
  const contentType = ctx.response.headers.get('content-type') || '';
  expect(contentType).toMatch(new RegExp(expectedType));
}

/**
 * Step helper: Assert response is valid JSON
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 */
export async function assertValidJson(expect, ctx) {
  const valid = await isValidJson(ctx.response);
  expect(valid).toBe(true);
}

/**
 * Step helper: Assert response is an object
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 */
export async function assertResponseIsObject(expect, ctx) {
  if (!ctx.responseData) {
    ctx.responseData = await parseJsonResponse(ctx.response);
  }
  expect(typeof ctx.responseData).toBe('object');
  expect(ctx.responseData).not.toBeNull();
  expect(Array.isArray(ctx.responseData)).toBe(false);
}

/**
 * Step helper: Assert response field contains value
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 * @param {string} field - Field name
 * @param {string} value - Expected value
 */
export async function assertFieldValue(expect, ctx, field, value) {
  if (!ctx.responseData) {
    ctx.responseData = await parseJsonResponse(ctx.response);
  }
  expect(ctx.responseData[field]).toBe(value);
}

/**
 * Step helper: Assert response has field
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 * @param {string} field - Field name
 */
export async function assertHasField(expect, ctx, field) {
  if (!ctx.responseData) {
    ctx.responseData = await parseJsonResponse(ctx.response);
  }
  expect(ctx.responseData).toHaveProperty(field);
}

/**
 * Step helper: Assert field is one of allowed values
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 * @param {string} field - Field name
 * @param {string} valuesStr - Comma-separated allowed values
 */
export async function assertFieldOneOf(expect, ctx, field, valuesStr) {
  if (!ctx.responseData) {
    ctx.responseData = await parseJsonResponse(ctx.response);
  }
  const allowedValues = valuesStr.split(',').map(v => v.trim());
  expect(allowedValues).toContain(ctx.responseData[field]);
}

/**
 * Step helper: Assert all responses have status code
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 * @param {number} expectedStatus - Expected HTTP status code
 */
export function assertAllResponsesStatus(expect, ctx, expectedStatus) {
  ctx.responses.forEach((response) => {
    expect(response.status).toBe(expectedStatus);
  });
}

/**
 * Step helper: Assert request completed within timeout
 * @param {Function} expect - Jest expect function
 * @param {Object} ctx - Test context
 */
export function assertNoTimeout(expect, ctx) {
  expect(ctx.timedOut).toBe(false);
}
