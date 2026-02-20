package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Mode        string `env:"mode"`
	TelegramBot struct {
		Url   string `env:"url"`
		Token string `env:"token"`
		Id    string `env:"id"`
	} `env:"telegramBot"`
	Gateway struct {
		Port string `env:"port"`
	} `env:"gateway"`
	Webhooks struct {
		Url   string `env:"gatewayWebhooksUrl"`
		Ep    string `env:"gatewayWebhooksEp"`
		Token string `env:"webhooksSecretToken"`
	} `json:"webhooks"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Failed to load config", err)
		return Config{}, err
	}

	var cfg Config
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Println("Failed to parse config", err)
		return Config{}, err
	}

	return cfg, nil

}
