package telegramemulator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

type From struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsPremium    bool   `json:"is_premium"`
}

type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type TelegramMessage struct {
	MessageID int    `json:"message_id"`
	From      From   `json:"from"`
	Chat      Chat   `json:"chat"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
}

type TelegramMessageEmulated struct {
	UpdateID int             `json:"update_id"`
	Message  TelegramMessage `json:"message"`
}

type TelegramMessagesEmulator struct{}

func (tme *TelegramMessagesEmulator) CreateMessage() TelegramMessageEmulated {
	currText := fmt.Sprintf("Fake Telegram Message: %s", time.Now())
	return TelegramMessageEmulated{
		UpdateID: rand.IntN(1_000_000),
		Message: TelegramMessage{
			MessageID: rand.IntN(1000),
			From: From{
				ID:           1361646730,
				IsBot:        false,
				FirstName:    "Maria",
				LastName:     "Lineva",
				Username:     "Maria_Lineva",
				LanguageCode: "ru",
				IsPremium:    true,
			},
			Chat: Chat{
				ID:        1361646730,
				FirstName: "Maria",
				LastName:  "Lineva",
				Username:  "Maria_Lineva",
				Type:      "private",
			},
			Date: time.Now().Unix(),
			Text: currText,
		},
	}
}

func (tme *TelegramMessagesEmulator) SendMessages(ctx context.Context, url string) {
	for {
		msg := tme.CreateMessage()
		jsonData, _ := json.Marshal(msg)

		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Telegram-Bot-Api-Secret-Token", "fwKklvjtDqHHP44KwZOMtoJXvbbf7tVt6fxJqjfwJ5mZBBvXqlv2rR19G5WdFT3w")

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		select {
		case <-ctx.Done():
			log.Printf("Stop sending Telegram Emulator Messages")
			return
		default:
		}
		time.Sleep(5 * time.Second)
	}
}

func NewTelegramMessagesEmulator() TelegramMessagesEmulator {

	tgMessageEmulator := TelegramMessagesEmulator{}
	return tgMessageEmulator
}
