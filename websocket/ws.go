package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins, modify this as needed
		return true
	},
}

// Read message from WebSocket
func readMessage(conn *websocket.Conn) (*WebSocketMessage, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	var wsMessage WebSocketMessage
	if err := json.Unmarshal(message, &wsMessage); err != nil {
		return nil, err
	}

	return &wsMessage, nil
}

// Send message to WebSocket
func sendMessage(conn *websocket.Conn, wsMessage *WebSocketMessage) error {
	message, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, message)
}

// WebSocket handler
func WebsocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	for {
		// Receive a message
		receivedMessage, err := readMessage(conn)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		log.Printf("Received message: %+v", receivedMessage)

		// Respond with a confirmation
		response := &WebSocketMessage{
			Type:    "response",
			Content: "Message received: " + receivedMessage.Content,
		}

		if err := sendMessage(conn, response); err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
			break
		}
	}
}
