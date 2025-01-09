package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Channel struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	ChannelName   string               `bson:"channel_name" json:"channel_name"`
	ActiveMembers []primitive.ObjectID `bson:"active_members" json:"active_members"`
	Messages      []primitive.ObjectID `bson:"messages" json:"messages"`
	MaxUserCount  int                  `bson:"max_user_count,omitempty" json:"max_user_count,omitempty"`
}
