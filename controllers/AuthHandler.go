package controllers

import (
	"RemiAPI/models"
	"RemiAPI/repository"
	"RemiAPI/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

var exptimeStr = utils.GetEnv("TOKEN_EXPIRES_IN", "3600")
var exptime, _ = strconv.Atoi(exptimeStr)

// Generate random user data
func generateRandomUserData(uuid string) models.User {
	return models.User{
		EmailID:     utils.GenerateRandomEmail(uuid),
		DisplayName: utils.GenerateUniqueName(),
	}
}

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

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token as a cookie
	c.SetCookie("session", token, exptime, "/", "", true, true)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// CreateTempUser creates a temporary user and returns a JWT token as a cookie
func CreateTempUser(c *gin.Context) {
	userUUID, _ := utils.GenerateUUID()
	newUser := generateRandomUserData(userUUID)
	newUser.Username, _ = utils.GenerateUUID()
	newUser.PasswordHash, _ = utils.GenerateUUID()
	newUser.Gender = "Others"
	newUser.DisplayPicture = "https://ui-avatars.com/api/?background=random?name=" + newUser.DisplayName

	// Hash the password
	hash := sha256.New()
	hash.Write([]byte(newUser.PasswordHash))
	newUser.PasswordHash = hex.EncodeToString(hash.Sum(nil))

	user, err := repository.CreateUser(context.Background(), newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set token as a cookie
	c.SetCookie("session", token, exptime, "/", "", true, true)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
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

	// Set token as a cookie
	c.SetCookie("session", token, exptime, "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
