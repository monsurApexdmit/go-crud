package controllers

import (
	"net/http"

	"go-crud/database"
	"go-crud/models"
	"go-crud/utils"

	"github.com/gin-gonic/gin"
)

func ListCoupons(c *gin.Context) {
	var coupons []models.Coupon
	if err := database.DB.Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve coupons"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupons retrieved successfully", "data": coupons})
}

func GetCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon

	if err := database.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon fetched successfully", "data": coupon})
}

func CreateCoupon(c *gin.Context) {
	var coupon models.Coupon

	if err := c.ShouldBind(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Handle Image Upload
	file, err := c.FormFile("image")
	if err == nil {
		path, err := utils.SaveUploadedFile(c, file, "uploads/coupons")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		coupon.Image = path
	}

	if err := database.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Coupon created successfully", "data": coupon})
}

func UpdateCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon

	if err := database.DB.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
		return
	}

	var updatedData models.Coupon
	if err := c.ShouldBind(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Handle Image Upload
	file, err := c.FormFile("image")
	if err == nil {
		path, err := utils.SaveUploadedFile(c, file, "uploads/coupons")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}
		coupon.Image = path
	}

	// Update fields
	coupon.CampaignName = updatedData.CampaignName
	coupon.Code = updatedData.Code
	coupon.Discount = updatedData.Discount
	coupon.Type = updatedData.Type
	coupon.StartDate = updatedData.StartDate
	coupon.EndDate = updatedData.EndDate
	coupon.Status = updatedData.Status
	// Image is already updated if a new file was uploaded

	if err := database.DB.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coupon updated successfully", "data": coupon})
}

func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Coupon{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	c.Status(http.StatusNoContent)
}
