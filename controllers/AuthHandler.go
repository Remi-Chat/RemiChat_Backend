package controllers

import (
	"RemiAPI/models"
	"RemiAPI/repository"
	"RemiAPI/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SignupHandler handles user registration
func SignupHandler(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hash := sha256.New()
	hash.Write([]byte(req.PasswordHash))
	req.PasswordHash = hex.EncodeToString(hash.Sum(nil))

	user, err := repository.CreateUser(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.Hex()})
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := repository.GetUserByEmail(context.Background(), req.ID.Hex())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify password
	hash := sha256.New()
	hash.Write([]byte(req.PasswordHash))
	if user.PasswordHash != hex.EncodeToString(hash.Sum(nil)) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
