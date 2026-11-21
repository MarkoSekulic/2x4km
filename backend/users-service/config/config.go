package config

import (
	"os"
	"smtp"
	"strconv"
)

type DBConfig struct {
	Address string
	MinPool uint64
	MaxPool uint64
}

type UserConfig struct {
	DefaultDataLimitMB uint64
	JWT_PRIVATE_KEY    string
	JWT_PUBLIC_KEY     string
}

type Config struct {
	Address string
	DBConfig
	UserConfig
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
			Address: os.Getenv("MONGODB_ADDR"),
			MinPool: uint64(getEnvInt("DB_MIN_POOL", 0)),
			MaxPool: uint64(getEnvInt("DB_MAX_POOL", 1)),
		},
		UserConfig: UserConfig{
			DefaultDataLimitMB: uint64(getEnvInt("USER_DEFAULT_DATA_LIMIT_MB", 10)),
			JWT_PRIVATE_KEY:    os.Getenv("JWT_PRIVATE_KEY"),
			JWT_PUBLIC_KEY:     os.Getenv("JWT_PUBLIC_KEY"),
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
