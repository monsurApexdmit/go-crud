package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"go-crud/database"
	"go-crud/models"
)

func ListVendors(c *gin.Context) {
	var vendors []models.Vendor

	query := database.DB.Preload("User").Preload("User.Role").Model(&models.Vendor{})

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&vendors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vendors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vendors retrieved successfully",
		"data":    vendors,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetVendor(c *gin.Context) {
	var vendor models.Vendor
	if err := database.DB.Preload("User").Preload("User.Role").First(&vendor, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vendor fetched successfully", "data": vendor})
}

func CreateVendor(c *gin.Context) {
	var vendor models.Vendor
	if err := c.ShouldBindJSON(&vendor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if vendor.Name == "" || vendor.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}

	// Create linked user with Vendor role (id=4)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user := models.User{
		Username: vendor.Name,
		Email:    vendor.Email,
		Password: string(hashedPassword),
		RoleID:   4,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
	vendor.UserID = &user.ID

	if err := database.DB.Create(&vendor).Error; err != nil {
		database.DB.Unscoped().Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vendor"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Vendor created successfully", "data": vendor})
}

func UpdateVendor(c *gin.Context) {
	var vendor models.Vendor
	if err := database.DB.First(&vendor, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := database.DB.Model(&vendor).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vendor"})
		return
	}

	// sync name/email to users table if changed
	userUpdate := map[string]interface{}{}
	if name, ok := data["name"]; ok {
		userUpdate["username"] = name
	}
	if email, ok := data["email"]; ok {
		userUpdate["email"] = email
	}
	if len(userUpdate) > 0 && vendor.UserID != nil {
		database.DB.Model(&models.User{}).Where("id = ?", *vendor.UserID).Updates(userUpdate)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor updated successfully", "data": vendor})
}

func DeleteVendor(c *gin.Context) {
	var vendor models.Vendor
	if err := database.DB.First(&vendor, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor not found"})
		return
	}

	database.DB.Delete(&vendor)

	if vendor.UserID != nil {
		database.DB.Delete(&models.User{}, *vendor.UserID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor deleted successfully"})
}
