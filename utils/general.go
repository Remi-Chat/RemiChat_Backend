package utils

import (
	"os"
	"strings"
	"time"

	"math/rand"

	ung "github.com/dillonstreator/go-unique-name-generator"
	"github.com/google/uuid"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to generate a unique name
func GenerateUniqueName() string {
	generator := ung.NewUniqueNameGenerator()

	name := generator.Generate()
	// "{adjective}_{color}_{name}"

	parts := strings.Split(name, "_")
	for i := range parts {
		parts[i] = strings.Title(parts[i])
	}

	// Shuffle and select 2 random parts
	rand.Seed(time.Now().UnixNano()) // Seed random number generator
	rand.Shuffle(len(parts), func(i, j int) { parts[i], parts[j] = parts[j], parts[i] })

	selectedParts := parts[:2] // Take the first 2 elements after shuffle

	name = strings.Join(selectedParts, " ")

	return name
}

// generate uuid

func GenerateUUID() (string, error) {
	// Generate a new UUID
	u, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	// Return the UUID as a string
	return u.String(), nil
}

// Generate email by taking uuid as param
func GenerateRandomEmail(uuid string) string {
	return uuid + "-temp-user-email@emeraldchat.com"
}
