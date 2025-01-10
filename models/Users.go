package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	EmailID        string             `bson:"email_id" json:"email_id"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	Username       string             `bson:"username" json:"username"`
	DisplayPicture string             `bson:"display_picture,omitempty" json:"display_picture,omitempty"`
	PasswordHash   string             `bson:"password_hash" json:"password_hash"`
	Gender         string             `bson:"gender,omitempty" json:"gender,omitempty"`
}
