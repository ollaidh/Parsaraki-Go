package main

type RequestHandler interface {
	Execute(BotMessage)
}

type NotACommandHandler struct{}

func (handler NotACommandHandler) Execute(botMsg BotMessage) {
	sendMessage("It's not a command, available commands are:\n/fooo, /bar, /buzzz", botMsg.Message.Chat.ID)
}

type StartHandler struct{}

func (handler StartHandler) Execute(botMsg BotMessage) {
	sendMessage("Started bot: HELLO, Good lookin!", botMsg.Message.Chat.ID)
}

type BotInfoHandler struct{}

func (handler BotInfoHandler) Execute(botMsg BotMessage) {
	sendMessage("Greetings from Parsaraki bot. I'm trained to do nothing so far!", botMsg.Message.Chat.ID)
}

type GetStatisticsHandler struct{}

func (handler GetStatisticsHandler) Execute(botMsg BotMessage) {
	sendMessage("Greetings from Parsaraki bot. I'm trained to do nothing so far!", botMsg.Message.Chat.ID)
}
