package main

import (
	"RemiAPI/db"
	"RemiAPI/middleware"
	"RemiAPI/routers"
	"RemiAPI/utils"
	"RemiAPI/ws"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var err = godotenv.Load()

var mongoURI = utils.GetEnv("MONGO_URI", "mongodb://localhost:27017")
var dbName = utils.GetEnv("DB_NAME", "remi")

func main() {

	// ================== INITIAL CONFIG ==================
	fmt.Println("Starting the application...")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Connecting to the database...", mongoURI, dbName)

	client, cleanup, err := db.ConnectToDB(mongoURI, dbName)

	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	defer cleanup()

	err = db.ConfigureChannels(client)

	if err != nil {
		log.Fatalf("Application setup failed: %v", err)
	}

	log.Println("Application setup completed successfully.")

	// ================== ROUTES ==================
	router := gin.Default()

	routers.AuthRoutes(router)
	routers.UserRoutes(router)

	router.GET("/ws", func(c *gin.Context) {
		// Apply authentication middleware
		authMiddleware := middleware.AuthMiddleware()
		authMiddleware(c)
		if c.IsAborted() {
			return
		}

		// Extract user data from context
		userID := c.GetString("user_id")

		// Pass user data to WebSocket handler
		ws.WebsocketHandler(c, userID)
	})

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Public route")
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}
