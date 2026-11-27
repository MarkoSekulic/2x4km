package config

import "os"

type Config struct {
	Address                     string
	TemplateServiceAddress      string
	UsersServiceAddress         string
	NotificationsServiceAddress string
	FilesServiceAddress         string
}

func GetConfig() Config {
	return Config{
		Address:                     os.Getenv("GATEWAY_ADDRESS"),
		TemplateServiceAddress:      os.Getenv("TEMPLATE_SERVICE_ADDRESS"),
		UsersServiceAddress:         os.Getenv("USERS_SERVICE_ADDRESS"),
		NotificationsServiceAddress: os.Getenv("NOTIFICATIONS_SERVICE_ADDRESS"),
		FilesServiceAddress:         os.Getenv("FILES_SERVICE_ADDRESS"),
	}
}
