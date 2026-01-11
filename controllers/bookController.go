package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
)

func ListBooks(c *gin.Context) {
    var books []models.Book
    if err := database.DB.Preload("Author").Find(&books).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve books"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Books retrieved successfully", "data": books})
}

func GetBook(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    if err := database.DB.Preload("Author").First(&book, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book fetched successfully", "data": book})
}

func CreateBook(c *gin.Context) {
    var book models.Book

    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
        return
    }

    // Check if author exists
    var author models.Author
    if err := database.DB.First(&author, book.AuthorID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
        return
    }

    if err := database.DB.Create(&book).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
        return
    }

        // Reload book with author populated
    var createdBook models.Book
    if err := database.DB.Preload("Author").First(&createdBook, book.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load created book"})
        return
    }


    c.JSON(http.StatusCreated, gin.H{"message": "Book created successfully", "data": createdBook})
}

func UpdateBook(c *gin.Context) {
    id := c.Param("id")
    var book models.Book

    // Check if book exists
    if err := database.DB.First(&book, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    var updatedData models.Book
    if err := c.ShouldBindJSON(&updatedData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
        return
    }

    // Validate new author
    var author models.Author
    if err := database.DB.First(&author, updatedData.AuthorID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Author not found"})
        return
    }

    // Update fields
    book.Title = updatedData.Title
    book.AuthorID = updatedData.AuthorID

    // Save updated book
    if err := database.DB.Save(&book).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
        return
    }

    // Reload book with author relationship
    var updatedBook models.Book
    if err := database.DB.
        Preload("Author").
        First(&updatedBook, book.ID).Error; err != nil {

        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load updated book"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "data": updatedBook})
}


func DeleteBook(c *gin.Context) {
    id := c.Param("id")

    if err := database.DB.Delete(&models.Book{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
        return
    }

    // No body should be sent for 204 response
    c.Status(http.StatusNoContent)
}
