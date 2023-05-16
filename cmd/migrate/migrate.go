package main

import (
	"log"
	"os"

	"github.com/tiagogauziski/go-todo/pkg/database"
	"github.com/tiagogauziski/go-todo/pkg/models"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URL"))
}

func main() {
	err := database.Database.AutoMigrate(&models.Todo{})

	if err != nil {
		log.Fatal("Failed to run database migrations.")
	}
}
