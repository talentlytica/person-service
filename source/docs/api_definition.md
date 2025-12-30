# Person Service API - OpenAPI 3.0 Specification

```yaml
openapi: 3.0.3
info:
  title: Person Service API
  description: |
    RESTful API for managing persons and their attributes.
    This service provides CRUD operations for person entities with support for dynamic attributes.
  version: 1.0.0
  contact:
    name: API Support
    email: support@example.com

servers:
  - url: http://localhost:8080/api/v1
    description: Development server
  - url: https://api.example.com/v1
    description: Production server

tags:
  - name: persons
    description: Person management operations
  - name: attributes
    description: Person attribute management operations
  - name: health
    description: Health check endpoints

paths:
  /health:
    get:
      tags:
        - health
      summary: Health check endpoint
      description: Returns the health status of the service
      operationId: getHealth
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/HealthResponse'

  /persons:
    post:
      tags:
        - persons
      summary: Create a new person
      description: Creates a new person entity with the provided information
      operationId: createPerson
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreatePersonRequest'
      responses:
        '201':
          description: Person created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PersonResponse'
        '400':
          description: Invalid request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /persons/{id}:
    get:
      tags:
        - persons
      summary: Get person by ID
      description: Retrieves a person entity by their unique identifier
      operationId: getPersonById
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
        - name: X-Trace-Id
          in: header
          description: Trace ID for request tracking
          schema:
            type: string
            format: uuid
      responses:
        '200':
          description: Person found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PersonResponse'
        '404':
          description: Person not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    put:
      tags:
        - persons
      summary: Update person
      description: Updates an existing person entity
      operationId: updatePerson
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdatePersonRequest'
      responses:
        '200':
          description: Person updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PersonResponse'
        '400':
          description: Invalid request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: Person not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    delete:
      tags:
        - persons
      summary: Delete person
      description: Deletes a person entity and all associated attributes
      operationId: deletePerson
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DeletePersonRequest'
      responses:
        '200':
          description: Person deleted successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SuccessResponse'
        '404':
          description: Person not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /persons/{id}/attributes:
    get:
      tags:
        - attributes
      summary: Get all attributes for a person
      description: Retrieves all attributes associated with a person
      operationId: getPersonAttributes
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
      responses:
        '200':
          description: Attributes retrieved successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AttributesResponse'
        '404':
          description: Person not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    post:
      tags:
        - attributes
      summary: Add attribute to person
      description: Adds a new attribute to a person entity
      operationId: addPersonAttribute
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AddAttributeRequest'
      responses:
        '201':
          description: Attribute added successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AttributeResponse'
        '400':
          description: Invalid request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: Person not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

  /persons/{id}/attributes/{attributeId}:
    put:
      tags:
        - attributes
      summary: Update person attribute
      description: Updates an existing attribute of a person
      operationId: updatePersonAttribute
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
        - name: attributeId
          in: path
          required: true
          description: Unique identifier of the attribute
          schema:
            type: integer
            format: int64
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/UpdateAttributeRequest'
      responses:
        '200':
          description: Attribute updated successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/AttributeResponse'
        '400':
          description: Invalid request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '404':
          description: Attribute not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

    delete:
      tags:
        - attributes
      summary: Delete person attribute
      description: Deletes a specific attribute from a person
      operationId: deletePersonAttribute
      parameters:
        - name: id
          in: path
          required: true
          description: Unique identifier of the person
          schema:
            type: integer
            format: int64
        - name: attributeId
          in: path
          required: true
          description: Unique identifier of the attribute
          schema:
            type: integer
            format: int64
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/DeleteAttributeRequest'
      responses:
        '200':
          description: Attribute deleted successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SuccessResponse'
        '404':
          description: Attribute not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Internal server error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  schemas:
    Meta:
      type: object
      required:
        - traceId
      properties:
        caller:
          type: string
          description: Identifier of the caller
          example: "user123"
        reason:
          type: string
          description: Reason for the request
          example: "create new user"
        traceId:
          type: string
          format: uuid
          description: Unique trace ID for request tracking
          example: "550e8400-e29b-41d4-a716-446655440000"

    Person:
      type: object
      required:
        - name
        - clientId
      properties:
        id:
          type: integer
          format: int64
          description: Unique identifier of the person
          readOnly: true
          example: 1
        name:
          type: string
          description: Full name of the person
          example: "John Doe"
        clientId:
          type: string
          description: Client identifier
          example: "1234567890"
        createdAt:
          type: string
          format: date-time
          description: Timestamp when the person was created
          readOnly: true
          example: "2024-03-20T10:00:00Z"
        updatedAt:
          type: string
          format: date-time
          description: Timestamp when the person was last updated
          readOnly: true
          example: "2024-03-20T10:00:00Z"
        attributes:
          type: array
          description: List of attributes associated with the person
          items:
            $ref: '#/components/schemas/Attribute'

    Attribute:
      type: object
      required:
        - key
        - value
      properties:
        id:
          type: integer
          format: int64
          description: Unique identifier of the attribute
          readOnly: true
          example: 1
        key:
          type: string
          description: Attribute key/name
          example: "email"
        value:
          type: string
          description: Attribute value
          example: "john.doe@example.com"
        createdAt:
          type: string
          format: date-time
          description: Timestamp when the attribute was created
          readOnly: true
          example: "2024-03-20T10:00:00Z"
        updatedAt:
          type: string
          format: date-time
          description: Timestamp when the attribute was last updated
          readOnly: true
          example: "2024-03-20T10:00:00Z"

    Error:
      type: object
      properties:
        code:
          type: string
          description: Error code
          example: "VALIDATION_ERROR"
        message:
          type: string
          description: Human-readable error message
          example: "Invalid input parameters"
        details:
          type: array
          description: Additional error details
          items:
            type: string
          example: ["name is required", "clientId must be numeric"]

    CreatePersonRequest:
      type: object
      required:
        - meta
        - person
      properties:
        meta:
          $ref: '#/components/schemas/Meta'
        person:
          type: object
          required:
            - name
            - clientId
          properties:
            name:
              type: string
              example: "John Doe"
            clientId:
              type: string
              example: "1234567890"

    UpdatePersonRequest:
      type: object
      required:
        - meta
        - person
      properties:
        meta:
          $ref: '#/components/schemas/Meta'
        person:
          type: object
          properties:
            name:
              type: string
              example: "John Doe Updated"
            clientId:
              type: string
              example: "1234567890"

    DeletePersonRequest:
      type: object
      required:
        - meta
      properties:
        meta:
          $ref: '#/components/schemas/Meta'

    AddAttributeRequest:
      type: object
      required:
        - meta
        - attribute
      properties:
        meta:
          $ref: '#/components/schemas/Meta'
        attribute:
          type: object
          required:
            - key
            - value
          properties:
            key:
              type: string
              example: "email"
            value:
              type: string
              example: "john.doe@example.com"

    UpdateAttributeRequest:
      type: object
      required:
        - meta
        - attribute
      properties:
        meta:
          $ref: '#/components/schemas/Meta'
        attribute:
          type: object
          properties:
            key:
              type: string
              example: "email"
            value:
              type: string
              example: "john.updated@example.com"

    DeleteAttributeRequest:
      type: object
      required:
        - meta
      properties:
        meta:
          $ref: '#/components/schemas/Meta'

    PersonResponse:
      type: object
      properties:
        meta:
          type: object
          properties:
            traceId:
              type: string
              format: uuid
              example: "550e8400-e29b-41d4-a716-446655440000"
        person:
          $ref: '#/components/schemas/Person'
        error:
          $ref: '#/components/schemas/Error'

    AttributeResponse:
      type: object
      properties:
        meta:
          type: object
          properties:
            traceId:
              type: string
              format: uuid
              example: "550e8400-e29b-41d4-a716-446655440000"
        attribute:
          $ref: '#/components/schemas/Attribute'
        error:
          $ref: '#/components/schemas/Error'

    AttributesResponse:
      type: object
      properties:
        meta:
          type: object
          properties:
            traceId:
              type: string
              format: uuid
              example: "550e8400-e29b-41d4-a716-446655440000"
        attributes:
          type: array
          items:
            $ref: '#/components/schemas/Attribute'
        error:
          $ref: '#/components/schemas/Error'

    SuccessResponse:
      type: object
      properties:
        meta:
          type: object
          properties:
            traceId:
              type: string
              format: uuid
              example: "550e8400-e29b-41d4-a716-446655440000"
        success:
          type: boolean
          example: true
        message:
          type: string
          example: "Operation completed successfully"
        error:
          $ref: '#/components/schemas/Error'

    ErrorResponse:
      type: object
      properties:
        meta:
          type: object
          properties:
            traceId:
              type: string
              format: uuid
              example: "550e8400-e29b-41d4-a716-446655440000"
        error:
          $ref: '#/components/schemas/Error'

    HealthResponse:
      type: object
      properties:
        status:
          type: string
          enum: [healthy, unhealthy]
          example: "healthy"
        timestamp:
          type: string
          format: date-time
          example: "2024-03-20T10:00:00Z"
        version:
          type: string
          example: "1.0.0"
```

## Usage

You can use this OpenAPI specification to:

1. **Generate Server Code**: Use tools like `openapi-generator`, `swagger-codegen`, or language-specific generators
   ```bash
   # Example for Go
   openapi-generator generate -i api_definition.md -g go-server -o ./generated

   # Example for Node.js/Express
   openapi-generator generate -i api_definition.md -g nodejs-express-server -o ./generated
   ```

2. **Generate Client SDKs**: Create client libraries in various languages
   ```bash
   openapi-generator generate -i api_definition.md -g javascript -o ./client-sdk
   ```

3. **API Documentation**: Use Swagger UI or ReDoc to generate interactive documentation
   ```bash
   # Using Docker with Swagger UI
   docker run -p 80:8080 -e SWAGGER_JSON=/api/api_definition.md -v $(pwd):/api swaggerapi/swagger-ui
   ```

4. **API Testing**: Import into Postman, Insomnia, or use automated testing tools

5. **Validation**: Validate requests/responses against the schema in your application
