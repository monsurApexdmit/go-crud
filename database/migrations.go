package database

import (
    "log"
    "go-crud/models"
)

func Migrate() {
    err := DB.AutoMigrate(&models.User{}, &models.Author{}, &models.Book{})
    if err != nil {
        log.Fatal("Migration failed:", err)
    }

    log.Println("Database migrated successfully.")
}
