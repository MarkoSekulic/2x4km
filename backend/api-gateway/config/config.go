package config

import "os"

type Config struct {
	Address                string
	TemplateServiceAddress string
	UsersServiceAddress    string
}

func GetConfig() Config {
	return Config{
		Address:                os.Getenv("GATEWAY_ADDRESS"),
		TemplateServiceAddress: os.Getenv("TEMPLATE_SERVICE_ADDRESS"),
		UsersServiceAddress:    os.Getenv("USERS_SERVICE_ADDRESS"),
	}
}
