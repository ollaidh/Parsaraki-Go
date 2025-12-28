package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func getReqUrl(method string) string {
	result, _ := url.JoinPath(
		CONFIG.TelegramBot.Url,
		"bot"+CONFIG.TelegramBot.Token,
		method)
	return result
}

func pingBot() {
	url := getReqUrl("getMe")
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}

	pingError := ""

	if response.Status != "200 OK" {
		pingError = fmt.Sprintf("Response status = %s", response.Status)

	}

	body, err := io.ReadAll(response.Body)
	var data map[string]interface{}

	err = json.Unmarshal(body, &data)
	if err == nil {
		receivedId := strconv.Itoa(int(data["result"].(map[string]interface{})["id"].(float64)))
		if receivedId != CONFIG.TelegramBot.Id {
			pingError = fmt.Sprintf("Incorrect Id=%s received in response, expected %s", receivedId, CONFIG.TelegramBot.Id)
		}

	} else {
		bodyString := string(body)
		pingError = fmt.Sprintf("Failed to parse response body: %s Error: %s", bodyString, err)
	}

	if pingError != "" {
		log.Fatalf("Bot Check failed: %s", pingError)
	} else {
		log.Printf("Bot Check Successful!")
	}

}

func setWebhook() {
	url := getReqUrl("setWebhook")

	payload := map[string]interface{}{
		"url":          CONFIG.TelegramBot.WebhooksUrl,
		"secret_token": CONFIG.TelegramBot.WebhooksSecretToken,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	rrr, _ := io.ReadAll(resp.Body)

	fmt.Println(string(rrr))
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()

}
