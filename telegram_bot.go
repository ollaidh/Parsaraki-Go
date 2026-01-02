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

type BotMessage struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
			IsPremium    bool   `json:"is_premium"`
		} `json:"from"`
		Chat struct {
			ID        int64  `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int64  `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

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
		"url":          CONFIG.Webhooks.GatewayWebhooksUrl + CONFIG.Webhooks.GatewayWebhooksEp,
		"secret_token": CONFIG.Webhooks.WebhooksSecretToken,
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

func sendMessage(msg string, chatID int64) {
	url := getReqUrl("sendMessage")

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    msg,
	}

	body, _ := json.Marshal(payload)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Print(err)
	}

	rrr, _ := io.ReadAll(response.Body)

	fmt.Println(string(rrr))
}
