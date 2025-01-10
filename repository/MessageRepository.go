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

// NewMessageRepository initializes a new MessageRepository
func CreateMessage(ctx context.Context, message models.Message) (primitive.ObjectID, error) {
	result, err := db.MessageCollection.InsertOne(ctx, message)
	if err != nil {
		return primitive.NilObjectID, err
	}
	log.Printf("Message created with ID: %v", result.InsertedID)
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetMessageByID retrieves a message by its ID
func GetMessageByID(ctx context.Context, id primitive.ObjectID) (*models.Message, error) {
	var message models.Message
	err := db.MessageCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateMessage updates an existing message in the database
func UpdateMessage(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := db.MessageCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no message found with the given ID")
	}

	log.Printf("Message updated with ID: %v", id)
	return nil
}

func SaveMessage(ctx context.Context, message *models.Message) (primitive.ObjectID, error) {
	result, err := db.MessageCollection.InsertOne(ctx, *message)
	if err != nil {
		return primitive.NilObjectID, err
	}
	log.Printf("Message saved with ID: %v", result.InsertedID)
	return result.InsertedID.(primitive.ObjectID), nil
}

// DeleteMessage deletes a message by its ID
func DeleteMessage(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := db.MessageCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no message found with the given ID")
	}

	log.Printf("Message deleted with ID: %v", id)
	return nil
}
