package main

import (
	"log"
	"net/http"
	"go-crud/database"
	"go-crud/routes"
)

func main() {
    database.Connect()

    r := routes.RegisterRoutes()

    log.Println("API running on :8004")
    log.Println("🔥 Hot reload is working!") // Add this line

    http.ListenAndServe(":8004", r)
}
