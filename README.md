# Sezzle Calculator

Full-stack calculator built for the Sezzle technical assessment. A React
and TypeScript frontend consumes a Go REST microservice that performs the
arithmetic.

```
├── backend/     Go service (net/http, no framework)
├── frontend/    React + TypeScript (Vite)
└── docker-compose.yaml
```

Supported operations: addition, subtraction, multiplication, division,
exponentiation, square root, percentage.

For details specific to each service, see [backend/README.md](backend/README.md)
and [frontend/README.md](frontend/README.md).

## Running with Docker Compose

This builds and starts both services together.

Requirements: Docker with the Compose plugin (`docker compose version` should
work).

```bash
docker compose up --build
```

- Frontend: http://localhost:3000
- Backend: http://localhost:8080

Stop with `docker compose down`, or add `-d` to `up` to run in the
background and use `docker compose logs -f` to follow logs.

`docker-compose.yaml` builds each service from its own `Dockerfile` and
wires them together:

```yaml
services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      ALLOWED_ORIGINS: "http://localhost:3000"

  frontend:
    build:
      context: ./frontend
      args:
        VITE_API_URL: "http://localhost:8080"
    ports:
      - "3000:80"
    depends_on:
      - backend
```

`ALLOWED_ORIGINS` tells the backend which origin is allowed to call it over
CORS; it is set to the frontend's compose address. `VITE_API_URL` is passed
as a build argument because Vite inlines environment variables into the
bundle at build time, not at container start time.

## Running without Docker

See the setup section of each service's README:

- [backend/README.md](backend/README.md#setup) - requires Go 1.22+
- [frontend/README.md](frontend/README.md#setup) - requires Node 20+

In short: run the backend first (`go run ./cmd/server`, listens on 8080),
then the frontend (`npm install && npm run dev`, serves on 5173).

## Tests

```bash
cd backend && go test ./... -cover
cd frontend && npm test
```

See each README for what the tests cover and how coverage is measured.

## API

Base URL: `http://localhost:8080`

| Method | Path | Purpose |
|--------|------|---------|
| GET | `/health` | Liveness probe |
| GET | `/api/operations` | List supported operations |
| POST | `/api/calculate` | Perform a calculation |

### POST /api/calculate

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `operation` | string | yes | One of the ids from `/api/operations` |
| `a` | number | yes | First operand |
| `b` | number | yes, except `sqrt` | Second operand |

```bash
curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":4,"b":2}'
# {"operation":"add","a":4,"b":2,"result":6}

curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"divide","a":10,"b":4}'
# {"operation":"divide","a":10,"b":4,"result":2.5}

curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"sqrt","a":81}'
# {"operation":"sqrt","a":81,"result":9}

curl -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"percentage","a":10,"b":200}'
# {"operation":"percentage","a":10,"b":200,"result":20}
```

Errors use a single envelope, `{"error": "..."}`:

```bash
curl -i -X POST http://localhost:8080/api/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"divide","a":10,"b":0}'
# HTTP/1.1 422 Unprocessable Entity
# {"error":"division by zero is not allowed"}
```

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Malformed JSON, missing/unknown fields, unsupported operation |
| 422 | Well-formed but mathematically invalid: division by zero, square root of a negative number, non-finite result |

## Design decisions

**Why this is a microservice.** The term describes scope, not quantity: one
independently deployable unit with a single responsibility. The backend does
arithmetic and nothing else, shares no state or code with the frontend, and
talks to it only over JSON/HTTP. It ships its own Dockerfile and exposes
`/health`, so it can be built, deployed and scaled independently.

**Go standard library, no web framework.** Since Go 1.22, `net/http.ServeMux`
routes by method and pattern (`"POST /api/calculate"`). For three endpoints a
third-party router adds a dependency without adding capability, so `go.mod`
has no external requirements.

**Arithmetic isolated from HTTP.** `internal/calculator` has no knowledge of
HTTP; it takes an operation and operands and returns a result or a domain
error. `internal/handlers` owns decoding, validation, and status codes. This
boundary is what makes the arithmetic testable without HTTP plumbing.

**400 versus 422.** Malformed input (bad JSON, missing `a`, unknown
operation) is a 400. A well-formed request that is mathematically impossible
- dividing by zero, the square root of a negative number - is a 422, because
the request was understood and rejected on its merits.

**Validation on both sides, not duplicated.** The frontend checks that
fields are present and numeric so it never sends a request that cannot
succeed. Mathematical rules live only on the backend; the UI shows whatever
message it returns.

**Percentage means "a% of b"**, e.g. `percentage(10, 200) = 20`, matching how
the `%` key behaves on a physical calculator. `sqrt` is the only unary
operation and omits `b` entirely rather than sending a meaningless zero.

## Assumptions

- Single-user, stateless service: no authentication, persistence, or
  calculation history.
- `ALLOWED_ORIGINS` defaults to `*` for ease of review; a real deployment
  would restrict it to known frontend hosts.
- Standard floating-point precision is used, as in any calculator built on
  IEEE-754 numbers. The frontend rounds displayed results to 12 significant
  digits so `0.1 + 0.2` reads as `0.3` instead of showing binary rounding
  noise.
