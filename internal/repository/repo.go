package repository

import "encoding/json"

type Repository interface {
	SaveBotRequest(json.RawMessage) error
}
