package data

import (
	"log"
	"fmt"
	"os"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using system environment variables")
    }

    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")

    if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" || dbPort == "" {
        log.Fatal("Missing one or more required database environment variables")
    }

    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
        dbHost, dbUser, dbPassword, dbName, dbPort,
    )

		var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Run migrations
    DB.AutoMigrate(&Set{}, &Card{}, &Listing{})
}
