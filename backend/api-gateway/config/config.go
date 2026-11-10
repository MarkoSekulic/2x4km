package config

import "os"

type Config struct {
	Address                string
	TemplateServiceAddress string
}

func GetConfig() Config {
	return Config{
		TemplateServiceAddress: os.Getenv("TEMPLATE_SERVICE_ADDRESS"),
		Address:                os.Getenv("GATEWAY_ADDRESS"),
	}
}
