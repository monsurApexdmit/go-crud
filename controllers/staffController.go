package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"go-crud/database"
	"go-crud/middlewares"
	"go-crud/models"
)

func ListStaff(c *gin.Context) {
	var staffList []models.Staff

	companyID, _ := middlewares.GetCompanyID(c)
	query := database.DB.Preload("User").Preload("User.Role").Model(&models.Staff{}).Where("company_id = ?", companyID)

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR contact LIKE ?", like, like, like)
	}
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if role := c.Query("role"); role != "" && role != "all" {
		query = query.Where("role = ?", role)
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

	if err := query.Offset(offset).Limit(limit).Find(&staffList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve staff"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Staff retrieved successfully",
		"data":    staffList,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetStaff(c *gin.Context) {
	var staff models.Staff
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).Preload("User").Preload("User.Role").First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Staff fetched successfully", "data": staff})
}

func CreateStaff(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if staff.Name == "" || staff.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)
	staff.CompanyID = companyID

	// Create linked user with Staff role (id=5)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user := models.User{
		Username: staff.Name,
		Email:    staff.Email,
		Password: string(hashedPassword),
		RoleID:   5,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
	staff.UserID = &user.ID

	if err := database.DB.Create(&staff).Error; err != nil {
		database.DB.Unscoped().Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create staff"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Staff created successfully", "data": staff})
}

func UpdateStaff(c *gin.Context) {
	var staff models.Staff
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := database.DB.Model(&staff).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update staff"})
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
	if len(userUpdate) > 0 && staff.UserID != nil {
		database.DB.Model(&models.User{}).Where("id = ?", *staff.UserID).Updates(userUpdate)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff updated successfully", "data": staff})
}

func DeleteStaff(c *gin.Context) {
	var staff models.Staff
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&staff).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	database.DB.Delete(&staff)

	if staff.UserID != nil {
		database.DB.Delete(&models.User{}, *staff.UserID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Staff deleted successfully"})
}
