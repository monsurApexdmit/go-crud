package controllers

import (
	"go-crud/database"
	"go-crud/dto"
	"go-crud/middlewares"
	"go-crud/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetBillingCompanyIDFromRequest extracts company_id from context or query param
func GetBillingCompanyIDFromRequest(c *gin.Context) (uint, error) {
	// Try to get from middleware context first (set by AuthMiddleware)
	if id, exists := c.Get("company_id"); exists {
		if floatID, ok := id.(float64); ok {
			return uint(floatID), nil
		}
	}

	// Try to get from query parameter
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

// GetPlans gets all available subscription plans
func GetPlans(c *gin.Context) {
	var plans []models.SubscriptionPlan

	if err := database.DB.Where("is_active = ?", true).Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to fetch plans",
		})
		return
	}

	var planDTOs []dto.SubscriptionPlanDTO
	for _, plan := range plans {
		planDTOs = append(planDTOs, dto.SubscriptionPlanDTO{
			ID:            plan.ID,
			Name:          plan.Name,
			Description:   plan.Description,
			Price:         plan.Price,
			BillingPeriod: plan.BillingPeriod,
			MaxUsers:      plan.MaxUsers,
			MaxProducts:   plan.MaxProducts,
			MaxBranches:   plan.MaxBranches,
			Features:      plan.Features,
			IsActive:      plan.IsActive,
			IsFeatured:    plan.IsFeatured,
			CreatedAt:     plan.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Plans retrieved",
		Data:    planDTOs,
	})
}

// GetCurrentSubscription gets current subscription
func GetCurrentSubscription(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var subscription models.Subscription
	// Get the latest subscription for the company (not filtered by status initially)
	if err := database.DB.Preload("Plan").Where("company_id = ?", companyID).Order("created_at DESC").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "No subscription found for this company",
		})
		return
	}

	// If subscription exists but not active, still return it (client can check status)
	// This allows frontend to show subscription details even if expired

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Subscription retrieved",
		Data: dto.SubscriptionDTO{
			ID:                 subscription.ID,
			PlanID:             subscription.PlanID,
			PlanName:           subscription.Plan.Name,
			Price:              subscription.Plan.Price,
			BillingPeriod:      subscription.Plan.BillingPeriod,
			Status:             subscription.Status,
			CurrentPeriodStart: subscription.CurrentPeriodStart,
			CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
			NextBillingDate:    subscription.NextBillingDate,
			AutoRenew:          subscription.AutoRenew,
			CreatedAt:          subscription.CreatedAt,
		},
	})
}

// GetPaymentHistory gets payment history
func GetPaymentHistory(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var payments []models.Payment
	if err := database.DB.Where("company_id = ?", companyID).Order("payment_date DESC").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to fetch payment history",
		})
		return
	}

	var paymentDTOs []dto.PaymentRecordDTO
	for _, payment := range payments {
		paymentDTOs = append(paymentDTOs, dto.PaymentRecordDTO{
			ID:            payment.ID,
			PaymentDate:   payment.PaymentDate,
			CreatedAt:     payment.CreatedAt,
			InvoiceNumber: payment.InvoiceNumber,
			Amount:        payment.Amount,
			Status:        payment.Status,
			InvoiceUrl:    payment.InvoiceUrl,
		})
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Payment history retrieved",
		Data: map[string]interface{}{
			"payments": paymentDTOs,
		},
	})
}

// RenewSubscription renews an expiring subscription
func RenewSubscription(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var req struct {
		SubscriptionID uint `json:"subscriptionId" binding:"required"`
		AutoRenew      bool `json:"autoRenew"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Find subscription
	var subscription models.Subscription
	if err := database.DB.Preload("Plan").Where("id = ? AND company_id = ?", req.SubscriptionID, uint(companyID)).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Subscription not found",
		})
		return
	}

	// Calculate new period
	newPeriodStart := subscription.CurrentPeriodEnd
	var newPeriodEnd time.Time
	if subscription.Plan.BillingPeriod == "monthly" {
		newPeriodEnd = newPeriodStart.AddDate(0, 1, 0)
	} else {
		newPeriodEnd = newPeriodStart.AddDate(1, 0, 0)
	}

	// Update subscription
	updateData := map[string]interface{}{
		"status":                "active",
		"current_period_start": newPeriodStart,
		"current_period_end":   newPeriodEnd,
		"next_billing_date":    newPeriodEnd,
		"auto_renew":           req.AutoRenew,
	}

	if err := database.DB.Model(&subscription).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to renew subscription",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Subscription renewed successfully",
		Data: map[string]interface{}{
			"currentPeriodEnd": newPeriodEnd,
			"nextBillingDate":  newPeriodEnd,
			"status":           "active",
		},
	})
}

// CancelSubscription cancels a subscription
func CancelSubscription(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var req struct {
		SubscriptionID uint `json:"subscriptionId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Find subscription
	var subscription models.Subscription
	if err := database.DB.Where("id = ? AND company_id = ?", req.SubscriptionID, uint(companyID)).First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Subscription not found",
		})
		return
	}

	// Update subscription status
	now := time.Now()
	if err := database.DB.Model(&subscription).Updates(map[string]interface{}{
		"status":       "cancelled",
		"cancelled_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to cancel subscription",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Subscription cancelled successfully",
		Data: map[string]interface{}{
			"status":     "cancelled",
			"cancelledAt": now,
		},
	})
}

// UpgradeSubscription upgrades to a new plan
func UpgradeSubscription(c *gin.Context) {
	// Get company_id from context (set by AuthMiddleware)
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok || companyID == 0 {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized - company ID not found in token",
		})
		return
	}

	var req struct {
		NewPlanID uint `json:"newPlanId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Find current subscription
	var subscription models.Subscription
	if err := database.DB.Where("company_id = ? AND status = ?", uint(companyID), "active").First(&subscription).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "No active subscription found",
		})
		return
	}

	// Verify new plan exists
	var newPlan models.SubscriptionPlan
	if err := database.DB.First(&newPlan, req.NewPlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Plan not found",
		})
		return
	}

	// Update subscription
	if err := database.DB.Model(&subscription).Update("plan_id", req.NewPlanID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to upgrade subscription",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Subscription upgraded successfully",
		Data: map[string]interface{}{
			"newPlanId":  newPlan.ID,
			"newPlanName": newPlan.Name,
		},
	})
}

// CreateSubscriptionForCompany creates a subscription for a company (for testing/seeding)
func CreateSubscriptionForCompany(c *gin.Context) {
	type Request struct {
		CompanyID uint `json:"companyId" binding:"required"`
		PlanID    uint `json:"planId"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Verify company exists
	var company models.Company
	if err := database.DB.First(&company, req.CompanyID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Company not found",
		})
		return
	}

	planID := req.PlanID
	if planID == 0 {
		planID = 1 // Default to trial plan
	}

	// Verify plan exists
	var plan models.SubscriptionPlan
	if err := database.DB.First(&plan, planID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Plan not found",
		})
		return
	}

	// Check if subscription already exists
	var existingSub models.Subscription
	if database.DB.Where("company_id = ?", req.CompanyID).First(&existingSub).Error == nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Message: "Subscription already exists for this company",
		})
		return
	}

	// Create subscription
	subscription := models.Subscription{
		CompanyID:          req.CompanyID,
		PlanID:             planID,
		Status:             "active",
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   time.Now().AddDate(0, 0, 10),
		NextBillingDate:    time.Now().AddDate(0, 0, 10),
		AutoRenew:          true,
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create subscription",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Message: "Subscription created successfully",
		Data: map[string]interface{}{
			"id":                  subscription.ID,
			"companyId":           subscription.CompanyID,
			"planId":              subscription.PlanID,
			"status":              subscription.Status,
			"currentPeriodStart":  subscription.CurrentPeriodStart,
			"currentPeriodEnd":    subscription.CurrentPeriodEnd,
			"nextBillingDate":     subscription.NextBillingDate,
			"autoRenew":           subscription.AutoRenew,
		},
	})
}
