package messages

import (
	"bisgo/app/models"
	"encoding/json"
	"errors"
)

// file reserved for messages type

func CreateMessagePayload(messageID int, method string, data any) ([]byte, error) {
	switch method {
	case "method1":
		tx, ok := data.(models.Transaction)

		if !ok {
			return nil, errors.New("error: can't create message payload")
		}

		return json.Marshal(tx)
	case "method2":
		return nil, errors.New("error: can't create message payload")
	default:
		return nil, errors.New("error: can't create message payload")
	}
}
