package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id" json:"user_id"` // Reference to User
	Content   string             `bson:"content,omitempty" json:"content,omitempty"`
	Image     string             `bson:"image,omitempty" json:"image,omitempty"` // URL or Path to Image
	GIF       string             `bson:"gif,omitempty" json:"gif,omitempty"`     // URL or Path to GIF
	Timestamp int64              `bson:"timestamp" json:"timestamp"`             // Unix Timestamp
}
