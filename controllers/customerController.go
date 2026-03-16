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

func ListCustomers(c *gin.Context) {
	var customers []models.Customer

	companyID, _ := middlewares.GetCompanyID(c)
	query := database.DB.Preload("User").Preload("User.Role").Model(&models.Customer{}).Where("company_id = ?", companyID)

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", like, like, like)
	}
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if customerType := c.Query("type"); customerType != "" && customerType != "all" {
		query = query.Where("customer_type = ?", customerType)
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

	if err := query.Offset(offset).Limit(limit).Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customers retrieved successfully",
		"data":    customers,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetCustomer(c *gin.Context) {
	var customer models.Customer
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).Preload("User").Preload("User.Role").First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer fetched successfully", "data": customer})
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if customer.Name == "" || customer.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and email are required"})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)
	customer.CompanyID = companyID

	// Create linked user with Customer role (id=3)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}
	user := models.User{
		Username: customer.Name,
		Email:    customer.Email,
		Password: string(hashedPassword),
		RoleID:   3,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}
	customer.UserID = &user.ID

	if err := database.DB.Create(&customer).Error; err != nil {
		// rollback: delete the user we just created
		database.DB.Unscoped().Delete(&user)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully", "data": customer})
}

func UpdateCustomer(c *gin.Context) {
	var customer models.Customer
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := database.DB.Model(&customer).Updates(data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
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
	if len(userUpdate) > 0 && customer.UserID != nil {
		database.DB.Model(&models.User{}).Where("id = ?", *customer.UserID).Updates(userUpdate)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated successfully", "data": customer})
}

func DeleteCustomer(c *gin.Context) {
	var customer models.Customer
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	database.DB.Delete(&customer)

	// soft-delete the linked user
	if customer.UserID != nil {
		database.DB.Delete(&models.User{}, *customer.UserID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
