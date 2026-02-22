package repository

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
	fmt.Printf("Saved to DB")
	md.messages = append(md.messages, botMsg)
	return nil
}
