package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/middlewares"
	"go-crud/models"
)

// GetSettings retrieves all settings
func GetSettings(c *gin.Context) {
	var settings models.Settings
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		// Return default settings if none exist
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default settings"})
			return
		}
	}

	response := gin.H{
		"general":       json.RawMessage(settings.GeneralSettings),
		"tax":           json.RawMessage(settings.TaxSettings),
		"shipping":      json.RawMessage(settings.ShippingSettings),
		"payment":       json.RawMessage(settings.PaymentSettings),
		"business":      json.RawMessage(settings.BusinessSettings),
		"regional":      json.RawMessage(settings.RegionalSettings),
		"notifications": json.RawMessage(settings.NotificationSettings),
		"storeHours":    json.RawMessage(settings.StoreHours),
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Settings retrieved successfully",
		"data":    response,
	})
}

// UpdateGeneralSettings updates general settings
func UpdateGeneralSettings(c *gin.Context) {
	var input models.GeneralSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		// Upsert: create if not exists
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetGeneralSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("general_settings", settings.GeneralSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "General settings updated successfully",
		"data":    input,
	})
}

// UpdateTaxSettings updates tax settings
func UpdateTaxSettings(c *gin.Context) {
	var input models.TaxSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		// Upsert: create if not exists
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetTaxSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("tax_settings", settings.TaxSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tax settings updated successfully",
		"data":    input,
	})
}

// UpdateShippingSettings updates shipping settings
func UpdateShippingSettings(c *gin.Context) {
	var input models.ShippingSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetShippingSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("shipping_settings", settings.ShippingSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shipping settings updated successfully",
		"data":    input,
	})
}

// UpdatePaymentSettings updates payment settings
func UpdatePaymentSettings(c *gin.Context) {
	var input models.PaymentSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetPaymentSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("payment_settings", settings.PaymentSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment settings updated successfully",
		"data":    input,
	})
}

// UpdateBusinessSettings updates business settings
func UpdateBusinessSettings(c *gin.Context) {
	var input models.BusinessSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetBusinessSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("business_settings", settings.BusinessSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Business settings updated successfully",
		"data":    input,
	})
}

// UpdateRegionalSettings updates regional settings
func UpdateRegionalSettings(c *gin.Context) {
	var input models.RegionalSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetRegionalSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("regional_settings", settings.RegionalSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Regional settings updated successfully",
		"data":    input,
	})
}

// UpdateNotificationSettings updates notification settings
func UpdateNotificationSettings(c *gin.Context) {
	var input models.NotificationSettings
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetNotificationSettings(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("notification_settings", settings.NotificationSettings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Notification settings updated successfully",
		"data":    input,
	})
}

// UpdateStoreHours updates store hours
func UpdateStoreHours(c *gin.Context) {
	var input models.StoreHoursData
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	if err := settings.SetStoreHours(input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("store_hours", settings.StoreHours).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Store hours updated successfully",
		"data":    input,
	})
}

// UploadLogo uploads store logo
func UploadLogo(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 5MB"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !validExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif, webp are allowed"})
		return
	}

	// Create uploads directory if not exists
	uploadDir := "uploads/logos"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("logo-%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update settings with logo URL
	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	logoURL := "/uploads/logos/" + filename
	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("logo_url", logoURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logo uploaded successfully",
		"data": gin.H{
			"logoUrl": logoURL,
		},
	})
}

// UploadBanner uploads store banner
func UploadBanner(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size must not exceed 10MB"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !validExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif, webp are allowed"})
		return
	}

	// Create uploads directory if not exists
	uploadDir := "uploads/banners"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("banner-%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Update settings with banner URL
	var settings models.Settings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		settings = initializeDefaultSettings()
		settings.CompanyID = companyID
		if err := database.DB.Create(&settings).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create settings"})
			return
		}
	}

	bannerURL := "/uploads/banners/" + filename
	if err := database.DB.Where("company_id = ?", companyID).Model(&settings).Update("banner_url", bannerURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Banner uploaded successfully",
		"data": gin.H{
			"bannerUrl": bannerURL,
		},
	})
}

// Helper function to initialize default settings
func initializeDefaultSettings() models.Settings {
	general := models.GeneralSettings{
		StoreName:    "My Store",
		StoreEmail:   "store@example.com",
		StorePhone:   "+1 234 567 8900",
		StoreAddress: "123 Main St, City, State",
		StoreDescription: "Welcome to our store",
	}

	tax := models.TaxSettings{
		DefaultTaxRate:     10,
		TaxInclusivePrice:  false,
		EnableGSTTracking:  false,
		EnableTaxExemption: false,
		DefaultShippingTax: 0,
	}

	shipping := models.ShippingSettings{
		EnableShipping:        true,
		DefaultShippingCost:   5.99,
		FreeShippingThreshold: 50,
		ShippingMethods: []models.ShippingMethod{
			{
				ID:            "standard",
				Name:          "Standard Shipping",
				Cost:          5.99,
				EstimatedDays: 5,
				IsActive:      true,
			},
		},
	}

	payment := models.PaymentSettings{
		EnableCash:        true,
		EnableCard:        true,
		EnableOnline:      true,
		CardProcessingFee: 2.9,
	}

	business := models.BusinessSettings{
		BusinessName:       "My Business",
		BusinessType:       "Retail",
		RegistrationNumber: "",
		GSTNumber:          "",
		Website:            "https://example.com",
		SocialLinks: models.SocialLinks{
			Facebook:  "",
			Instagram: "",
			Twitter:   "",
		},
	}

	regional := models.RegionalSettings{
		Language: "en-US",
		Currency: "USD",
		Timezone: "UTC",
	}

	notifications := models.NotificationSettings{
		EmailNotifications: true,
		OrderNotifications: true,
		MarketingEmails:    false,
	}

	storeHours := models.StoreHoursData{
		Monday:    models.StoreHours{Open: "09:00", Close: "17:00", IsOpen: true},
		Tuesday:   models.StoreHours{Open: "09:00", Close: "17:00", IsOpen: true},
		Wednesday: models.StoreHours{Open: "09:00", Close: "17:00", IsOpen: true},
		Thursday:  models.StoreHours{Open: "09:00", Close: "17:00", IsOpen: true},
		Friday:    models.StoreHours{Open: "09:00", Close: "17:00", IsOpen: true},
		Saturday:  models.StoreHours{Open: "10:00", Close: "15:00", IsOpen: true},
		Sunday:    models.StoreHours{Open: "00:00", Close: "00:00", IsOpen: false},
	}

	settings := models.Settings{}
	settings.SetGeneralSettings(general)
	settings.SetTaxSettings(tax)
	settings.SetShippingSettings(shipping)
	settings.SetPaymentSettings(payment)
	settings.SetBusinessSettings(business)
	settings.SetRegionalSettings(regional)
	settings.SetNotificationSettings(notifications)
	settings.SetStoreHours(storeHours)

	return settings
}
