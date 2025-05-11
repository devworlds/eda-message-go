FROM golang:1.24 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . ./

RUN apt-get update && apt-get install -y librdkafka-dev gcc libc-dev

RUN CGO_ENABLED=1 GOOS=linux go build -o persistence-service ./persistence/cmd/main.go

# Use uma imagem distroless baseada em Debian 12 (bookworm)
FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=builder /app/persistence-service ./

CMD ["./persistence-service"]