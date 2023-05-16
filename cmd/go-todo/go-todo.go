package main

import (
	"log"
	"os"

	internal "github.com/tiagogauziski/go-todo/pkg"
	"github.com/tiagogauziski/go-todo/pkg/database"

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
	router := internal.SetupRouter()

	router.Run()
}
