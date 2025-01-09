package main

import (
	"RemiAPI/db"
	"RemiAPI/handlers"
	"RemiAPI/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

const mongoURI = "mongodb://localhost:27017"
const dbName = "remi"

func initializeDatabase() (*mongo.Database, func(), error) {
	return db.ConnectToDB(mongoURI, dbName)
}

func setupApplication(client *mongo.Database) error {
	return db.ConfigureChannels(client)
}

func main() {
	fmt.Println("Starting the application...")

	client, cleanup, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Initialization failed: %v", err)
	}
	defer cleanup()

	// Other initialization tasks
	err = setupApplication(client)
	if err != nil {
		log.Fatalf("Application setup failed: %v", err)
	}

	log.Println("Application setup completed successfully.")

	userRepo := repository.NewUserRepository(client)
	authHandler := handlers.NewAuthHandler(userRepo)

	router := mux.NewRouter()
	router.HandleFunc("/signup", authHandler.SignupHandler).Methods("POST")
	router.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
