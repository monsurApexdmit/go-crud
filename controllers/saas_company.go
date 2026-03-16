package controllers

import (
	"go-crud/database"
	"go-crud/dto"
	"go-crud/middlewares"
	"go-crud/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCompanyIDFromRequest extracts company_id from context using safe type assertion
func GetCompanyIDFromRequest(c *gin.Context) (uint, error) {
	// Use the safe GetCompanyID from middlewares
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		// Try to get from query parameter as fallback
		companyIDStr := c.Query("company_id")
		if companyIDStr == "" {
			return 0, nil // Return 0 if not provided
		}

		parsedID, err := strconv.ParseUint(companyIDStr, 10, 32)
		if err != nil {
			return 0, err
		}

		return uint(parsedID), nil
	}

	return companyID, nil
}

// GetCompanyProfile gets company profile
func GetCompanyProfile(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, companyID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Company not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Company profile retrieved",
		Data: dto.CompanyDTO{
			ID:          company.ID,
			Name:        company.Name,
			Industry:    company.Industry,
			Website:     company.Website,
			Phone:       company.Phone,
			Address:     company.Address,
			City:        company.City,
			State:       company.State,
			ZipCode:     company.ZipCode,
			Country:     company.Country,
			Logo:        company.Logo,
			Description: company.Description,
			Status:      company.Status,
			CreatedAt:   company.CreatedAt,
			UpdatedAt:   company.UpdatedAt,
		},
	})
}

// UpdateCompanyProfile updates company profile
func UpdateCompanyProfile(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var req dto.UpdateCompanyProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	updateData := map[string]interface{}{}
	if req.CompanyProfile.Name != "" {
		updateData["name"] = req.CompanyProfile.Name
	}
	if req.CompanyProfile.Industry != "" {
		updateData["industry"] = req.CompanyProfile.Industry
	}
	if req.CompanyProfile.Website != "" {
		updateData["website"] = req.CompanyProfile.Website
	}
	if req.CompanyProfile.Phone != "" {
		updateData["phone"] = req.CompanyProfile.Phone
	}
	if req.CompanyProfile.Address != "" {
		updateData["address"] = req.CompanyProfile.Address
	}
	if req.CompanyProfile.City != "" {
		updateData["city"] = req.CompanyProfile.City
	}
	if req.CompanyProfile.State != "" {
		updateData["state"] = req.CompanyProfile.State
	}
	if req.CompanyProfile.ZipCode != "" {
		updateData["zip_code"] = req.CompanyProfile.ZipCode
	}
	if req.CompanyProfile.Country != "" {
		updateData["country"] = req.CompanyProfile.Country
	}
	if req.CompanyProfile.Description != "" {
		updateData["description"] = req.CompanyProfile.Description
	}

	var company models.Company
	if err := database.DB.Model(&company).Where("id = ?", uint(companyID)).Updates(updateData).First(&company, uint(companyID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update company",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Company profile updated",
		Data: dto.CompanyDTO{
			ID:          company.ID,
			Name:        company.Name,
			Industry:    company.Industry,
			Website:     company.Website,
			Phone:       company.Phone,
			Address:     company.Address,
			City:        company.City,
			State:       company.State,
			ZipCode:     company.ZipCode,
			Country:     company.Country,
			Logo:        company.Logo,
			Description: company.Description,
			Status:      company.Status,
			CreatedAt:   company.CreatedAt,
			UpdatedAt:   company.UpdatedAt,
		},
	})
}

// GetCompanyStatus gets company status
func GetCompanyStatus(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var company models.Company
	if err := database.DB.First(&company, companyID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Company not found",
		})
		return
	}

	// Count users
	var userCount int64
	database.DB.Model(&models.SaasUser{}).Where("company_id = ?", uint(companyID)).Count(&userCount)

	// Get subscription plan for max users
	var subscription models.Subscription
	var maxUsers int = 5
	if err := database.DB.Where("company_id = ?", uint(companyID)).First(&subscription).Error; err == nil {
		var plan models.SubscriptionPlan
		if err := database.DB.First(&plan, subscription.PlanID).Error; err == nil {
			maxUsers = plan.MaxUsers
		}
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Company status retrieved",
		Data: dto.CompanyStatusDTO{
			ID:       company.ID,
			Name:     company.Name,
			Status:   company.Status,
			UserCount: int(userCount),
			MaxUsers: maxUsers,
			CreatedAt: company.CreatedAt,
		},
	})
}

// GetCompanySettings gets company settings
func GetCompanySettings(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var settings models.CompanySettings
	if err := database.DB.Where("company_id = ?", companyID).First(&settings).Error; err != nil {
		// Return default/empty settings if none exist
		c.JSON(http.StatusOK, dto.APIResponse{
			Message: "Company settings retrieved",
			Data: dto.CompanySettingsDTO{
				ID:          0,
				CompanyID:   companyID,
				CompanyName: "",
				TaxID:       "",
				TaxIDType:   "",
				TaxRate:     0,
				Currency:    "",
				Timezone:    "",
				Language:    "",
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Company settings retrieved",
		Data: dto.CompanySettingsDTO{
			ID:          settings.ID,
			CompanyID:   settings.CompanyID,
			CompanyName: settings.CompanyName,
			TaxID:       settings.TaxID,
			TaxIDType:   settings.TaxIDType,
			TaxRate:     settings.TaxRate,
			Currency:    settings.Currency,
			Timezone:    settings.Timezone,
			Language:    settings.Language,
			CreatedAt:   settings.CreatedAt,
			UpdatedAt:   settings.UpdatedAt,
		},
	})
}

// UpdateCompanySettings updates company settings
func UpdateCompanySettings(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var req dto.UpdateCompanySettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Build update data map
	updateData := map[string]interface{}{
		"company_id": companyID,
	}
	if req.CompanyName != "" {
		updateData["company_name"] = req.CompanyName
	}
	if req.TaxID != "" {
		updateData["tax_id"] = req.TaxID
	}
	if req.TaxIDType != "" {
		updateData["tax_id_type"] = req.TaxIDType
	}
	if req.TaxRate > 0 {
		updateData["tax_rate"] = req.TaxRate
	}
	if req.Currency != "" {
		updateData["currency"] = req.Currency
	}
	if req.Timezone != "" {
		updateData["timezone"] = req.Timezone
	}
	if req.Language != "" {
		updateData["language"] = req.Language
	}

	// Use FirstOrCreate + Update (upsert pattern)
	var settings models.CompanySettings
	if err := database.DB.Where("company_id = ?", companyID).
		Assign(updateData).
		FirstOrCreate(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update settings",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Company settings updated",
		Data: dto.CompanySettingsDTO{
			ID:          settings.ID,
			CompanyID:   settings.CompanyID,
			CompanyName: settings.CompanyName,
			TaxID:       settings.TaxID,
			TaxIDType:   settings.TaxIDType,
			TaxRate:     settings.TaxRate,
			Currency:    settings.Currency,
			Timezone:    settings.Timezone,
			Language:    settings.Language,
			CreatedAt:   settings.CreatedAt,
			UpdatedAt:   settings.UpdatedAt,
		},
	})
}
