package main

import (
	"encoding/json"
	"io"
	"net/http"
)

func parseBotMessage(request *http.Request) (BotMessage, error) {
	defer request.Body.Close()
	var botMsg BotMessage

	body, _ := io.ReadAll(request.Body)
	err := json.Unmarshal(body, &botMsg)
	if err != nil {
		return BotMessage{}, err
	}
	return botMsg, nil
}

func parseCheckBotMessage(request *http.Response) (BotGetMeResponse, error) {
	defer request.Body.Close()
	var botGetMeResponse BotGetMeResponse

	body, _ := io.ReadAll(request.Body)
	err := json.Unmarshal(body, &botGetMeResponse)
	if err != nil {
		return BotGetMeResponse{}, err
	}
	return botGetMeResponse, nil
}

func getCommand(botMsg *BotMessage) (string, error) {
	entities := botMsg.Message.Entities

	var cmd string
	if len(entities) > 0 {
		cmd = botMsg.Message.Text
	} else {
		cmd = "others"
	}

	return cmd, nil

}
