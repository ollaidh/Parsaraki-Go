package main

import (
	"encoding/json"
	"log"
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

func loadConfig() Config {
	data, err := os.ReadFile("config.json")

	if err != nil {
		log.Println("Failed lo load config from file")
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal(data, &config)
	return config

}
