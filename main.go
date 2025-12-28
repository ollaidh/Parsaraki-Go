package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	url, err := url.JoinPath(
		config.TelegramBot.Url,
		"bot"+config.TelegramBot.Token,
		"getMe")
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}

	body, err := io.ReadAll(response.Body)
	bodyString := string(body)
	fmt.Printf("Status: %s, Body: %s", response.Status, bodyString)

}
