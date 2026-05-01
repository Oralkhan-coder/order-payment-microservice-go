package config

import (
	"os"
)

type Config struct {
	PaymentServiceHost string
	PaymentServicePort string
}

func InitConfig() *Config {
	return &Config{
		PaymentServiceHost: getEnv("PAYMENT_SERVICE_HOST", "localhost"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", "9091"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
