package database

import (
	"log"

	"github.com/tiagogauziski/go-todo/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDatabase(dsn string) {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	Database = db
}

func CreateDatabase() {
	if Database == nil {
		log.Fatal("Database not initialised.")
	}

	// TODO: create database if not exists
}

func RunMigrations() {
	err := Database.AutoMigrate(&models.Todo{})

	if err != nil {
		log.Fatal("Failed to run database migrations.")
	}
}
