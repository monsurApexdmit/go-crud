package controllers

import (
	"go-crud/database"
	"go-crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListLocations retrieves all locations
func ListLocations(c *gin.Context) {
	var locations []models.Location
	if err := database.DB.Find(&locations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve locations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Locations retrieved successfully", "data": locations})
}

// GetLocation retrieves a single location by ID
func GetLocation(c *gin.Context) {
	id := c.Param("id")
	var location models.Location

	if err := database.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location fetched successfully", "data": location})
}

// CreateLocation creates a new location
func CreateLocation(c *gin.Context) {
	var location models.Location

	if err := c.ShouldBind(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Save to DB
	if err := database.DB.Create(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Location created successfully",
		"data":    location,
	})
}

// UpdateLocation updates an existing location
func UpdateLocation(c *gin.Context) {
	id := c.Param("id")

	var location models.Location
	if err := database.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	// updates := make(map[string]interface{})

	// if v := c.PostForm("name"); v != "" {
	// 	updates["name"] = v
	// }
	// if v := c.PostForm("address"); v != "" {
	// 	updates["address"] = v
	// }
	// if v := c.PostForm("contact_person"); v != "" {
	// 	updates["contact_person"] = v
	// }
	// if v := c.PostForm("is_default"); v != "" {
	// 	updates["is_default"] = v == "true"
	// }

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := database.DB.Model(&location).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update location"})
		return
	}

	// if len(location) == 0 {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
	// 	return
	// }

	// if err := database.DB.Model(&location).Updates(location).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// Fetch updated location to return
	// database.DB.First(&location, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Location updated successfully",
		"data":    location,
	})
}

// DeleteLocation deletes a location
func DeleteLocation(c *gin.Context) {
	id := c.Param("id")
	var location models.Location

	if err := database.DB.First(&location, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Location not found"})
		return
	}

	if err := database.DB.Delete(&location).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Location deleted successfully",
	})
}
