# Etapa 1: build do Go
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copia arquivos do backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend .

# Compila o binário
RUN go build -ldflags="-s -w" -o bin/main ./cmd/main.go

# Etapa 2: imagem final
FROM alpine:3.20

WORKDIR /app

# Copia binário
COPY --from=builder /app/bin/main .

# Define variáveis default (sobrescrevíveis via docker-compose)
ENV SERVER_PORT=8080
ENV MONGODB_URI=mongodb://localhost:27017
ENV MONGODB_DATABASE=home_server_hub
ENV DOCKER_HOST=unix:///var/run/docker.sock

EXPOSE 8080

CMD ["./main"]

