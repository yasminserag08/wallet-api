package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/gorm"

	"wallet-api/models"

	"github.com/joho/godotenv" // to avoid hardcoding the password
	"gorm.io/driver/postgres"
)

func Connect() (*gorm.DB, error) {
	// in case there's no .env file
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment variables")
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{}); err != nil {
		return nil, err
	}

	return db, nil
}
