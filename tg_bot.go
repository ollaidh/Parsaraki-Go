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

// message from bot structure - received by webhook update
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
		Date     int64  `json:"date"`
		Text     string `json:"text"`
		Entities []struct {
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			Type   string `json:"type"`
		} `json:"entities"`
	} `json:"message"`
}

type TelegramClient struct {
	config Config
}

func NewTelegramClient(config Config) TelegramClient {
	return TelegramClient{
		config: config,
	}
}

// get bot url to use methods
func (tc *TelegramClient) getRequestUrl(method string) string {
	result, _ := url.JoinPath(
		tc.config.TelegramBot.Url,
		"bot"+tc.config.TelegramBot.Token,
		method)
	return result
}

// check if bot is functioning
func (tc *TelegramClient) pingBot() {
	url := tc.getRequestUrl("getMe")
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
		if receivedId != tc.config.TelegramBot.Id {
			pingError = fmt.Sprintf("Incorrect Id=%s received in response, expected %s", receivedId, tc.config.TelegramBot.Id)
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

// set bot updates webhook
func (tc *TelegramClient) setWebhook() {
	url := tc.getRequestUrl("setWebhook")

	payload := map[string]interface{}{
		"url":          tc.config.Webhooks.GatewayWebhooksUrl + tc.config.Webhooks.GatewayWebhooksEp,
		"secret_token": tc.config.Webhooks.WebhooksSecretToken,
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

// send message from bot to chat
func (tc *TelegramClient) sendMessage(msg string, chatID int64) {
	url := tc.getRequestUrl("sendMessage")

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

// ep registered at webhook service
func (tc *TelegramClient) ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != tc.config.Webhooks.WebhooksSecretToken {
		http.Error(w, "Incorrect secret token in request header!", http.StatusBadRequest)
	} else {
		println("GOT REQUEST")
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)

		var botMsg BotMessage
		err = json.Unmarshal(body, &botMsg)

		if err != nil {
			http.Error(w, "Fail to read request body", http.StatusBadRequest)
			return
		}

		commandHandlers := map[string]RequestHandler{
			"/start":      StartHandler{},
			"/info":       BotInfoHandler{},
			"/statistics": GetStatisticsHandler{},
			"other":       NotACommandHandler{},
		}

		entities := botMsg.Message.Entities

		var cmd string

		if len(entities) > 0 {
			cmd = botMsg.Message.Text
		} else {
			cmd = "others"
		}

		commandHandler, ok := commandHandlers[cmd]
		if !ok {
			commandHandler = NotACommandHandler{}
		}

		commandHandler.Execute(botMsg)

	}

}
