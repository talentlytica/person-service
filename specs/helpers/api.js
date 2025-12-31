import fetch from 'node-fetch';

/**
 * Reusable API Helper Module
 * Provides common HTTP operations for integration tests
 */

/**
 * Get the service URL from global test environment
 * @returns {string} Service URL
 */
export function getServiceUrl() {
  if (!global.__TEST_ENV__) {
    throw new Error('Test environment not initialized');
  }
  return global.__TEST_ENV__.getServiceUrl();
}

/**
 * Send a GET request to the service
 * @param {string} path - API endpoint path
 * @param {Object} options - Fetch options (optional)
 * @returns {Promise<Response>} Fetch response
 */
export async function sendGetRequest(path, options = {}) {
  const serviceUrl = getServiceUrl();
  const url = `${serviceUrl}${path}`;
  return await fetch(url, {
    method: 'GET',
    ...options
  });
}

/**
 * Send a GET request with timeout
 * @param {string} path - API endpoint path
 * @param {number} timeoutMs - Timeout in milliseconds
 * @returns {Promise<{response: Response, timedOut: boolean}>}
 */
export async function sendGetRequestWithTimeout(path, timeoutMs) {
  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeoutMs);
  
  try {
    const response = await sendGetRequest(path, {
      signal: controller.signal
    });
    clearTimeout(timeoutId);
    return { response, timedOut: false };
  } catch (error) {
    clearTimeout(timeoutId);
    if (error.name === 'AbortError') {
      return { response: null, timedOut: true };
    }
    throw error;
  }
}

/**
 * Send a POST request to the service
 * @param {string} path - API endpoint path
 * @param {Object} body - Request body (will be JSON stringified)
 * @param {Object} options - Fetch options (optional)
 * @returns {Promise<Response>} Fetch response
 */
export async function sendPostRequest(path, body = {}, options = {}) {
  const serviceUrl = getServiceUrl();
  const url = `${serviceUrl}${path}`;
  return await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    },
    body: JSON.stringify(body),
    ...options
  });
}

/**
 * Send a PUT request to the service
 * @param {string} path - API endpoint path
 * @param {Object} body - Request body (will be JSON stringified)
 * @param {Object} options - Fetch options (optional)
 * @returns {Promise<Response>} Fetch response
 */
export async function sendPutRequest(path, body = {}, options = {}) {
  const serviceUrl = getServiceUrl();
  const url = `${serviceUrl}${path}`;
  return await fetch(url, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    },
    body: JSON.stringify(body),
    ...options
  });
}

/**
 * Send a DELETE request to the service
 * @param {string} path - API endpoint path
 * @param {Object} body - Request body (will be JSON stringified, optional)
 * @param {Object} options - Fetch options (optional)
 * @returns {Promise<Response>} Fetch response
 */
export async function sendDeleteRequest(path, body = null, options = {}) {
  const serviceUrl = getServiceUrl();
  const url = `${serviceUrl}${path}`;
  const fetchOptions = {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    },
    ...options
  };
  
  if (body !== null) {
    fetchOptions.body = JSON.stringify(body);
  }
  
  return await fetch(url, fetchOptions);
}

/**
 * Send multiple concurrent GET requests
 * @param {string} path - API endpoint path
 * @param {number} count - Number of concurrent requests
 * @returns {Promise<Response[]>} Array of responses
 */
export async function sendConcurrentGetRequests(path, count) {
  const requests = Array(count).fill(null).map(() => sendGetRequest(path));
  return await Promise.all(requests);
}

/**
 * Parse JSON response body
 * @param {Response} response - Fetch response
 * @returns {Promise<any>} Parsed JSON data
 */
export async function parseJsonResponse(response) {
  return await response.json();
}

/**
 * Get content type from response headers
 * @param {Response} response - Fetch response
 * @returns {string} Content type header value
 */
export function getContentType(response) {
  return response.headers.get('content-type') || '';
}

/**
 * Check if response is valid JSON
 * @param {Response} response - Fetch response
 * @returns {Promise<boolean>} True if valid JSON
 */
export async function isValidJson(response) {
  try {
    await response.clone().json();
    return true;
  } catch {
    return false;
  }
}

