Person Attributes Web UI

Quick static UI to list, create, update and delete person attributes.

Files
- [source/webui/index.html](source/webui/index.html)
- [source/webui/app.js](source/webui/app.js)
- [source/webui/styles.css](source/webui/styles.css)

Usage
1. Serve the folder (recommended) to avoid CORS issues. From the `source/webui` folder run:

```bash
# Python 3 simple server
python3 -m http.server 8000
# then open http://localhost:8000/
```

2. Enter a `Person ID` (UUID) and optionally an `API base` (for example `http://localhost:8080/`). Click `Load`.
3. Use the create form to add a new attribute. Use `Save` / `Delete` buttons on each attribute row.

Notes
- The UI expects the service endpoints described in the project:
  - GET  /persons/:personId/attributes
  - POST /persons/:personId/attributes  (body {key,value,meta})
  - PUT  /persons/:personId/attributes/:attributeId
  - DELETE /persons/:personId/attributes/:attributeId
- `meta` is auto-populated by the UI with a simple traceId; adjust as needed.

Makefile helper

You can start the backend and the proxied web UI together using the top-level `Makefile` target:

```bash
make webui
```

Stop the services with:

```bash
make webui-down
```
