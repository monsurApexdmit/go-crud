package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"go-crud/database"
	"go-crud/dto"
	"go-crud/middlewares"
	"go-crud/models"
	"go-crud/services"
	"go-crud/utils"

	"github.com/gin-gonic/gin"
)

// ListCoupons retrieves all coupons
func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("company_id = ?", companyID).Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupons retrieved successfully", "data": coupons})
}

// GetCoupon retrieves a single coupon by ID
func GetCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon

	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon fetched successfully", "data": coupon})
}

// GetCouponByCode retrieves a coupon by its code (public endpoint)
func GetCouponByCode(c *gin.Context) {
	code := c.Param("code")

	coupon, err := services.GetCouponByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found or inactive"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon found",
		"data":    coupon,
	})
}

// CreateCoupon creates a new coupon with JSON support
func CreateCoupon(c *gin.Context) {
	var request dto.CreateCouponRequest

	// Use JSON binding
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleValidationErrors(c, err)
		return
	}

	// Validate date range
	if request.EndDate.Before(request.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "End date must be after start date",
		})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)

	// Map DTO to Model
	coupon := models.Coupon{
		CampaignName:           request.CampaignName,
		Code:                   request.Code,
		Discount:               request.Discount,
		Type:                   request.Type,
		StartDate:              request.StartDate,
		EndDate:                request.EndDate,
		Status:                 request.Status,
		UsageLimit:             request.UsageLimit,
		UsageLimitPerUser:      request.UsageLimitPerUser,
		MinOrderAmount:         request.MinOrderAmount,
		MaxDiscount:            request.MaxDiscount,
		ApplicableToCategories: request.ApplicableToCategories,
		ApplicableToProducts:   request.ApplicableToProducts,
		FreeShipping:           request.FreeShipping,
		Stackable:              request.Stackable,
		AutoApply:              request.AutoApply,
		Priority:               request.Priority,
		CompanyID:              companyID,
	}

	// Save to DB
	if err := database.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create coupon: " + err.Error(),
		})
		return
	}

	log.Println("✅ Coupon created successfully:", coupon.Code)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Coupon created successfully",
		"data":    coupon,
	})
}

// CreateCouponWithImage creates a coupon with form-data (for image upload)
func CreateCouponWithImage(c *gin.Context) {
	var coupon models.Coupon

	if err := c.ShouldBind(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)
	coupon.CompanyID = companyID

	startStr := c.PostForm("start_date")
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start_date format (use RFC3339)",
		})
		return
	}

	endStr := c.PostForm("end_date")
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end_date format (use RFC3339)",
		})
		return
	}

	coupon.StartDate = start
	coupon.EndDate = end

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		path, err := utils.SaveUploadedFile(c, file, "uploads/coupons")
		if err != nil {
			log.Println("❌ Image upload error:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save image",
			})
			return
		}
		coupon.Image = path
	}

	if err := database.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Coupon created successfully",
		"data":    coupon,
	})
}

// UpdateCoupon updates an existing coupon
func UpdateCoupon(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)
	var coupon models.Coupon
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	var request dto.UpdateCouponRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleValidationErrors(c, err)
		return
	}

	// Update fields if provided
	if request.CampaignName != nil {
		coupon.CampaignName = *request.CampaignName
	}
	if request.Code != nil {
		coupon.Code = *request.Code
	}
	if request.Discount != nil {
		coupon.Discount = *request.Discount
	}
	if request.Type != nil {
		coupon.Type = *request.Type
	}
	if request.StartDate != nil {
		coupon.StartDate = *request.StartDate
	}
	if request.EndDate != nil {
		coupon.EndDate = *request.EndDate
	}
	if request.Status != nil {
		coupon.Status = *request.Status
	}
	if request.UsageLimit != nil {
		coupon.UsageLimit = request.UsageLimit
	}
	if request.UsageLimitPerUser != nil {
		coupon.UsageLimitPerUser = request.UsageLimitPerUser
	}
	if request.MinOrderAmount != nil {
		coupon.MinOrderAmount = *request.MinOrderAmount
	}
	if request.MaxDiscount != nil {
		coupon.MaxDiscount = request.MaxDiscount
	}
	if request.ApplicableToCategories != nil {
		coupon.ApplicableToCategories = *request.ApplicableToCategories
	}
	if request.ApplicableToProducts != nil {
		coupon.ApplicableToProducts = *request.ApplicableToProducts
	}
	if request.FreeShipping != nil {
		coupon.FreeShipping = *request.FreeShipping
	}
	if request.Stackable != nil {
		coupon.Stackable = *request.Stackable
	}
	if request.AutoApply != nil {
		coupon.AutoApply = *request.AutoApply
	}
	if request.Priority != nil {
		coupon.Priority = *request.Priority
	}

	// Validate date range if both are updated
	if coupon.EndDate.Before(coupon.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "End date must be after start date",
		})
		return
	}

	if err := database.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon updated successfully",
		"data":    coupon,
	})
}

// UpdateCouponWithImage updates coupon with form-data (for image replacement)
func UpdateCouponWithImage(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)
	var coupon models.Coupon
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&coupon).Error; err != nil {
		c.JSON(404, gin.H{"error": "Coupon not found"})
		return
	}

	updates := make(map[string]interface{})

	if v := c.PostForm("campaign_name"); v != "" {
		updates["campaign_name"] = v
	}
	if v := c.PostForm("code"); v != "" {
		updates["code"] = v
	}
	if v := c.PostForm("type"); v != "" {
		updates["type"] = v
	}
	if v := c.PostForm("status"); v != "" {
		updates["status"] = v == "true"
	}
	if v := c.PostForm("start_date"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			updates["start_date"] = t
		}
	}
	if v := c.PostForm("end_date"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			updates["end_date"] = t
		}
	}

	// Handle image replacement
	file, err := c.FormFile("image")
	if err == nil {
		if coupon.Image != "" {
			_ = os.Remove(coupon.Image)
		}

		path, err := utils.SaveUploadedFile(c, file, "uploads/coupons")
		if err != nil {
			c.JSON(500, gin.H{"error": "Image upload failed"})
			return
		}
		updates["image"] = path
	}

	if len(updates) == 0 {
		c.JSON(400, gin.H{"error": "No fields to update"})
		return
	}

	if err := database.DB.Model(&coupon).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "Coupon updated successfully",
		"data":    updates,
	})
}

// DeleteCoupon soft deletes a coupon
func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon

	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	// Delete image if exists
	if coupon.Image != "" {
		_ = os.Remove(coupon.Image)
	}

	if err := database.DB.Delete(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coupon deleted successfully",
	})
}

// ValidateCoupon validates a coupon code for checkout
func ValidateCoupon(c *gin.Context) {
	var request dto.ValidateCouponRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleValidationErrors(c, err)
		return
	}

	// Convert DTO to service request
	serviceReq := services.CouponValidationRequest{
		Code:        request.Code,
		CustomerID:  request.CustomerID,
		OrderAmount: request.OrderAmount,
		CartItems:   make([]services.CartItem, len(request.CartItems)),
	}

	for i, item := range request.CartItems {
		serviceReq.CartItems[i] = services.CartItem{
			ProductID:  item.ProductID,
			CategoryID: item.CategoryID,
			Price:      item.Price,
			Quantity:   item.Quantity,
		}
	}

	// Validate coupon
	result, err := services.ValidateCoupon(serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Validation failed: " + err.Error()})
		return
	}

	if !result.Valid {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"valid":      false,
			"error_code": result.ErrorCode,
			"message":    result.Message,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetCouponUsageStats retrieves usage analytics for a coupon
func GetCouponUsageStats(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)
	// Verify coupon exists and belongs to company
	var coupon models.Coupon
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&coupon).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	// Convert string ID to uint
	var couponID uint
	if _, err := fmt.Sscanf(id, "%d", &couponID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID"})
		return
	}

	stats, err := services.GetCouponUsageStats(couponID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usage statistics retrieved successfully",
		"data":    stats,
	})
}
