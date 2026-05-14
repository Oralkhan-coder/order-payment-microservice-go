package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Oralkhan-coder/notification-service/internal/consumer"
	"github.com/Oralkhan-coder/notification-service/internal/infrastructure"
	"github.com/Oralkhan-coder/notification-service/internal/provider"
	"github.com/Oralkhan-coder/notification-service/internal/provider/real"
	"github.com/Oralkhan-coder/notification-service/internal/provider/simulated"
	"github.com/Oralkhan-coder/notification-service/internal/repository"
	"github.com/joho/godotenv"
	goredis "github.com/redis/go-redis/v9"
)

func main() {
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Printf("warning: .env file not loaded: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	idempotencyTTL := getEnvDuration("IDEMPOTENCY_TTL", 24*time.Hour)
	providerMode := getEnv("PROVIDER_MODE", "SIMULATED")
	maxAttempts := getEnvInt("RETRY_MAX_ATTEMPTS", 4)

	rmq, err := infrastructure.NewRabbitMQ(rabbitURL)
	if err != nil {
		log.Fatalf("failed to initialize rabbitmq: %v", err)
	}
	defer func() {
		_ = rmq.Channel.Close()
		_ = rmq.Conn.Close()
	}()

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	defer redisClient.Close()

	store := repository.NewRedisIdempotencyStore(redisClient, idempotencyTTL)

	var sender provider.EmailSender
	if providerMode == "REAL" {
		sender = real.New(
			getEnv("SMTP_HOST", ""),
			getEnv("SMTP_PORT", "587"),
			getEnv("SMTP_USER", ""),
			getEnv("SMTP_PASS", ""),
			getEnv("SMTP_FROM", "noreply@example.com"),
		)
		log.Println("email provider: REAL (SMTP)")
	} else {
		sender = simulated.New()
		log.Println("email provider: SIMULATED")
	}

	c := consumer.NewNotificationConsumer(rmq.Channel, store, sender, maxAttempts)
	log.Println("notification consumer started")
	if err := c.Start(ctx); err != nil {
		log.Fatalf("notification consumer stopped with error: %v", err)
	}

	log.Println("notification service shutting down gracefully")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
