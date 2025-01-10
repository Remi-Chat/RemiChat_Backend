package main

import (
	"RemiAPI/db"
	"RemiAPI/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var mongoURI = utils.GetEnv("MONGO_URI", "mongodb://localhost:27017")
var dbName = utils.GetEnv("DB_NAME", "remi")

func main() {
	fmt.Println("Starting the application...")

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

	router := gin.Default()

	log.Println("Server running on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}
