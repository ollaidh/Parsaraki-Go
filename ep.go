package main

import (
	"fmt"
	"io"
	"net/http"
)

func ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != CONFIG.Webhooks.WebhooksSecretToken {
		http.Error(w, "Incorrect secret token in request header!", http.StatusBadRequest)
	} else {
		defer request.Body.Close()

		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(w, "Fail to read request body", http.StatusBadRequest)
			return
		}

		fmt.Println("Got request from Telegram Bot", string(body))

	}

}
