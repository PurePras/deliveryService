# --- build stage ---
FROM golang:1.25-alpine AS build
WORKDIR /src

# Cached separately from source so `go mod download` only reruns when deps change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 for a static binary that runs on the distroless image below (no libc).
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# --- runtime stage ---
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/api ./api

# Migrations are embedded into the binary (migrations/embed.go), so nothing else
# needs to be copied in — the image only ever holds the one binary.
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/api"]
