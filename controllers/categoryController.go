package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
)


// ==========================
// GET ALL CATEGORIES (TREE)
// ==========================
func ListCategories(c *gin.Context) {
	var categories []models.Category

	err := database.DB.
		Preload("Children").
		Where("parent_id IS NULL").
		Find(&categories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}


// ==========================
// GET SINGLE CATEGORY
// ==========================
func GetCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := database.DB.
		Preload("Children").
		First(&category, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category fetched successfully",
		"data":    category,
	})
}


// ==========================
// CREATE CATEGORY
// ==========================
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if category.CategoryName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category name is required"})
		return
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}


// ==========================
// UPDATE CATEGORY
// ==========================
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := database.DB.Model(&category).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    category,
	})
}


// ==========================
// DELETE CATEGORY
// ==========================
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}
