package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Address string
	MinPool uint64
	MaxPool uint64
}

type Config struct {
	Address string
	DBConfig
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
	}
}
