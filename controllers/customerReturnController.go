package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
	"gorm.io/gorm"
)

// ListCustomerReturns retrieves customer returns with filtering and pagination
func ListCustomerReturns(c *gin.Context) {
	var returns []models.CustomerReturn
	query := database.DB.Model(&models.CustomerReturn{})

	// Search by return number, customer name, or order number
	if search := c.Query("search"); search != "" {
		query = query.Where("return_number LIKE ? OR customer_name LIKE ? OR order_number LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Filter by status
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Filter by customer ID
	if customerID := c.Query("customer_id"); customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// Order by most recent first
	query = query.Order("request_date DESC")

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).Preload("Customer").Preload("Items").Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer returns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer returns retrieved successfully",
		"data":    returns,
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

// GetCustomerReturn retrieves a single customer return by ID
func GetCustomerReturn(c *gin.Context) {
	var ret models.CustomerReturn
	if err := database.DB.Preload("Customer").Preload("Items").First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer return not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer return fetched successfully", "data": ret})
}

// CreateCustomerReturn creates a new customer return
func CreateCustomerReturn(c *gin.Context) {
	var ret models.CustomerReturn
	if err := c.ShouldBindJSON(&ret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Generate return number if not provided
	if ret.ReturnNumber == "" {
		ret.ReturnNumber = generateCustomerReturnNumber()
	}

	// Set request date to now if not provided
	if ret.RequestDate.IsZero() {
		ret.RequestDate = time.Now()
	}

	// Validate required fields
	if ret.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer name is required"})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "approved", "rejected", "completed"}
	if ret.Status != "" && !contains(validStatuses, ret.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	if ret.Status == "" {
		ret.Status = "pending"
	}

	// Validate refund method
	validMethods := []string{"cash", "store_credit", "original_payment"}
	if !contains(validMethods, ret.RefundMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refund method"})
		return
	}

	// Create return and items in transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&ret).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "Return number already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer return"})
		return
	}

	database.DB.Preload("Customer").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Customer return created successfully", "data": ret})
}

// UpdateCustomerReturn updates an existing customer return
func UpdateCustomerReturn(c *gin.Context) {
	var ret models.CustomerReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer return not found"})
		return
	}

	var input models.CustomerReturn
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update allowed fields
	if input.Status != "" {
		validStatuses := []string{"pending", "approved", "rejected", "completed"}
		if !contains(validStatuses, input.Status) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		ret.Status = input.Status
	}
	if input.Notes != "" {
		ret.Notes = input.Notes
	}
	if input.ProcessedBy != "" {
		ret.ProcessedBy = input.ProcessedBy
	}

	if err := database.DB.Save(&ret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer return"})
		return
	}

	database.DB.Preload("Customer").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Customer return updated successfully", "data": ret})
}

// ApproveCustomerReturn approves a customer return
func ApproveCustomerReturn(c *gin.Context) {
	var ret models.CustomerReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer return not found"})
		return
	}

	if ret.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only approve pending returns"})
		return
	}

	var input struct {
		ProcessedBy string `json:"processedBy" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProcessedBy is required"})
		return
	}

	now := time.Now()
	ret.Status = "approved"
	ret.ProcessedDate = &now
	ret.ProcessedBy = input.ProcessedBy

	// Load items for inventory update
	var items []models.CustomerReturnItem
	if err := database.DB.Where("return_id = ?", ret.ID).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load return items"})
		return
	}

	// Update return and restock inventory in a transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Save return with approved status
		if err := tx.Save(&ret).Error; err != nil {
			return err
		}

		// Restock items back to inventory
		for _, item := range items {
			if item.VariantID != nil {
				// For variant products: update variant_inventory
				// Find existing inventory record or create new one
				var inventory models.VariantInventory
				result := tx.Where("variant_id = ?", *item.VariantID).First(&inventory)

				if result.Error == gorm.ErrRecordNotFound {
					// Create new inventory record (assuming default location)
					inventory = models.VariantInventory{
						VariantID: *item.VariantID,
						LocationID: 1, // Default location - adjust as needed
						Quantity: item.Quantity,
					}
					if err := tx.Create(&inventory).Error; err != nil {
						return err
					}
				} else if result.Error != nil {
					return result.Error
				} else {
					// Update existing inventory
					inventory.Quantity += item.Quantity
					if err := tx.Save(&inventory).Error; err != nil {
						return err
					}
				}
			} else if item.ProductID != nil {
				// For simple products: update products.stock
				if err := tx.Model(&models.Product{}).
					Where("id = ?", *item.ProductID).
					UpdateColumn("stock", gorm.Expr("stock + ?", item.Quantity)).
					Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve return and restock inventory"})
		return
	}

	database.DB.Preload("Customer").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Customer return approved and inventory restocked successfully", "data": ret})
}

// RejectCustomerReturn rejects a customer return
func RejectCustomerReturn(c *gin.Context) {
	var ret models.CustomerReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer return not found"})
		return
	}

	if ret.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Can only reject pending returns"})
		return
	}

	var input struct {
		ProcessedBy string `json:"processedBy" binding:"required"`
		Notes       string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProcessedBy is required"})
		return
	}

	now := time.Now()
	ret.Status = "rejected"
	ret.ProcessedDate = &now
	ret.ProcessedBy = input.ProcessedBy
	if input.Notes != "" {
		ret.Notes = input.Notes
	}

	if err := database.DB.Save(&ret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject return"})
		return
	}

	database.DB.Preload("Customer").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Customer return rejected successfully", "data": ret})
}

// DeleteCustomerReturn soft-deletes a customer return
func DeleteCustomerReturn(c *gin.Context) {
	var ret models.CustomerReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer return not found"})
		return
	}

	database.DB.Delete(&ret)
	c.JSON(http.StatusOK, gin.H{"message": "Customer return deleted successfully"})
}

// GetCustomerReturnStats returns statistics for customer returns
func GetCustomerReturnStats(c *gin.Context) {
	var stats struct {
		Total            int64   `json:"total"`
		Pending          int64   `json:"pending"`
		Approved         int64   `json:"approved"`
		Rejected         int64   `json:"rejected"`
		Completed        int64   `json:"completed"`
		TotalRefundAmount float64 `json:"totalRefundAmount"`
	}

	database.DB.Model(&models.CustomerReturn{}).Count(&stats.Total)
	database.DB.Model(&models.CustomerReturn{}).Where("status = ?", "pending").Count(&stats.Pending)
	database.DB.Model(&models.CustomerReturn{}).Where("status = ?", "approved").Count(&stats.Approved)
	database.DB.Model(&models.CustomerReturn{}).Where("status = ?", "rejected").Count(&stats.Rejected)
	database.DB.Model(&models.CustomerReturn{}).Where("status = ?", "completed").Count(&stats.Completed)
	database.DB.Model(&models.CustomerReturn{}).
		Where("status IN ?", []string{"approved", "completed"}).
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalRefundAmount)

	c.JSON(http.StatusOK, gin.H{"message": "Statistics retrieved successfully", "data": stats})
}

// GetCustomerReturnsByCustomer retrieves returns for a specific customer
func GetCustomerReturnsByCustomer(c *gin.Context) {
	var returns []models.CustomerReturn
	customerID := c.Param("customerId")

	if err := database.DB.Where("customer_id = ?", customerID).
		Preload("Customer").Preload("Items").
		Order("request_date DESC").
		Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer returns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer returns retrieved successfully", "data": returns})
}

// generateCustomerReturnNumber generates a unique return number
func generateCustomerReturnNumber() string {
	return fmt.Sprintf("RET-%d", time.Now().UnixNano()/1000)
}
