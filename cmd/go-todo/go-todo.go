package main

import (
	"log"
	"os"

	"github.com/tiagogauziski/go-todo/pkg"
	"github.com/tiagogauziski/go-todo/pkg/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env", ".env.dev")

	if err != nil {
		log.Println("WARN: Unable to locate .env files.")
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URL"))
}

func main() {
	router := pkg.SetupRouter()

	router.Run()
}
