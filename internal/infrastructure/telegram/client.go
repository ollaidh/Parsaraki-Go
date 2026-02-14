package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"parsaraki-go/config"
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

type BotGetMeResponse struct {
	Ok     bool `json:"ok"`
	Result struct {
		Id                      int    `json:"id"`
		IsBot                   bool   `json:"is_bot"`
		FirstName               string `json:"first_name"`
		Username                string `json:"username"`
		CanJoinGroups           bool   `json:"can_join_groups"`
		CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
		SupportsInlineQueries   bool   `json:"supports_inline_queries"`
		CanConnectToBusiness    bool   `json:"can_connect_to_business"`
		HasMainWebApp           bool   `json:"has_main_web_app"`
		HasTopicsEnabled        bool   `json:"has_topics_enabled"`
	}
}

type TelegramClient struct {
	Config config.Config
}

func NewTelegramClient(config *config.Config) TelegramClient {
	tgClient := TelegramClient{
		Config: *config,
	}
	tgClient.pingBot()
	tgClient.setWebhook()

	return tgClient
}

// get bot url to use methods
func (tc *TelegramClient) getRequestUrl(method string) string {
	result, _ := url.JoinPath(
		tc.Config.TelegramBot.Url,
		"bot"+tc.Config.TelegramBot.Token,
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

	botGetMeResponse, err := parseCheckBotMessage(response)

	if err == nil {
		if strconv.Itoa(botGetMeResponse.Result.Id) != tc.Config.TelegramBot.Id {
			pingError = fmt.Sprintf("Incorrect Id=%s received in response, expected %s", botGetMeResponse.Result.Id, tc.Config.TelegramBot.Id)
		}

	} else {
		pingError = fmt.Sprintf("Failed to parse response body: %s Error: %s", response.Body, err)
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
		"url":          tc.Config.Webhooks.Url + tc.Config.Webhooks.Ep,
		"secret_token": tc.Config.Webhooks.Token,
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
		"sendMessage": &MessagePayloadGetter{},
		"sendPhoto":   &PhotoPayloadGetter{},
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
