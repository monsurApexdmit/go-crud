package controllers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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

	// log.Println("campaign_name:", c.PostForm("campaign_name"))

	if err := c.ShouldBind(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	startStr := c.PostForm("start_date")
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid start_date format",
		})
		return
	}

	endStr := c.PostForm("end_date")
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid end_date format",
		})
		return
	}

	// Assign parsed dates
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
	} else {
		log.Println("ℹ️ No image uploaded")
	}

	// Save to DB
	if err := database.DB.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("✅ Coupon created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message": "Coupon created successfully",
		"data":    coupon,
	})
}


func UpdateCoupon(c *gin.Context) {
	id := c.Param("id")

	var coupon models.Coupon
	if err := database.DB.First(&coupon, id).Error; err != nil {
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

	if v := c.PostForm("discount"); v != "" {
		if d, err := strconv.ParseFloat(v, 64); err == nil {
			updates["discount"] = d
		}
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


func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon

	if err := database.DB.First(&coupon, id).Error; err != nil {
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
		"message": "Cupon deleted successfully",
	})

}
