package repository

import "parsaraki-go/internal/infrastructure/telegram"

type Repository interface {
	SaveBotRequest(botMsg telegram.BotMessage) error
}
