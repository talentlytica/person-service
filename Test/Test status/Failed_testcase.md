# Failed Test Cases

**Total: 21**

## 1. SPEC: Should handle Unicode and special characters

- **File:** `API_Tests/specification_tests/person_attributes_security_spec.test.js`
- **Suite:** SPECIFICATION: Person Attributes - Security & Edge Cases
- **Duration:** 121 ms

**Failure message:**
```
AxiosError: Request failed with status code 500
    at settle (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\settle.js:19:12)
    at IncomingMessage.handleStreamEnd (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\adapters\http.js:793:11)
    at IncomingMessage.emit (node:events:530:35)
    at endReadableNT (node:internal/streams/readable:1698:12)
    at processTicksAndRejections (node:internal/process/task_queues:90:21)
    at Axios.request (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\Axios.js:45:41)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\specification_tests\person_attributes_security_spec.test.js:282:24)
```

## 2. Should return 404 for non-existent personId

- **File:** `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: PUT /persons/:personId/attributes/:attributeId
- **Duration:** 4 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\PUT_persons_personId_attributes_attributeId_negative.test.js:111:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 3. Should return 404 for non-existent attributeId

- **File:** `API_Tests/negative_tests/PUT_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: PUT /persons/:personId/attributes/:attributeId
- **Duration:** 2 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\PUT_persons_personId_attributes_attributeId_negative.test.js:123:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 4. Should handle or reject extremely long key (10KB)

- **File:** `API_Tests/negative_tests/POST_api_key-value_negative.test.js`
- **Suite:** NEGATIVE: POST /api/key-value - Create/Update Key-Value
- **Duration:** 7 ms

**Failure message:**
```
Error: expect(received).toContain(expected) // indexOf

Expected value: 500
Received array: [400, 413]
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\POST_api_key-value_negative.test.js:240:26)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 5. Should handle or reject extremely long value (1MB)

- **File:** `API_Tests/negative_tests/POST_api_key-value_negative.test.js`
- **Suite:** NEGATIVE: POST /api/key-value - Create/Update Key-Value
- **Duration:** 16 ms

**Failure message:**
```
Error: expect(received).toContain(expected) // indexOf

Expected value: 500
Received array: [400, 413]
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\POST_api_key-value_negative.test.js:256:26)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 6. Should reject key with only whitespace

- **File:** `API_Tests/negative_tests/POST_api_key-value_negative.test.js`
- **Suite:** NEGATIVE: POST /api/key-value - Create/Update Key-Value
- **Duration:** 5 ms

**Failure message:**
```
TypeError: Cannot read properties of undefined (reading 'status')
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\POST_api_key-value_negative.test.js:268:29)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 7. Should return 404 for non-existent personId

- **File:** `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: DELETE /persons/:personId/attributes/:attributeId
- **Duration:** 4 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:107:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 8. Should return 404 for non-existent attributeId

- **File:** `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: DELETE /persons/:personId/attributes/:attributeId
- **Duration:** 2 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:118:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 9. Should handle DELETE of already deleted attribute (idempotent)

- **File:** `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: DELETE /persons/:personId/attributes/:attributeId
- **Duration:** 2 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:130:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 10. Should ignore request body on DELETE

- **File:** `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: DELETE /persons/:personId/attributes/:attributeId
- **Duration:** 3 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:219:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 11. Should handle multiple concurrent DELETEs of same attribute

- **File:** `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: DELETE /persons/:personId/attributes/:attributeId
- **Duration:** 5 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:314:47
    at Array.forEach (<anonymous>)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:312:13)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 12. Should ignore invalid Content-Type header on DELETE

- **File:** `API_Tests/negative_tests/DELETE_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: DELETE /persons/:personId/attributes/:attributeId
- **Duration:** 2 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\DELETE_persons_personId_attributes_attributeId_negative.test.js:336:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 13. Should return 404 for non-existent personId

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes/:attributeId
- **Duration:** 3 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_attributeId_negative.test.js:94:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 14. Should return 404 for non-existent attributeId

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes/:attributeId
- **Duration:** 1 ms

**Failure message:**
```
Error: expect(received).toBe(expected) // Object.is equality

Expected: 404
Received: 400
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_attributeId_negative.test.js:105:37)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 15. Should ignore request body on GET

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes/:attributeId
- **Duration:** 2 ms

**Failure message:**
```
AxiosError: Request failed with status code 400
    at settle (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\settle.js:19:12)
    at IncomingMessage.handleStreamEnd (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\adapters\http.js:793:11)
    at IncomingMessage.emit (node:events:530:35)
    at endReadableNT (node:internal/streams/readable:1698:12)
    at processTicksAndRejections (node:internal/process/task_queues:90:21)
    at Axios.request (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\Axios.js:45:41)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_attributeId_negative.test.js:234:22)
```

## 16. Should handle invalid Accept header

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_attributeId_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes/:attributeId
- **Duration:** 7 ms

**Failure message:**
```
Error: expect(received).toContain(expected) // indexOf

Expected value: 400
Received array: [404, 406]
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_attributeId_negative.test.js:257:28)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

## 17. Should return empty or 404 for non-existent personId

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes - Get All Attributes
- **Duration:** 2 ms

**Failure message:**
```
AxiosError: Request failed with status code 404
    at settle (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\settle.js:19:12)
    at IncomingMessage.handleStreamEnd (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\adapters\http.js:793:11)
    at IncomingMessage.emit (node:events:530:35)
    at endReadableNT (node:internal/streams/readable:1698:12)
    at processTicksAndRejections (node:internal/process/task_queues:90:21)
    at Axios.request (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\Axios.js:45:41)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_negative.test.js:84:22)
```

## 18. Should reject request with body (GET should not have body)

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes - Get All Attributes
- **Duration:** 3 ms

**Failure message:**
```
AxiosError: Request failed with status code 404
    at settle (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\settle.js:19:12)
    at IncomingMessage.handleStreamEnd (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\adapters\http.js:793:11)
    at IncomingMessage.emit (node:events:530:35)
    at endReadableNT (node:internal/streams/readable:1698:12)
    at processTicksAndRejections (node:internal/process/task_queues:90:21)
    at Axios.request (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\Axios.js:45:41)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_negative.test.js:150:22)
```

## 19. Should handle invalid Accept header

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes - Get All Attributes
- **Duration:** 2 ms

**Failure message:**
```
AxiosError: Request failed with status code 404
    at settle (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\settle.js:19:12)
    at IncomingMessage.handleStreamEnd (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\adapters\http.js:793:11)
    at IncomingMessage.emit (node:events:530:35)
    at endReadableNT (node:internal/streams/readable:1698:12)
    at processTicksAndRejections (node:internal/process/task_queues:90:21)
    at Axios.request (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\Axios.js:45:41)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_negative.test.js:163:22)
```

## 20. Should ignore invalid query parameters

- **File:** `API_Tests/negative_tests/GET_persons_personId_attributes_negative.test.js`
- **Suite:** NEGATIVE: GET /persons/:personId/attributes - Get All Attributes
- **Duration:** 2 ms

**Failure message:**
```
AxiosError: Request failed with status code 404
    at settle (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\settle.js:19:12)
    at IncomingMessage.handleStreamEnd (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\adapters\http.js:793:11)
    at IncomingMessage.emit (node:events:530:35)
    at endReadableNT (node:internal/streams/readable:1698:12)
    at processTicksAndRejections (node:internal/process/task_queues:90:21)
    at Axios.request (C:\RepoGit\person-service - v2\Test\node_modules\axios\lib\core\Axios.js:45:41)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_persons_personId_attributes_negative.test.js:179:22)
```

## 21. Should handle very long key in URL

- **File:** `API_Tests/negative_tests/GET_api_key-value_key_negative.test.js`
- **Suite:** NEGATIVE: GET /api/key-value/:key - Retrieve Value
- **Duration:** 5 ms

**Failure message:**
```
Error: expect(received).toContain(expected) // indexOf

Expected value: 500
Received array: [400, 404, 414]
    at Object.<anonymous> (C:\RepoGit\person-service - v2\Test\API_Tests\negative_tests\GET_api_key-value_key_negative.test.js:179:31)
    at processTicksAndRejections (node:internal/process/task_queues:105:5)
```

