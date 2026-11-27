package config

import (
	"os"
	"shared/smtp"
	"strconv"
)

type DBConfig struct {
	Address  string
	Username string
	Password string
	MinPool  uint64
	MaxPool  uint64
}

type NotificationConfig struct {
	JWT_PUBLIC_KEY string
}

type Config struct {
	Address string
	DBConfig
	NotificationConfig
	SMTPConfig smtp.Config
}

func getEnvInt(env string, defaultVal int) int {
	val, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		return defaultVal
	}
	return val
}

func GetConfig() Config {
	return Config{
		Address: os.Getenv("SERVICE_PORT"),
		DBConfig: DBConfig{
			Address: os.Getenv("CASSANDRA_ADDR"),
			MinPool: uint64(getEnvInt("DB_MIN_POOL", 0)),
			MaxPool: uint64(getEnvInt("DB_MAX_POOL", 1)),
		},
		NotificationConfig: NotificationConfig{
			JWT_PUBLIC_KEY: os.Getenv("JWT_PUBLIC_KEY"),
		},
		SMTPConfig: smtp.Config{
			Host:     "mailhog",
			Port:     1025,
			From:     "no-reply@localhost.test",
			Username: "",
			Password: "",
		},
	}
}
