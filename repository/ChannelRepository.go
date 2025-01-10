package repository

import (
	"RemiAPI/db"
	"RemiAPI/models"
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewChannelRepository initializes a new ChannelRepository
func CreateChannel(ctx context.Context, channel models.Channel) (primitive.ObjectID, error) {
	result, err := db.ChannelCollection.InsertOne(ctx, channel)
	if err != nil {
		return primitive.NilObjectID, err
	}
	log.Printf("Channel created with ID: %v", result.InsertedID)
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetChannelByID retrieves a channel by its ID
func GetChannelByID(ctx context.Context, id primitive.ObjectID) (*models.Channel, error) {
	var channel models.Channel
	err := db.ChannelCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

// UpdateChannel updates an existing channel in the database
func UpdateChannel(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := db.ChannelCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no channel found with the given ID")
	}

	log.Printf("Channel updated with ID: %v", id)
	return nil
}

// DeleteChannel deletes a channel by its ID
func DeleteChannel(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := db.ChannelCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no channel found with the given ID")
	}

	log.Printf("Channel deleted with ID: %v", id)
	return nil
}
