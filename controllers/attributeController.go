package controllers

import (
	"go-crud/database"
	"go-crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListAttributes retrieves all attributes
func ListAttributes(c *gin.Context) {
	var attributes []models.Attribute
	if err := database.DB.Find(&attributes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attributes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attributes retrieved successfully",
		"data":    attributes,
	})
}

// GetAttribute retrieves a single attribute by ID
func GetAttribute(c *gin.Context) {
	id := c.Param("id")
	var attribute models.Attribute

	if err := database.DB.First(&attribute, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute fetched successfully",
		"data":    attribute,
	})
}

// CreateAttribute creates a new attribute
func CreateAttribute(c *gin.Context) {
	var attribute models.Attribute

	if err := c.ShouldBindJSON(&attribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if attribute.Name == "" || attribute.DisplayName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and DisplayName are required"})
		return
	}

	if err := database.DB.Create(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attribute"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Attribute created successfully",
		"data":    attribute,
	})
}

// UpdateAttribute updates an existing attribute
func UpdateAttribute(c *gin.Context) {
	id := c.Param("id")
	var attribute models.Attribute

	if err := database.DB.First(&attribute, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := database.DB.Model(&attribute).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attribute"})
		return
	}

	// Fetch updated attribute
	database.DB.First(&attribute, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute updated successfully",
		"data":    attribute,
	})
}

// DeleteAttribute deletes an attribute
func DeleteAttribute(c *gin.Context) {
	id := c.Param("id")

	result := database.DB.Delete(&models.Attribute{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attribute"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute deleted successfully",
	})
}
