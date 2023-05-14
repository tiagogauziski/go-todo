package main

import (
	"log"
	"os"

	"github.com/tiagogauziski/go-todo/internal"
	"github.com/tiagogauziski/go-todo/internal/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env", "../../.env.dev")

	if err != nil {
		log.Fatal("Error loading .env file.")
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URL"))
}

func main() {
	internal.StartServer()
}
