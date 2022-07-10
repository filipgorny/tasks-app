package organizer

import (
	"os"

	"github.com/joho/godotenv"
)

type OrganizerServiceConfig struct {
	Protocol string
	Domain   string
	Port     string
	Token    string
}

func LoadConfig() OrganizerServiceConfig {
	godotenv.Load(".env")

	config := OrganizerServiceConfig{}

	config.Protocol = os.Getenv("ORG_SERVICE_PROTOCOL")
	config.Domain = os.Getenv("ORG_SERVICE_DOMAIN")
	config.Port = os.Getenv("ORG_SERVICE_PORT")
	config.Token = os.Getenv("ORG_SERVICE_TOKEN")

	return config
}
