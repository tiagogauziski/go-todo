package main

import (
	"log"
	"os"

	"github.com/tiagogauziski/go-todo/pkg"
	"github.com/tiagogauziski/go-todo/pkg/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env", "../.env", "../../.env")

	if err != nil {
		log.Println("WARN: Unable to locate .env files.")
	}

	database.ConnectDatabase(os.Getenv("DATABASE_URI"))
}

func main() {
	router := pkg.SetupRouter()

	router.Run()
}
