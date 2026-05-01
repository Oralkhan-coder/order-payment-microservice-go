package config

import (
	"os"
	"strconv"
)

type Config struct {
	Db                 *PostgresConfig
	PaymentServiceHost string
	PaymentServicePort string
}

func InitConfig() *Config {
	dbCfg := PostgresConfig{
		Database: getEnv("POSTGRESQL_DB", "order_db"),
		Host:     getEnv("POSTGRESQL_URI", "localhost"),
		Port:     getEnvAsUint16("POSTGRESQL_PORT", 5432),
		Username: getEnv("POSTGRESQL_USERNAME", "postgres"),
		Password: getEnv("POSTGRESQL_PASSWORD", "postgres"),
	}

	return &Config{
		Db:                 &dbCfg,
		PaymentServiceHost: getEnv("PAYMENT_SERVICE_HOST", "localhost"),
		PaymentServicePort: getEnv("PAYMENT_SERVICE_PORT", "9091"),
	}
}

type PostgresConfig struct {
	Database string `env:"POSTGRESQL_DB"`
	Host     string `env:"POSTGRESQL_URI"`
	Port     uint16 `env:"POSTGRESQL_PORT"`
	Username string `env:"POSTGRESQL_USERNAME"`
	Password string `env:"POSTGRESQL_PASSWORD"`
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvAsUint16(key string, fallback uint16) uint16 {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	parsed, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return fallback
	}
	return uint16(parsed)
}
