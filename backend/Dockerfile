# syntax=docker/dockerfile:1

# --- Build stage -----------------------------------------------------------
# Compiles a static binary. Using the full golang image here (not -alpine)
# avoids libc mismatches; CGO is disabled below so the binary has no
# dynamic dependencies anyway, keeping this stage's choice irrelevant to
# the final image size.
FROM golang:1.25 AS build

WORKDIR /src

# Cache module downloads separately from source changes.
COPY go.mod ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produces a fully static binary so it can run on the
# minimal distroless/scratch base below with no shared libraries.
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/calc-server ./cmd/server

# --- Test stage --------------------------------------------------------
# Optional intermediate stage: `docker build --target test` runs the full
# suite inside a clean container, useful for CI without installing Go
# locally.
FROM build AS test
RUN go test ./...

# --- Runtime stage -----------------------------------------------------
# distroless gives us a minimal, low-CVE-surface base with no shell and
# no package manager, appropriate for a small production microservice.
FROM gcr.io/distroless/static-debian12:nonroot AS runtime

WORKDIR /app
COPY --from=build /out/calc-server /app/calc-server

ENV PORT=8080
EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/calc-server"]
