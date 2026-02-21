package database

import (
	"log"

	"go-api/backend-api/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB


func ConnectDB() *gorm.DB {
    dsn := config.GetEnv("DB_URL", "")

    // Check if the DSN is empty and log a fatal error if it is
    if dsn == "" {
        log.Fatal("Database URL is empty. Please check your .env file.")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        PrepareStmt: false,
    })
    if err != nil {
        panic(err)
    }
    
    DB = db // Assign to package-level variable
    return db
}