FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./auth/cmd/main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /app
COPY --from=builder /app/auth-service .
EXPOSE 8081
CMD ["./auth-service"]