package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
)

// ListSalaryPayments returns salary records with optional filters.
// Query params: staff_id, month, status, page, limit
func ListSalaryPayments(c *gin.Context) {
	query := database.DB.Preload("Staff").Model(&models.SalaryPayment{})

	if staffID := c.Query("staff_id"); staffID != "" {
		query = query.Where("staff_id = ?", staffID)
	}
	if month := c.Query("month"); month != "" {
		query = query.Where("month = ?", month)
	}
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
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

	var payments []models.SalaryPayment
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve salary payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Salary payments retrieved successfully",
		"data":    payments,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func GetSalaryPayment(c *gin.Context) {
	var payment models.SalaryPayment
	if err := database.DB.Preload("Staff").First(&payment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Salary payment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Salary payment fetched successfully", "data": payment})
}

// paymentRequest is the body for POST /salary-payments
type paymentRequest struct {
	StaffID       uint    `json:"staffId" binding:"required"`
	Month         string  `json:"month" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	PaidAmount    float64 `json:"paidAmount"`
	PaymentDate   string  `json:"paymentDate"`
	PaymentMethod string  `json:"paymentMethod"`
	Notes         string  `json:"notes"`
}

// CreateSalaryPayment creates a new salary record for a staff member + month.
// The unique constraint (staff_id, month) ensures one record per staff per month.
func CreateSalaryPayment(c *gin.Context) {
	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// verify staff exists
	var staff models.Staff
	if err := database.DB.First(&staff, req.StaffID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff not found"})
		return
	}

	payment := models.SalaryPayment{
		StaffID:       req.StaffID,
		Month:         req.Month,
		Amount:        req.Amount,
		PaidAmount:    req.PaidAmount,
		Status:        calcPaymentStatus(req.Amount, req.PaidAmount),
		PaymentDate:   req.PaymentDate,
		PaymentMethod: req.PaymentMethod,
		Notes:         req.Notes,
	}

	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create salary payment"})
		return
	}

	database.DB.Preload("Staff").First(&payment, payment.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Salary payment created successfully", "data": payment})
}

// UpdateSalaryPayment — used when marking a payment as paid / partial.
// Accepts paidAmount and recalculates status automatically.
func UpdateSalaryPayment(c *gin.Context) {
	var payment models.SalaryPayment
	if err := database.DB.First(&payment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Salary payment not found"})
		return
	}

	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	updates := map[string]interface{}{
		"paid_amount":    req.PaidAmount,
		"status":         calcPaymentStatus(payment.Amount, req.PaidAmount),
		"payment_method": req.PaymentMethod,
		"notes":          req.Notes,
	}
	if req.PaymentDate != "" {
		updates["payment_date"] = req.PaymentDate
	}

	if err := database.DB.Model(&payment).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update salary payment"})
		return
	}

	database.DB.Preload("Staff").First(&payment, payment.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Salary payment updated successfully", "data": payment})
}

func DeleteSalaryPayment(c *gin.Context) {
	var payment models.SalaryPayment
	if err := database.DB.First(&payment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Salary payment not found"})
		return
	}

	database.DB.Delete(&payment)
	c.JSON(http.StatusOK, gin.H{"message": "Salary payment deleted successfully"})
}

// calcPaymentStatus returns Paid / Partial / Pending based on amounts.
func calcPaymentStatus(amount, paidAmount float64) string {
	if paidAmount >= amount {
		return "Paid"
	}
	if paidAmount > 0 {
		return "Partial"
	}
	return "Pending"
}
