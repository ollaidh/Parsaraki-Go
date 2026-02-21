package telegramapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"parsaraki-go/config"
	"parsaraki-go/internal/infrastructure/telegram"
)

type MessageBroker interface {
	WriteMessage(message json.RawMessage)
}

func NewTelegramHandler(msgBroker MessageBroker, config config.Config) *TelegramHandler {
	return &TelegramHandler{
		msgProducer: msgBroker,
		config:      config,
	}
}

type TelegramHandler struct {
	msgProducer MessageBroker
	config      config.Config
}

// ep registered at webhook service
func (th *TelegramHandler) ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != th.config.Webhooks.Token {
		http.Error(w, "Incorrect secret token in request header!", http.StatusBadRequest)
	} else {

		botMsg, err := telegram.ParseBotMessage(request)
		fmt.Printf("\nGOT REQUEST to /bot-message endpoint: %s", botMsg.Message.Text)

		if err != nil {
			http.Error(w, "Fail to read request body", http.StatusBadRequest)
			println("Failed to process bot message")
			return
		}

		msgForBroker, err := json.Marshal(botMsg)

		th.msgProducer.WriteMessage(msgForBroker)

	}

}
