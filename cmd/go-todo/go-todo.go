package main

import (
	"log"
	"os"

	"github.com/tiagogauziski/go-todo/internal"
	"github.com/tiagogauziski/go-todo/internal/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env", "../.env", ".env")

	if err != nil {
		log.Println("WARN: Unable to locate .env files.")
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URI"))
}

func main() {
	router := internal.SetupRouter()

	router.Run()
}
