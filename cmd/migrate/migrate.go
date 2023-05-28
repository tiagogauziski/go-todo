package main

import (
	"log"
	"os"

	"github.com/tiagogauziski/go-todo/internal/database"
	"github.com/tiagogauziski/go-todo/internal/models"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env", "../.env", ".env")

	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URI"))
}

func main() {
	err := database.Database.AutoMigrate(&models.Todo{})

	if err != nil {
		log.Fatal("Failed to run database migrations.")
	}
}
