package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type EmailConfig struct {
	Host         string `mapstructure:"SMTP_HOST"`
	Port         int    `mapstructure:"SMTP_PORT"`
	AuthEmail    string `mapstructure:"SMTP_AUTH_EMAIL"`
	AuthPassword string `mapstructure:"SMTP_AUTH_PASSWORD"`
}

func NewEmailConfig() EmailConfig {
	return EmailConfig{
		Host:         os.Getenv("SMTP_HOST"),
		Port:         587,
		AuthEmail:    os.Getenv("SMTP_AUTH_EMAIL"),
		AuthPassword: os.Getenv("SMTP_AUTH_PASSWORD"),
	}
}
