package controllers

import (
	"RemiAPI/models"
	"RemiAPI/repository"
	"RemiAPI/utils"
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConnectionPool manages active WebSocket connections
type ConnectionPool struct {
	connections map[string]map[*websocket.Conn]*models.User // channelID -> connections
	mutex       sync.RWMutex
}

// Create a global connection pool
var Pool = &ConnectionPool{
	connections: make(map[string]map[*websocket.Conn]*models.User),
}

// Message types
const (
	JoinChannel  = "join_channel"
	LeaveChannel = "leave_channel"
	SendMessage  = "send_message"
	UserJoined   = "user_joined"
	UserLeft     = "user_left"
	NewMessage   = "new_message"
	Error        = "error"
)

// Add connection to pool
func (pool *ConnectionPool) AddConnection(channelID string, conn *websocket.Conn, user *models.User) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if pool.connections[channelID] == nil {
		pool.connections[channelID] = make(map[*websocket.Conn]*models.User)
	}
	pool.connections[channelID][conn] = user
}

// Remove connection from pool
func (pool *ConnectionPool) RemoveConnection(channelID string, conn *websocket.Conn) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	if connections, exists := pool.connections[channelID]; exists {
		delete(connections, conn)
		if len(connections) == 0 {
			delete(pool.connections, channelID)
		}
	}
}

// Broadcast message to all connections in a channel except sender
func (pool *ConnectionPool) Broadcast(channelID string, message *models.Message, sender *websocket.Conn) {
	pool.mutex.RLock()
	defer pool.mutex.RUnlock()

	if connections, exists := pool.connections[channelID]; exists {
		for conn := range connections {
			if conn != sender {
				utils.SendMessage(conn, message)
			}
		}
	}
}

// Handle joining a channel
func HandleJoinChannel(conn *websocket.Conn, msg *models.Message, user *models.User) {
	var channelData struct {
		ChannelID string `json:"channel_id"`
	}
	if err := json.Unmarshal([]byte(msg.Content), &channelData); err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Invalid channel data"})
		return
	}

	channelObjID, err := primitive.ObjectIDFromHex(channelData.ChannelID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Invalid channel ID"})
		return
	}

	// Get channel from database
	channel, err := repository.GetChannelByID(context.Background(), channelObjID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Channel not found"})
		return
	}

	// Check if channel is full
	if channel.MaxUserCount > 0 && len(channel.ActiveMembers) >= channel.MaxUserCount {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Channel is full"})
		return
	}

	// Add user to channel's active members
	err = repository.AddUserToChannel(context.Background(), channelObjID, user.ID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Failed to join channel"})
		return
	}

	// Add connection to pool
	Pool.AddConnection(channelData.ChannelID, conn, user)

	// Notify other users
	userJoinedMsg := &models.Message{
		Type:    UserJoined,
		Content: user.DisplayName + " joined the channel",
	}
	Pool.Broadcast(channelData.ChannelID, userJoinedMsg, conn)
}

// Handle leaving a channel
func HandleLeaveChannel(conn *websocket.Conn, msg *models.Message, user *models.User) {
	var channelData struct {
		ChannelID string `json:"channel_id"`
	}
	if err := json.Unmarshal([]byte(msg.Content), &channelData); err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Invalid channel data"})
		return
	}

	channelObjID, err := primitive.ObjectIDFromHex(channelData.ChannelID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Invalid channel ID"})
		return
	}

	// Remove user from channel's active members
	err = repository.RemoveUserFromChannel(context.Background(), channelObjID, user.ID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Failed to leave channel"})
		return
	}

	// Remove connection from pool
	Pool.RemoveConnection(channelData.ChannelID, conn)

	// Notify other users
	userLeftMsg := &models.Message{
		Type:    UserLeft,
		Content: user.DisplayName + " left the channel",
	}
	Pool.Broadcast(channelData.ChannelID, userLeftMsg, conn)
}

// Handle sending a message
func HandleSendMessage(conn *websocket.Conn, msg *models.Message, user *models.User) {
	var messageData struct {
		ChannelID string `json:"channel_id"`
		Content   string `json:"content,omitempty"`
		Image     string `json:"image,omitempty"`
		GIF       string `json:"gif,omitempty"`
	}
	if err := json.Unmarshal([]byte(msg.Content), &messageData); err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Invalid message data"})
		return
	}

	channelObjID, err := primitive.ObjectIDFromHex(messageData.ChannelID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Invalid channel ID"})
		return
	}

	// Create new message
	newMessage := models.Message{
		UserID:    user.ID,
		Content:   messageData.Content,
		Image:     messageData.Image,
		GIF:       messageData.GIF,
		Timestamp: time.Now().Unix(),
	}

	// Save message to database
	messageID, err := repository.SaveMessage(context.Background(), &newMessage)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Failed to save message"})
		return
	}

	// Add message to channel
	err = repository.AddMessageToChannel(context.Background(), channelObjID, messageID)
	if err != nil {
		utils.SendMessage(conn, &models.Message{Type: Error, Content: "Failed to add message to channel"})
		return
	}

	// Broadcast message to all users in channel
	messageResponse := struct {
		MessageID   string `json:"message_id"`
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
		Content     string `json:"content,omitempty"`
		Image       string `json:"image,omitempty"`
		GIF         string `json:"gif,omitempty"`
		Timestamp   int64  `json:"timestamp"`
	}{
		MessageID:   messageID.Hex(),
		UserID:      user.ID.Hex(),
		DisplayName: user.DisplayName,
		Content:     newMessage.Content,
		Image:       newMessage.Image,
		GIF:         newMessage.GIF,
		Timestamp:   newMessage.Timestamp,
	}

	responseBytes, _ := json.Marshal(messageResponse)
	broadcastMsg := &models.Message{
		Type:    NewMessage,
		Content: string(responseBytes),
	}
	Pool.Broadcast(messageData.ChannelID, broadcastMsg, nil)
}
