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

// ListVendorReturns retrieves vendor returns with filtering and pagination
func ListVendorReturns(c *gin.Context) {
	var returns []models.VendorReturn
	query := database.DB.Model(&models.VendorReturn{})

	// Search by return number or vendor name
	if search := c.Query("search"); search != "" {
		query = query.Where("return_number LIKE ? OR vendor_name LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Filter by status
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Filter by vendor ID
	if vendorID := c.Query("vendor_id"); vendorID != "" {
		query = query.Where("vendor_id = ?", vendorID)
	}

	// Order by most recent first
	query = query.Order("return_date DESC")

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
	if err := query.Offset(offset).Limit(perPage).Preload("Vendor").Preload("Items").Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vendor returns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Vendor returns retrieved successfully",
		"data":    returns,
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

// GetVendorReturn retrieves a single vendor return by ID
func GetVendorReturn(c *gin.Context) {
	var ret models.VendorReturn
	if err := database.DB.Preload("Vendor").Preload("Items").First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor return not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Vendor return fetched successfully", "data": ret})
}

// CreateVendorReturn creates a new vendor return
func CreateVendorReturn(c *gin.Context) {
	var ret models.VendorReturn
	if err := c.ShouldBindJSON(&ret); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Generate return number if not provided
	if ret.ReturnNumber == "" {
		ret.ReturnNumber = generateVendorReturnNumber()
	}

	// Set return date to now if not provided
	if ret.ReturnDate.IsZero() {
		ret.ReturnDate = time.Now()
	}

	// Validate required fields
	if ret.VendorName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vendor name is required"})
		return
	}
	if ret.CreatedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CreatedBy is required"})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "shipped", "received_by_vendor", "completed"}
	if ret.Status != "" && !contains(validStatuses, ret.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}
	if ret.Status == "" {
		ret.Status = "pending"
	}

	// Validate credit type
	validCreditTypes := []string{"refund", "credit_note", "replacement"}
	if !contains(validCreditTypes, ret.CreditType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credit type"})
		return
	}

	// Create return, deduct inventory in transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create return record
		if err := tx.Create(&ret).Error; err != nil {
			return err
		}

		// Deduct items from inventory when returning to vendor
		for _, item := range ret.Items {
			if item.VariantID != nil {
				// For variant products: deduct from variant_inventory
				var inventory models.VariantInventory
				if err := tx.Where("variant_id = ?", *item.VariantID).First(&inventory).Error; err != nil {
					// If no inventory record found, skip (can't deduct from nothing)
					continue
				}

				if inventory.Quantity < item.Quantity {
					return fmt.Errorf("insufficient inventory for variant %s", item.VariantName)
				}

				inventory.Quantity -= item.Quantity
				if err := tx.Save(&inventory).Error; err != nil {
					return err
				}
			} else if item.ProductID != nil {
				// For simple products: deduct from products.stock
				var product models.Product
				if err := tx.First(&product, *item.ProductID).Error; err != nil {
					continue
				}

				if product.Stock < item.Quantity {
					return fmt.Errorf("insufficient stock for product %s", item.ProductName)
				}

				if err := tx.Model(&models.Product{}).
					Where("id = ?", *item.ProductID).
					UpdateColumn("stock", gorm.Expr("stock - ?", item.Quantity)).
					Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "Return number already exists"})
			return
		}
		if strings.Contains(err.Error(), "insufficient") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vendor return"})
		return
	}

	database.DB.Preload("Vendor").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Vendor return created and inventory deducted successfully", "data": ret})
}

// UpdateVendorReturn updates an existing vendor return
func UpdateVendorReturn(c *gin.Context) {
	var ret models.VendorReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor return not found"})
		return
	}

	var input models.VendorReturn
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update allowed fields
	if input.Status != "" {
		validStatuses := []string{"pending", "shipped", "received_by_vendor", "completed"}
		if !contains(validStatuses, input.Status) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		ret.Status = input.Status

		// Set completed date when status becomes completed
		if input.Status == "completed" && ret.CompletedDate == nil {
			now := time.Now()
			ret.CompletedDate = &now
		}
	}
	if input.Notes != "" {
		ret.Notes = input.Notes
	}

	if err := database.DB.Save(&ret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vendor return"})
		return
	}

	database.DB.Preload("Vendor").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Vendor return updated successfully", "data": ret})
}

// UpdateVendorReturnStatus updates only the status of a vendor return
func UpdateVendorReturnStatus(c *gin.Context) {
	var ret models.VendorReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor return not found"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	validStatuses := []string{"pending", "shipped", "received_by_vendor", "completed"}
	if !contains(validStatuses, input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	ret.Status = input.Status

	// Set completed date when status becomes completed
	if input.Status == "completed" && ret.CompletedDate == nil {
		now := time.Now()
		ret.CompletedDate = &now
	}

	if err := database.DB.Save(&ret).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	database.DB.Preload("Vendor").Preload("Items").First(&ret, ret.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully", "data": ret})
}

// DeleteVendorReturn soft-deletes a vendor return
func DeleteVendorReturn(c *gin.Context) {
	var ret models.VendorReturn
	if err := database.DB.First(&ret, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Vendor return not found"})
		return
	}

	database.DB.Delete(&ret)
	c.JSON(http.StatusOK, gin.H{"message": "Vendor return deleted successfully"})
}

// GetVendorReturnStats returns statistics for vendor returns
func GetVendorReturnStats(c *gin.Context) {
	var stats struct {
		Total           int64   `json:"total"`
		Pending         int64   `json:"pending"`
		Shipped         int64   `json:"shipped"`
		Completed       int64   `json:"completed"`
		TotalCreditAmount float64 `json:"totalCreditAmount"`
	}

	database.DB.Model(&models.VendorReturn{}).Count(&stats.Total)
	database.DB.Model(&models.VendorReturn{}).Where("status = ?", "pending").Count(&stats.Pending)
	database.DB.Model(&models.VendorReturn{}).Where("status = ?", "shipped").Count(&stats.Shipped)
	database.DB.Model(&models.VendorReturn{}).Where("status = ?", "completed").Count(&stats.Completed)
	database.DB.Model(&models.VendorReturn{}).
		Where("status = ?", "completed").
		Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TotalCreditAmount)

	c.JSON(http.StatusOK, gin.H{"message": "Statistics retrieved successfully", "data": stats})
}

// GetVendorReturnsByVendor retrieves returns for a specific vendor
func GetVendorReturnsByVendor(c *gin.Context) {
	var returns []models.VendorReturn
	vendorID := c.Param("vendorId")

	if err := database.DB.Where("vendor_id = ?", vendorID).
		Preload("Vendor").Preload("Items").
		Order("return_date DESC").
		Find(&returns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve vendor returns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vendor returns retrieved successfully", "data": returns})
}

// generateVendorReturnNumber generates a unique vendor return number
func generateVendorReturnNumber() string {
	return fmt.Sprintf("VRT-%d", time.Now().UnixNano()/1000)
}
