# Backend

Go REST service that performs the calculator's arithmetic. Built with the
standard library only (`net/http`); no third-party dependencies.

## Setup

Requires Go 1.22 or newer.

```bash
cd backend
go run ./cmd/server
```

The server listens on port 8080 by default.

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Port to listen on |
| `ALLOWED_ORIGINS` | `*` | Comma-separated list of allowed CORS origins |

Example:

```bash
PORT=9000 ALLOWED_ORIGINS=http://localhost:5173 go run ./cmd/server
```

Stop the server with Ctrl+C; it shuts down gracefully, finishing in-flight
requests before exiting.

## Running with Docker

```bash
docker build -t sezzle-backend .
docker run -p 8080:8080 sezzle-backend
```

To change the port or CORS origins at runtime:

```bash
docker run -p 8080:8080 -e ALLOWED_ORIGINS=http://localhost:3000 sezzle-backend
```

The `Dockerfile` is a multi-stage build: it compiles a static binary in a
`golang` build stage, then copies only that binary into a minimal
`gcr.io/distroless/static-debian12` runtime image, which has no shell and no
package manager.

## Project structure

```
cmd/server/main.go        entrypoint: env config, graceful shutdown
internal/
  server/                 router and middleware (CORS, logging)
  handlers/               HTTP handlers, request validation
  calculator/             arithmetic engine, no HTTP knowledge
  httpx/                  shared JSON response helpers
```

`internal/calculator` never imports `net/http`. It exposes a single
function, `Calculate(operation, a, b)`, that returns a result or one of
three sentinel errors (`ErrDivisionByZero`, `ErrNegativeSqrt`,
`ErrUnsupportedOperation`). Keeping the arithmetic free of HTTP concerns is
what makes it trivial to unit test in isolation.

`internal/handlers` owns everything HTTP-specific: decoding the JSON body,
validating required fields, and mapping domain errors to status codes.

## API

### GET /health

Liveness probe. Returns `{"status": "ok"}`.

### GET /api/operations

Lists supported operation ids: `{"operations": ["add", "subtract", ...]}`.

### POST /api/calculate

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `operation` | string | yes | One of the ids from `/api/operations` |
| `a` | number | yes | First operand |
| `b` | number | yes, except `sqrt` | Second operand; omit for `sqrt` |

Success response:

```json
{ "operation": "add", "a": 4, "b": 2, "result": 6 }
```

Error response, same envelope for every failure:

```json
{ "error": "division by zero is not allowed" }
```

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Malformed JSON, missing/unknown fields, unsupported operation |
| 422 | Well-formed but mathematically invalid (division by zero, square root of a negative number, non-finite result) |

See the [root README](../README.md#api) for curl examples of every
operation.

## Tests

```bash
go test ./...              # run all tests
go test ./... -v           # verbose output
go test ./... -cover       # coverage summary per package
```

Coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Formatting and static checks:

```bash
gofmt -l .    # no output means clean
go vet ./...
```

### What the tests check

`internal/calculator/calculator_test.go` - table-driven tests for every
operation (add, subtract, multiply, divide, power, sqrt, percentage) and
their error cases (division by zero, square root of a negative number,
unsupported operation). This package reaches 100% coverage because it is
pure logic with no external dependencies to mock.

`internal/handlers/calculate_test.go` - uses `net/http/httptest` to send
requests directly to the handler functions, without a running server. Covers
successful calculations, validation failures (missing fields, malformed
JSON, unknown fields), and how domain errors from the calculator map to HTTP
status codes.

`internal/httpx/json_test.go` - the small JSON response helpers used by
every handler.

`internal/server/router_test.go` - builds the full router with
`NewRouter(...)` and sends requests through it end to end, including CORS
headers and preflight (`OPTIONS`) requests.

### Coverage by package

| Package | Coverage |
|---------|----------|
| `internal/calculator` | 100% |
| `internal/server` | 100% |
| `internal/httpx` | 100% |
| `internal/handlers` | 92.2% |
| `cmd/server` | 0% (entrypoint wiring only, no branching logic) |

The gap in `internal/handlers` is intentional. A few lines guard against
`NaN`/`Infinity` operands, but standard JSON has no literal for either value,
so `encoding/json` rejects such a payload before those lines can execute.
They stay in the code as a safety net for future changes, documented in
place, rather than covered by a test that could not happen through the
actual API.
