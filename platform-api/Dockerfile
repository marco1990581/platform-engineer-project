# ==========================
# Build Stage
# ==========================
FROM golang:1.25.0 AS builder

WORKDIR /app

# Descargar dependencias primero para aprovechar la cache
COPY go.mod ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o platform-api ./cmd/api

# ==========================
# Runtime Stage
# ==========================
FROM debian:bookworm-slim

WORKDIR /app

# Certificados para HTTPS
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /app/config

COPY --from=builder /app/platform-api .

ENV PLATFORM_API_USERS_FILE=/app/config/users.json

EXPOSE 8081

ENTRYPOINT ["./platform-api"]
