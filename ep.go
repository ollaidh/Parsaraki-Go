package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != CONFIG.Webhooks.WebhooksSecretToken {
		http.Error(w, "Incorrect secret token in request header!", http.StatusBadRequest)
	} else {
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)

		var botMsg BotMessage
		err = json.Unmarshal(body, &botMsg)

		if err != nil {
			http.Error(w, "Fail to read request body", http.StatusBadRequest)
			return
		}

		msg := fmt.Sprintf("Got message from username=%s, message_id=%s, message: %s", botMsg.Message.From.Username, strconv.Itoa(botMsg.Message.MessageID), botMsg.Message.Text)

		fmt.Println(msg)
		fmt.Println(string(body))

		msgBack := fmt.Sprintf("Thank you for your message, you sent '%s'", botMsg.Message.Text)
		sendMessage(msgBack, botMsg.Message.Chat.ID)

	}

}
