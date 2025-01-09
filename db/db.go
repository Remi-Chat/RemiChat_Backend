package db

import (
	"RemiAPI/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDB(uri, dbName string) (*mongo.Database, func(), error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}

	db := client.Database(dbName)
	return db, cleanup, nil
}

func ConfigureChannels(db *mongo.Database) error {
	channelCollection := db.Collection("channels")
	channels := []models.Channel{
		{ChannelName: "General Chat", MaxUserCount: 40},
		{ChannelName: "Roleplay Chat", MaxUserCount: 40},
		{ChannelName: "Regional Chat", MaxUserCount: 40},
		{ChannelName: "Tech Chat", MaxUserCount: 40},
		{ChannelName: "Bots :)", MaxUserCount: 40},
	}

	for _, channel := range channels {
		if err := upsertChannel(channelCollection, channel); err != nil {
			return err
		}
	}

	return nil
}

func upsertChannel(collection *mongo.Collection, channel models.Channel) error {
	filter := bson.M{"channel_name": channel.ChannelName}
	update := bson.M{"$set": channel}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	log.Printf("Configured or updated channel: %s", channel.ChannelName)
	return nil
}
