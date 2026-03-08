package inmemoryrepo

import (
	"fmt"
	"parsaraki-go/internal/infrastructure/telegram"
)

type InMemoryRepo struct {
	messages []telegram.BotMessage
}

func NewMemoryDB() InMemoryRepo {
	return InMemoryRepo{}
}

func (md *InMemoryRepo) SaveBotRequest(botMsg telegram.BotMessage) error {
	md.messages = append(md.messages, botMsg)
	fmt.Printf("Message saved: %v\n", md.messages)
	return nil
}
