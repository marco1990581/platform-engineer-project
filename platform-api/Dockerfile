# ==========================
# Build Stage
# ==========================
FROM golang:1.24 AS builder

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
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/platform-api .

EXPOSE 8081

ENTRYPOINT ["./platform-api"]
