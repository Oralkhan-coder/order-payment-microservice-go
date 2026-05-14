package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	PaymentServiceHost string
	PaymentServicePort string
	RedisHost          string
	RedisPort          string
	CacheTTL           time.Duration
	RateLimitRequests  int
	RateLimitWindow    time.Duration
}

func InitConfig() *Config {
	return &Config{
		PaymentServiceHost: getEnv("PAYMENT_SERVICE_HOST", "localhost"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", "9091"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		CacheTTL:           getEnvDuration("CACHE_TTL", 5*time.Minute),
		RateLimitRequests:  getEnvInt("RATE_LIMIT_REQUESTS", 10),
		RateLimitWindow:    getEnvDuration("RATE_LIMIT_WINDOW", time.Minute),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}
