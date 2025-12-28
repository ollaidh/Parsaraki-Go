package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	TelegramBot struct {
		Url   string `json:"url"`
		Token string `json:"token"`
	} `json:"telegramBot"`
}

func loadConfig() (Config, error) {
	data, err := os.ReadFile("config.json")

	if err != nil {
		log.Println("Failed lo load config from file", err)
		return Config{}, err
	}
	var config Config
	json.Unmarshal(data, &config)
	return config, nil

}
