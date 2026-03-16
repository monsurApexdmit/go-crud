package controllers

import (
	"go-crud/database"
	"go-crud/dto"
	"go-crud/middlewares"
	"go-crud/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetBillingContact gets billing contact information
func GetBillingContact(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var billingContact models.BillingContact
	result := database.DB.Where("company_id = ?", companyID).First(&billingContact)

	// If no billing contact exists, return empty/default response
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, dto.APIResponse{
			Message: "Billing contact not set",
			Data: dto.BillingContactDTO{
				CompanyID: companyID,
				IsDefault: false,
			},
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Billing contact retrieved",
		Data: dto.BillingContactDTO{
			ID:        billingContact.ID,
			CompanyID: billingContact.CompanyID,
			Email:     billingContact.Email,
			Phone:     billingContact.Phone,
			Address:   billingContact.Address,
			City:      billingContact.City,
			State:     billingContact.State,
			ZipCode:   billingContact.ZipCode,
			Country:   billingContact.Country,
			TaxID:     billingContact.TaxID,
			TaxIDType: billingContact.TaxIDType,
			IsDefault: billingContact.IsDefault,
			CreatedAt: billingContact.CreatedAt,
			UpdatedAt: billingContact.UpdatedAt,
		},
	})
}

// UpdateBillingContact updates billing contact information
func UpdateBillingContact(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var req dto.BillingContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Check if billing contact exists for company
	var billingContact models.BillingContact
	existingContact := database.DB.Where("company_id = ?", companyID).First(&billingContact).RowsAffected > 0

	updateData := map[string]interface{}{
		"email":          req.Email,
		"phone":          req.Phone,
		"address":        req.Address,
		"city":           req.City,
		"state":          req.State,
		"zip_code":       req.ZipCode,
		"country":        req.Country,
		"tax_id":         req.TaxID,
		"tax_id_type":    req.TaxIDType,
	}

	if existingContact {
		// Update existing
		if err := database.DB.Model(&billingContact).Updates(updateData).Where("company_id = ?", uint(companyID)).First(&billingContact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to update billing contact",
			})
			return
		}
	} else {
		// Create new
		newContact := models.BillingContact{
			CompanyID: companyID,
			Email:     req.Email,
			Phone:     req.Phone,
			Address:   req.Address,
			City:      req.City,
			State:     req.State,
			ZipCode:   req.ZipCode,
			Country:   req.Country,
			TaxID:     req.TaxID,
			TaxIDType: req.TaxIDType,
			IsDefault: true,
		}
		if err := database.DB.Create(&newContact).Error; err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to create billing contact",
			})
			return
		}
		billingContact = newContact
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Billing contact updated",
		Data: dto.BillingContactDTO{
			ID:        billingContact.ID,
			CompanyID: billingContact.CompanyID,
			Email:     billingContact.Email,
			Phone:     billingContact.Phone,
			Address:   billingContact.Address,
			City:      billingContact.City,
			State:     billingContact.State,
			ZipCode:   billingContact.ZipCode,
			Country:   billingContact.Country,
			TaxID:     billingContact.TaxID,
			TaxIDType: billingContact.TaxIDType,
			IsDefault: billingContact.IsDefault,
			CreatedAt: billingContact.CreatedAt,
			UpdatedAt: billingContact.UpdatedAt,
		},
	})
}
