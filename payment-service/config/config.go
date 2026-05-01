package config

import (
	"os"
)

type Config struct {
	RabbitMQURL string
	HTTPPort    string
	GRPCPort    string
}

func InitConfig() *Config {

	return &Config{
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		HTTPPort:    getEnv("PAYMENT_HTTP_PORT", "8081"),
		GRPCPort:    getEnv("PAYMENT_GRPC_PORT", "9091"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
