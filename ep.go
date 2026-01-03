package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// ep registered at webhook service
func ProcessBotMessage(w http.ResponseWriter, request *http.Request) {
	telegramBotApiSecretToken := request.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	if telegramBotApiSecretToken != CONFIG.Webhooks.WebhooksSecretToken {
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
