package controllers

import (
	"net/http"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
)

func ListUsers(c *gin.Context) {
    var users []models.User
    if err :=  database.DB.Preload("Role").Find(&users).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Users retrieved successfully", "data": users})
}

func GetUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User

    if err := database.DB.Preload("Role").First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User fetched successfully", "data": user})
}

func CreateUser(c *gin.Context) {
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
        return
    }

    if user.Username == "" || user.Email == "" || user.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email and password are required"})
        return
    }
    
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
        return
    }
    user.Password = string(hashedPassword)

    if err := database.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    if err := database.DB.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load role"})
		return
	}

    user.Password = ""

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": user})
}

func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    var updatedData models.User
    if err := c.ShouldBindJSON(&updatedData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
        return
    }

    user.Username = updatedData.Username
    user.Email = updatedData.Email
    user.Address = updatedData.Address

    if updatedData.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword(
            []byte(updatedData.Password),
            bcrypt.DefaultCost,
        )
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
            return
        }
        user.Password = string(hashedPassword)
    }

    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    user.Password = ""

    c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "data": user})
}


func DeleteUser(c *gin.Context) {
    id := c.Param("id")

    if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    c.Status(http.StatusNoContent)
}
