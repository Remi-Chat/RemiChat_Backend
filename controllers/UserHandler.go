package controllers

import (
	"RemiAPI/repository"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateProfileHandler handles updates to user profiles
func UpdateProfileHandler(c *gin.Context) {
	// Get the user ID from the context set by the AuthMiddleware
	userIDHex := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Parse the request body for updates
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// If password is being updated, hash it
	if password, ok := updates["password"]; ok {
		hash := sha256.New()
		hash.Write([]byte(password.(string)))
		updates["password_hash"] = hex.EncodeToString(hash.Sum(nil))
		delete(updates, "password") // Remove plain password key
	}

	// Perform the update
	err = repository.UpdateUser(context.Background(), userID, bson.M(updates))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetUserDetailsHandler(c *gin.Context) {
	// Get the user ID from the context set by the AuthMiddleware
	userIDHex := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve user details
	user, err := repository.GetUserByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user details
	c.JSON(http.StatusOK, gin.H{
		"user": map[string]interface{}{
			"email_id":        user.EmailID,
			"display_name":    user.DisplayName,
			"username":        user.Username,
			"gender":          user.Gender,
			"display_picture": user.DisplayPicture,
		},
	})
}
