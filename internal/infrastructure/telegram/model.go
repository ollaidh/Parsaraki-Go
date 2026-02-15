package telegram

// ADD move PayloadGetter to separate file
type PayloadGetter interface {
	GetPayload(string, int64) map[string]interface{}
}

// ADD maybe each realization to separate file
type MessagePayloadGetter struct {
}

// send message from bot to chat
func (pg *MessagePayloadGetter) GetPayload(msg string, chatID int64) map[string]interface{} {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    msg,
	}

	return payload
}

type PhotoPayloadGetter struct {
}

// send picture from bot to chat
func (pg *PhotoPayloadGetter) GetPayload(photo string, chatID int64) map[string]interface{} {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"photo":   photo,
	}

	return payload
}
