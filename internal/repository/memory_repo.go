package repository

import (
	"encoding/json"
	"fmt"
)

type MemoryDB struct {
}

func NewMemoryDB() MemoryDB {
	return MemoryDB{}
}

func (md *MemoryDB) SaveBotRequest(botMsg json.RawMessage) error {
	fmt.Printf("Saved to DB: %s", botMsg)
	return nil
}
