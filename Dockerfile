FROM golang:1.21-alpine AS builder
WORKDIR /app

# Copiar go.mod
COPY go.mod ./
RUN go mod download || true

# Copiar el resto del codigo fuente
COPY . .

# Compilar todos los comandos principales (cmd/main.go) dinámicamente
RUN mkdir -p /app/bin && \
    for cmd in $(find plugins -name "main.go" | grep "/cmd/"); do \
        plugin_name=$(echo "$cmd" | awk -F'/' '{print $(NF-2)}'); \
        echo "Compilando plugin: $plugin_name"; \
        go build -o /app/bin/$plugin_name "$cmd"; \
    done

FROM telegraf:alpine
# Copiar todos los binarios compilados
COPY --from=builder /app/bin/* /usr/local/bin/
