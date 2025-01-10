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

func CreateChannel(ctx context.Context, channel models.Channel) (primitive.ObjectID, error) {
	result, err := db.ChannelCollection.InsertOne(ctx, channel)
	if err != nil {
		return primitive.NilObjectID, err
	}
	log.Printf("Channel created with ID: %v", result.InsertedID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func GetChannelByID(ctx context.Context, id primitive.ObjectID) (*models.Channel, error) {
	var channel models.Channel
	err := db.ChannelCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

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

func AddUserToChannel(ctx context.Context, channelID primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": channelID}
	update := bson.M{"$addToSet": bson.M{"users": userID}}

	result, err := db.ChannelCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no channel found with the given ID")
	}

	log.Printf("User added to channel with ID: %v", channelID)
	return nil
}

func RemoveUserFromChannel(ctx context.Context, channelID primitive.ObjectID, userID primitive.ObjectID) error {
	filter := bson.M{"_id": channelID}
	update := bson.M{"$pull": bson.M{"users": userID}}

	result, err := db.ChannelCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no channel found with the given ID")
	}

	log.Printf("User removed from channel with ID: %v", channelID)
	return nil
}

func AddMessageToChannel(ctx context.Context, channelID primitive.ObjectID, messageID primitive.ObjectID) error {
	filter := bson.M{"_id": channelID}
	update := bson.M{"$addToSet": bson.M{"messages": messageID}}

	result, err := db.ChannelCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no channel found with the given ID")
	}

	log.Printf("Message added to channel with ID: %v", channelID)
	return nil
}

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
