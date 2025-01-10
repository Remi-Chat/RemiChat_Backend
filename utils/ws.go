package utils

import (
	"RemiAPI/models"
	"encoding/json"

	"github.com/gorilla/websocket"
)

func SendMessage(conn *websocket.Conn, wsMessage *models.Message) error {
	message, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}

func ReadMessage(conn *websocket.Conn) (*models.Message, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var wsMessage models.Message
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return nil, err
	}

	return &wsMessage, nil
}
