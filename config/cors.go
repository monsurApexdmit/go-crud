package config

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig returns CORS middleware configuration
func CORSConfig() gin.HandlerFunc {
	// Get allowed origins from environment variable
	allowedOrigins := getAllowedOrigins()

	config := cors.Config{
		// AllowOrigins: List of allowed origins
		AllowOrigins: allowedOrigins,

		// AllowMethods: HTTP methods allowed
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		// AllowHeaders: Headers that can be used in the actual request
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
			"X-Tenant-Key",      // For multi-tenant support
			"X-API-Key",         // For API key authentication
			"X-Client-Version",  // Client version tracking
		},

		// ExposeHeaders: Headers that browsers are allowed to access
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Content-Disposition",
			"X-Total-Count",     // For pagination
			"X-Page-Count",      // For pagination
			"X-Current-Page",    // For pagination
			"X-Rate-Limit",      // Rate limit info
			"X-Rate-Remaining",  // Remaining requests
		},

		// AllowCredentials: Allow cookies to be sent
		AllowCredentials: true,

		// MaxAge: How long the results of a preflight request can be cached
		MaxAge: 12 * time.Hour,
	}

	return cors.New(config)
}

// getAllowedOrigins returns list of allowed origins from environment
func getAllowedOrigins() []string {
	// Get from environment variable
	originsEnv := os.Getenv("ALLOWED_ORIGINS")

	// Default origins for development
	defaultOrigins := []string{
		// HTTP origins
		"http://localhost:3000",      // React default
		"http://localhost:3001",      // Alternative React/Next.js
		"http://localhost:4200",      // Angular
		"http://localhost:8080",      // Vue
		"http://localhost:5173",      // Vite
		"http://localhost:5174",      // Alternative Vite
		"http://127.0.0.1:3000",
		"http://127.0.0.1:5173",
		// HTTPS origins (for local development with SSL)
		"https://localhost:3000",
		"https://localhost:3001",
		"https://localhost:4200",
		"https://localhost:5173",
		"https://127.0.0.1:3000",
		"https://127.0.0.1:3001",
	}

	// If environment variable is set, use it
	if originsEnv != "" {
		origins := strings.Split(originsEnv, ",")
		// Trim spaces
		for i, origin := range origins {
			origins[i] = strings.TrimSpace(origin)
		}
		return origins
	}

	// Development mode: allow all origins (use with caution!)
	if os.Getenv("ENV") == "development" || os.Getenv("GIN_MODE") != "release" {
		return []string{"*"} // Allow all origins in development
	}

	// Production: use default origins
	return defaultOrigins
}

// CustomCORSMiddleware is a custom CORS implementation for more control
func CustomCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if isOriginAllowed(origin) {
			// Set CORS headers
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-Key, X-API-Key")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Total-Count, X-Page-Count, X-Current-Page")
			c.Writer.Header().Set("Access-Control-Max-Age", "43200") // 12 hours
		}

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list
func isOriginAllowed(origin string) bool {
	allowedOrigins := getAllowedOrigins()

	// If wildcard is set, allow all
	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
	}

	return false
}
