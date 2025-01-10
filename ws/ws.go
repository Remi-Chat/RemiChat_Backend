package ws

import (
	"RemiAPI/controllers"
	"RemiAPI/repository"
	"RemiAPI/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins, modify this as needed
		return true
	},
}

func WebsocketHandler(c *gin.Context, userID string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Printf("Invalid userID: %v", err)
		return
	}

	user, err := repository.GetUserByID(context.Background(), objectID)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return
	}

	log.Printf("User %s connected", user.DisplayName)

	for {
		receivedMessage, err := utils.ReadMessage(conn)
		if err != nil {
			log.Printf("Error reading WebSocket message: %v", err)
			break
		}

		switch receivedMessage.Type {
		case controllers.JoinChannel:
			controllers.HandleJoinChannel(conn, receivedMessage, user)
		case controllers.LeaveChannel:
			controllers.HandleLeaveChannel(conn, receivedMessage, user)
		case controllers.SendMessage:
			controllers.HandleSendMessage(conn, receivedMessage, user)
		default:
			return
		}
	}
}
