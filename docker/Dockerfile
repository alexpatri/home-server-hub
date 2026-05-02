# Etapa 1: build do Go
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copia arquivos do backend
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Compila o binário (Go puro, sem CGO)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/main ./cmd/main.go

# Etapa 2: imagem final
FROM alpine:3.20

WORKDIR /app

# Diretório padrão para banco e imagens
RUN mkdir -p /app/data/images

# Copia binário
COPY --from=builder /app/bin/main .

# Define variáveis default (sobrescrevíveis via docker-compose)
ENV SERVER_PORT=8000
ENV SQLITE_PATH=/app/data/home_server_hub.db
ENV IMAGES_DIR=/app/data/images
ENV DOCKER_HOST=unix:///var/run/docker.sock

EXPOSE 8000

CMD ["./main"]
