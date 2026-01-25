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
		"url":          tc.config.Webhooks.Url + tc.config.Webhooks.Ep,
		"secret_token": tc.config.Webhooks.Token,
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

func (tc *TelegramClient) sendContent(chatID int64, actionType string, content string) (string, error) {
	url := tc.getRequestUrl(actionType) // ex actionType = "sendMessage"

	payloadGetters := map[string]PayloadGetter{
		"sendMessage": MessagePayloadGetter{},
		"sendPhoto":   PhotoPayloadGetter{},
	}

	payloadGetter, ok := payloadGetters[actionType]

	if !ok {
		return "", fmt.Errorf("Unknown action type: %v", actionType)
	}

	payload := payloadGetter.GetPayload(content, chatID)

	body, _ := json.Marshal(payload)

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Print(err)
	}

	resp, _ := io.ReadAll(response.Body)
	fmt.Println(string(resp))
	return string(resp), nil

}

// ep registered at webhook service
func (tc *TelegramClient) ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != tc.config.Webhooks.Token {
		http.Error(w, "Incorrect secret token in request header!", http.StatusBadRequest)
	} else {
		println("GOT REQUEST")

		botMsg, err := parseBotMessage(request)

		if err != nil {
			http.Error(w, "Fail to read request body", http.StatusBadRequest)
			println("Failed to process bot message")
			return
		}

		//cmd, err := getCommand(botMsg)

		//.... выбор команды и соответственно отправки

		action := "sendMessage"
		content := "testMessage"

		kc := NewKafkaClient("test-topic")

		kc.writeMessage("test-topic", content)
		kc.readMessage("test-topic")

		tc.sendContent(botMsg.Message.Chat.ID, action, content)

	}

}
