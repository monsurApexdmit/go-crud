package main

import (
	"go-crud/database"
	"go-crud/routes"
	"log"
)


func main() {
    database.Connect()
    database.Migrate() 
    r := routes.RegisterRoutes()

    log.Println("API running on :8004")
    log.Println("🔥 Hot reload is working!") // Add this line

    r.Run(":8004")
}
