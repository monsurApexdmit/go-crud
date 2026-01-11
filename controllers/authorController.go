package controllers

import (
	"net/http"

	"go-crud/database"
	"go-crud/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func ListAuthors(c *gin.Context) {
	var authors []models.Author
	result := database.DB.Find(&authors)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve authors"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Authors retrieved successfully", "data": authors})
}


func CreateAuthor(c *gin.Context) {
	var author models.Author

	if err := c.ShouldBindJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	result := database.DB.Create(&author)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create author"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "Author created successfully", "data": author})
}


func GetAuthor(c *gin.Context) {
	id := c.Param("id")


	var author models.Author
	result := database.DB.First(&author, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve author"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Author retrieved successfully", "data": author})
}

func UpdateAuthor(c *gin.Context) {
	id := c.Param("id")
	var author models.Author

	if err := database.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}
	
	var updatedData models.Author
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	author.Name = updatedData.Name
	author.Email = updatedData.Email

	if err := database.DB.Save(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update author"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Author updated successfully", "data": author})
}

func DeleteAuthor(c *gin.Context) {
	id := c.Param("id")
	var author models.Author

	if err := database.DB.First(&author, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
		return
	}

	if err := database.DB.Delete(&author).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete author"})
		return
	}

	c.Status(http.StatusNoContent)
}

