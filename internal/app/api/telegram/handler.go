package telegramapi

import (
	"net/http"
	"parsaraki-go/internal/infrastructure/kafka"
	"parsaraki-go/internal/infrastructure/telegram"
)

type TelegramHandler struct {
	TgClient      telegram.TelegramClient
	KafkaProducer kafka.KafkaProducer
}

// ep registered at webhook service
func (th *TelegramHandler) ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != th.TgClient.Config.Webhooks.Token {
		http.Error(w, "Incorrect secret token in request header!", http.StatusBadRequest)
	} else {
		println("GOT REQUEST")

		botMsg, err := telegram.ParseBotMessage(request)

		if err != nil {
			http.Error(w, "Fail to read request body", http.StatusBadRequest)
			println("Failed to process bot message")
			return
		}

		println(botMsg.Message.Text)

		// PARSE TELEGRAM MESSAGE USING FUNCTIONS FROM client.go
		// SEND TELEGRAM MESSAGE TO KAFKA USING FUNCTIONS FROM producer.go

		// action := "sendMessage"
		// content := botMsg.Message.Text

		// msgResponseToBot := th.TgClient.msgProcessor.ProcessMsg(content)

		// tc.sendContent(botMsg.Message.Chat.ID, action, msgResponseToBot)

	}

}
