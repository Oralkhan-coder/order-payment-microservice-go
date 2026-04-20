package config

type Config struct {
	Db *PostgresConfig
}

func InitConfig() *Config {
	dbCfg := PostgresConfig{
		Database: "order_db",
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "postgres",
	}

	return &Config{
		Db: &dbCfg,
	}
}

type PostgresConfig struct {
	Database string `env:"POSTGRESQL_DB"`
	Host     string `env:"POSTGRESQL_URI"`
	Port     uint16 `env:"POSTGRESQL_PORT"`
	Username string `env:"POSTGRESQL_USERNAME"`
	Password string `env:"POSTGRESQL_PASSWORD"`
}
