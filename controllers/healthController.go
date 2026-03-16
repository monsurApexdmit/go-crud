package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck - API health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"message":   "API is running",
	})
}

// CORSTest - Test CORS configuration
func CORSTest(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	method := c.Request.Method

	c.JSON(http.StatusOK, gin.H{
		"message": "CORS is working correctly",
		"cors_info": gin.H{
			"origin":             origin,
			"method":             method,
			"allowed_origins":    os.Getenv("ALLOWED_ORIGINS"),
			"credentials":        "true",
			"max_age":            "43200",
			"allowed_methods":    "GET, POST, PUT, PATCH, DELETE, OPTIONS",
			"allowed_headers":    c.Writer.Header().Get("Access-Control-Allow-Headers"),
			"exposed_headers":    c.Writer.Header().Get("Access-Control-Expose-Headers"),
			"preflight_request":  method == "OPTIONS",
		},
		"headers_received": c.Request.Header,
		"timestamp":        time.Now().Unix(),
	})
}

// APIInfo - API information endpoint
func APIInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api_name":    "Go CRUD E-Commerce API",
		"version":     "1.0.0",
		"environment": os.Getenv("ENV"),
		"features": gin.H{
			"authentication":   true,
			"authorization":    true,
			"pagination":       true,
			"search":           true,
			"file_upload":      true,
			"coupon_system":    true,
			"multi_tenant":     false, // Future feature
		},
		"endpoints": gin.H{
			"health":     "/health",
			"cors_test":  "/cors-test",
			"api_info":   "/api/info",
			"docs":       "/api/docs", // Future: API documentation
		},
		"timestamp": time.Now().Unix(),
	})
}
