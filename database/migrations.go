package database

import (
    "log"
    "go-crud/models"
)

func Migrate() {
    err := DB.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatal("Migration failed:", err)
    }

    log.Println("Database migrated successfully.")
}
