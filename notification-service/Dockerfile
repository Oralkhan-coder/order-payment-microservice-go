FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/notification-service ./cmd/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app

COPY --from=builder /bin/notification-service /app/notification-service

ENTRYPOINT ["/app/notification-service"]
