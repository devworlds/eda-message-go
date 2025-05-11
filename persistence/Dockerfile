FROM golang:1.24 AS builder

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./persistence ./persistence

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o persistence-service ./persistence/cmd/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /app
COPY --from=builder /app/persistence-service .
EXPOSE 8080
CMD ["./persistence-service"]