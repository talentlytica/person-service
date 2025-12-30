# API Documentation

This directory contains the API documentation for the Person Service.

## Files

- `openapi.yaml` - OpenAPI 3.0 specification
- `swagger-ui.html` - Interactive API documentation UI
- `api_definition.md` - Complete API specification with usage instructions

## Viewing Interactive Documentation

### Option 1: Using Docker (Recommended)

Run Swagger UI with Docker:

```bash
# From the project root
docker run -p 8081:8080 \
  -e SWAGGER_JSON=/docs/openapi.yaml \
  -v $(pwd)/source/docs:/docs \
  swaggerapi/swagger-ui
```

Then open: http://localhost:8081

### Option 2: Using Python HTTP Server

If you have Python installed:

```bash
# Navigate to the docs directory
cd source/docs

# Start a simple HTTP server
python3 -m http.server 8082
```

Then open: http://localhost:8082/swagger-ui.html

### Option 3: Using Node.js HTTP Server

If you have Node.js installed:

```bash
# Install http-server globally (one-time)
npm install -g http-server

# Navigate to the docs directory
cd source/docs

# Start the server
http-server -p 8082
```

Then open: http://localhost:8082/swagger-ui.html

### Option 4: Using Go HTTP Server

Create a simple Go server:

```bash
# Navigate to docs directory
cd source/docs

# Run a simple file server
go run -e 'package main; import ("net/http"; "log"); func main() { log.Println("Server starting on :8082"); http.Handle("/", http.FileServer(http.Dir("."))); log.Fatal(http.ListenAndServe(":8082", nil)) }'
```

Or use the included Makefile command (if available):

```bash
make docs
```

## Features of Interactive Documentation

The Swagger UI provides:

- 📖 **Interactive API exploration** - Try out API endpoints directly from the browser
- 🔍 **Search and filter** - Quickly find specific endpoints
- 📝 **Request/Response examples** - See sample data for all operations
- ✅ **Schema validation** - Understand required fields and data types
- 🎨 **Syntax highlighting** - Easy-to-read JSON/YAML
- 📥 **Download OpenAPI spec** - Export the specification

## Using the API Documentation

1. **Explore endpoints** - Click on any endpoint to expand details
2. **Try it out** - Click "Try it out" button to test the API
3. **Fill parameters** - Enter required parameters and request body
4. **Execute** - Click "Execute" to send the request
5. **View response** - See the actual response from your server

## Code Generation

Use the `openapi.yaml` file to generate code:

```bash
# Generate Go server
openapi-generator generate -i openapi.yaml -g go-server -o ./generated

# Generate JavaScript client
openapi-generator generate -i openapi.yaml -g javascript -o ./client-sdk

# Generate TypeScript client
openapi-generator generate -i openapi.yaml -g typescript-axios -o ./client-sdk
```

## Importing into Tools

### Postman
1. Open Postman
2. Click "Import"
3. Select `openapi.yaml`
4. All endpoints will be imported as a collection

### Insomnia
1. Open Insomnia
2. Click "Import/Export"
3. Select "Import Data" → "From File"
4. Choose `openapi.yaml`

### VS Code REST Client
Use the OpenAPI specification with REST Client extension for inline testing.

## Updating Documentation

When you update the API:

1. Edit `openapi.yaml` with your changes
2. Refresh the browser to see updated documentation
3. The Swagger UI will automatically reload the specification

## API Server Configuration

Make sure your API server is running and accessible at the configured URL:
- Development: `http://localhost:8080/api/v1`
- Production: Update the `servers` section in `openapi.yaml`

## Support

For questions or issues with the API, contact: support@example.com

