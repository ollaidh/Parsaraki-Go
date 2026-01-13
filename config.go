package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Mode        string `json:"mode"`
	TelegramBot struct {
		Url                 string `json:"url"`
		Token               string `json:"token"`
		Id                  string `json:"id"`
		WebhooksUrl         string `json:"webhooksUrl"`
		WebhooksEp          string `json:"webhooksEp"`
		WebhooksSecretToken string `json:"webhooksSecretToken"`
	} `json:"telegramBot"`
	Gateway struct {
		Port string `json:"port"`
	} `json:"gateway"`
	Webhooks struct {
		GatewayWebhooksUrl  string `json:"gatewayWebhooksUrl"`
		GatewayWebhooksEp   string `json:"gatewayWebhooksEp"`
		WebhooksSecretToken string `json:"webhooksSecretToken"`
	} `json:"webhooks"`
}

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.json")

	if err != nil {
		return Config{}, err
	}
	var config Config
	json.Unmarshal(data, &config)
	return config, nil

}
