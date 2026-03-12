/**
 * Tests for GET /version - Service Version Check
 *
 * Gherkin:
 *   Feature: Service Version
 *     Scenario: Check current service version
 *       Given the person-service is running
 *       When I send a GET request to "/version"
 *       Then the response status should be 200
 *       And the response should contain a "version" field
 *       And the version should match the VERSION file
 *
 *     Scenario: Version format is valid semver with v prefix
 *       Given the person-service is running
 *       When I send a GET request to "/version"
 *       Then the version should match the pattern "v<major>.<minor>.<patch>"
 */

import { describe, test, expect, beforeAll } from '@jest/globals';
import axios from 'axios';
import { readFileSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';
import { config } from 'dotenv';

config();

const __dirname = dirname(fileURLToPath(import.meta.url));

describe('GET /version - Service Version', () => {
  const BASE_URL = process.env.BASE_URL || 'http://localhost:3000';
  let expectedVersion;

  beforeAll(() => {
    const versionFile = resolve(__dirname, '../../../VERSION');
    expectedVersion = readFileSync(versionFile, 'utf-8').trim();
  });

  test('Version endpoint returns 200 with version field', async () => {
    const response = await axios.get(`${BASE_URL}/version`);

    expect(response.status).toBe(200);
    expect(response.data).toHaveProperty('version');
  });

  test('Version matches the VERSION file', async () => {
    const response = await axios.get(`${BASE_URL}/version`);

    expect(response.data.version).toBe(expectedVersion);
  });

  test('Version follows v<major>.<minor>.<patch> format', async () => {
    const response = await axios.get(`${BASE_URL}/version`);

    expect(response.data.version).toMatch(/^v\d+\.\d+\.\d+$/);
  });
});
