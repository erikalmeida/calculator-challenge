# Sezzle Calculator

A full-stack calculator for the take-home assignment. The REST API is written in Go with the standard library, and the user interface is built with React and TypeScript.

## Requirements

- Go 1.27 or newer.
- Node.js 26 LTS or newer.
- npm 11 or newer.
- `curl` for optional API checks.
- Docker Desktop with Docker Compose v2 for the optional containerized setup.

## Run with Docker

From the repository root, build and start both services:

```bash
docker compose up --build
```

- Frontend: `http://localhost:3000`
- Backend API and Swagger UI: `http://localhost:8080/docs`
- Backend health check: `http://localhost:8080/health`

The frontend container uses Nginx to serve the React build and proxy `/api` requests to the backend container. This keeps the browser on one origin and avoids an additional CORS configuration for the containerized setup.

Stop the services with:

```bash
docker compose down
```

## Run locally

Start the backend in one terminal:

```bash
cd backend
go run ./cmd/api
```

Start the frontend in another terminal:

```bash
cd frontend
npm install
npm run dev
```

Open the URL printed by Vite, normally `http://localhost:5173`. During development, Vite proxies `/api` requests to `http://localhost:8080`.

`npm run dev` starts the frontend locally. `npm run build` does not start a server; it only validates the TypeScript project and creates an optimized production build in `frontend/dist`.

To inspect that production build locally after running `npm run build`, use:

```bash
cd frontend
npm run preview
```

## Interactive API documentation

With the backend running, open Swagger UI at:

```text
http://localhost:8080/docs
```

The interface reads the OpenAPI 3.1 specification served by the API at:

```text
http://localhost:8080/api/v1/openapi.yaml
```

Swagger UI assets are loaded from the Swagger UI CDN, while the API contract itself is versioned locally at `backend/internal/adapters/in/httpapi/docs/openapi.yaml`.

## API

```http
GET /health
```

Health response (`200`):

```json
{
  "status": "ok"
}
```

```http
POST /api/v1/calculate
Content-Type: application/json
```

Example request:

```bash
curl -i -X POST http://localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"multiply","operands":[2,3,4]}'
```

Successful response (`200`):

```json
{
  "operation": "multiply",
  "operands": [2, 3, 4],
  "result": 24
}
```

Error response (`400`), for example division by zero:

```json
{
  "error": {
    "code": "division_by_zero",
    "message": "The divisor cannot be zero."
  }
}
```

Every response includes a generated `X-Request-ID` header. Include this value when reporting a problem so the corresponding server log entry can be identified quickly.

## Error handling and logging

The API uses a consistent JSON error envelope. Expected client errors return `400`, an unsupported `Content-Type` returns `415`, an oversized body returns `413`, and unexpected failures return a safe `500 internal_error` response.

The HTTP adapter includes recovery middleware. If an unexpected Go `panic` occurs, the server logs the technical error and stack trace, keeps the request ID, and returns a safe JSON response without exposing internals to the client.

The backend writes structured JSON logs for server lifecycle events and completed requests. Request logs contain request ID, method, path, status, duration, error code, operation, and operand count. They deliberately exclude operand values and request bodies.

## Calculator contract

| Operation | Operand count | Example |
| --- | --- | --- |
| `add` | 2 to 20 | `[1, 2, 3]` → `6` |
| `subtract` | 2 to 20 | `[10, 3, 2]` → `5` |
| `multiply` | 2 to 20 | `[2, 3, 4]` → `24` |
| `divide` | exactly 2 | `[42, 6]` → `7` |

The API accepts JSON decimal numbers only. Use an optional leading minus sign and a period as the decimal separator:

```text
Accepted values: `12` · `-12` · `0.5` · `-0.5` · `0.0000000001`
Rejected values: `1,000` · `1,5` · `1e3` · `null`
```

The following limits intentionally keep the `float64` contract readable and predictable:

| Rule | Value |
| --- | --- |
| Maximum operands | 20 |
| Maximum absolute operand | `100000000000000` (`1e14`) |
| Minimum non-zero absolute operand | `0.0000000001` (`1e-10`) |
| Maximum decimal places | 10 |
| Maximum significant digits | 15 |
| Maximum operand input length | 17 characters |
| Maximum request body | 1 MiB |

`float64` is appropriate for this calculator assignment, but it is not an exact decimal representation. A financial production system should use a money type based on minor units or arbitrary-precision decimal arithmetic with explicit rounding rules.

## Tests and quality checks

Backend, from `backend/`:

```bash
go test ./...
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go vet ./...
```

Frontend, from `frontend/`:

```bash
npm run lint
npm run test -- --run
npm run test:coverage
npm run build
```

`npm run build` is a release-quality check, not the command used to start the application during development.

## AI usage

See [AI_USAGE.md](AI_USAGE.md) for a concise disclosure of how AI assisted with planning, implementation, review, and learning during this take-home project.

## Architecture

The backend uses a lightweight hexagonal architecture:

```text
HTTP/JSON inbound adapter
        ↓
CalculateUseCase application port
        ↓
CalculationStrategyResolver output port
        ↓
InMemoryCalculationStrategyRegistry
        ↓
AdditionStrategy / SubtractionStrategy / MultiplicationStrategy / DivisionStrategy
```

- Each operation is an independent `CalculationStrategy`.
- The application service depends on `CalculationStrategyResolver`, not on the in-memory registry implementation.
- The HTTP adapter depends only on `CalculateUseCase`.
- The backend is the source of truth for business rules and validation; the frontend mirrors input rules only to provide faster feedback.
- The Docker setup uses a multi-stage Go build and an Nginx runtime that proxies browser API requests to the backend container.
- No external Go dependencies, database, framework, remote repository, or commit have been created for this project.
