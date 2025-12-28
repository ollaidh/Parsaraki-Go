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
		Id    string `json:"id"`
	} `json:"telegramBot"`
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
