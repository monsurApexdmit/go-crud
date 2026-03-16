package main

import (
	"go-crud/database"
	"go-crud/routes"
	"log"
	"os"
)

func main() {
	database.Connect()
	database.Migrate()
	r := routes.RegisterRoutes()

	// Check if SSL certificates exist for HTTPS
	certFile := "./certs/cert.pem"
	keyFile := "./certs/key.pem"

	// Check if running in HTTPS mode
	useHTTPS := os.Getenv("USE_HTTPS")

	if useHTTPS == "true" {
		// Check if certificate files exist
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Fatal("Certificate file not found at ./certs/cert.pem")
		}
		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			log.Fatal("Private key file not found at ./certs/key.pem")
		}

		port := os.Getenv("PORT")
		if port == "" {
			port = "8004"
		}

		log.Println("🔒 HTTPS Server starting on https://localhost:" + port)
		log.Println("🔥 Hot reload is working!")

		// Run with HTTPS
		if err := r.RunTLS(":"+port, certFile, keyFile); err != nil {
			log.Fatal("Failed to start HTTPS server: ", err)
		}
	} else {
		// Run with HTTP (default)
		port := os.Getenv("PORT")
		if port == "" {
			port = "8004"
		}

		log.Println("🌐 HTTP Server starting on http://localhost:" + port)
		log.Println("🔥 Hot reload is working!")

		if err := r.Run(":" + port); err != nil {
			log.Fatal("Failed to start HTTP server: ", err)
		}
	}
}
