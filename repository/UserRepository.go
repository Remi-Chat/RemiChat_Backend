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

// CreateUser inserts a new user into the database
func CreateUser(ctx context.Context, user models.User) (primitive.ObjectID, error) {
	result, err := db.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	log.Printf("User created with ID: %v", result.InsertedID)
	return result.InsertedID.(primitive.ObjectID), nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := db.UserCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email ID
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := db.UserCollection.FindOne(ctx, bson.M{"email_id": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(ctx context.Context, id primitive.ObjectID, updates bson.M) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := db.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("no user found with the given ID")
	}

	log.Printf("User updated with ID: %v", id)
	return nil
}

// DeleteUser deletes a user by their ID
func DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	result, err := db.UserCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no user found with the given ID")
	}

	log.Printf("User deleted with ID: %v", id)
	return nil
}
