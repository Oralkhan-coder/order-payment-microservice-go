FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/order-service ./cmd/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /bin/order-service /app/order-service
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080 9090
ENTRYPOINT ["/app/order-service"]
