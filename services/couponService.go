package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-crud/database"
	"go-crud/models"

	"gorm.io/gorm"
)

// CouponValidationRequest represents the request to validate a coupon
type CouponValidationRequest struct {
	Code        string     `json:"code" binding:"required"`
	CustomerID  *uint      `json:"customer_id"`
	OrderAmount float64    `json:"order_amount" binding:"required,gt=0"`
	CartItems   []CartItem `json:"cart_items"`
}

// CartItem represents a product in the cart
type CartItem struct {
	ProductID  uint    `json:"product_id"`
	CategoryID uint    `json:"category_id"`
	Price      float64 `json:"price"`
	Quantity   int     `json:"quantity"`
}

// CouponValidationResult represents the validation result
type CouponValidationResult struct {
	Valid          bool           `json:"valid"`
	Coupon         *models.Coupon `json:"coupon,omitempty"`
	DiscountAmount float64        `json:"discount_amount"`
	FinalAmount    float64        `json:"final_amount"`
	FreeShipping   bool           `json:"free_shipping"`
	ErrorCode      string         `json:"error_code,omitempty"`
	Message        string         `json:"message"`
}

// ValidateCoupon validates a coupon code against order details
func ValidateCoupon(req CouponValidationRequest) (*CouponValidationResult, error) {
	var coupon models.Coupon

	// 1. Check if coupon exists
	if err := database.DB.Where("code = ?", req.Code).First(&coupon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &CouponValidationResult{
				Valid:     false,
				ErrorCode: "NOT_FOUND",
				Message:   "Coupon code not found",
			}, nil
		}
		return nil, err
	}

	// 2. Check if coupon is active
	if !coupon.Status {
		return &CouponValidationResult{
			Valid:     false,
			ErrorCode: "INACTIVE",
			Message:   "This coupon is currently inactive",
		}, nil
	}

	// 3. Check date range
	now := time.Now()
	if now.Before(coupon.StartDate) {
		return &CouponValidationResult{
			Valid:     false,
			ErrorCode: "NOT_STARTED",
			Message:   fmt.Sprintf("This coupon is not valid until %s", coupon.StartDate.Format("2006-01-02")),
		}, nil
	}
	if now.After(coupon.EndDate) {
		return &CouponValidationResult{
			Valid:     false,
			ErrorCode: "EXPIRED",
			Message:   fmt.Sprintf("This coupon expired on %s", coupon.EndDate.Format("2006-01-02")),
		}, nil
	}

	// 4. Check global usage limit
	if coupon.UsageLimit != nil && coupon.TimesUsed >= *coupon.UsageLimit {
		return &CouponValidationResult{
			Valid:     false,
			ErrorCode: "USAGE_LIMIT_REACHED",
			Message:   "This coupon has reached its maximum usage limit",
		}, nil
	}

	// 5. Check per-user usage limit
	if req.CustomerID != nil && coupon.UsageLimitPerUser != nil {
		var usageCount int64
		database.DB.Model(&models.CouponUsage{}).
			Where("coupon_id = ? AND customer_id = ?", coupon.ID, *req.CustomerID).
			Count(&usageCount)

		if int(usageCount) >= *coupon.UsageLimitPerUser {
			return &CouponValidationResult{
				Valid:     false,
				ErrorCode: "USER_LIMIT_REACHED",
				Message:   fmt.Sprintf("You have already used this coupon the maximum allowed times (%d)", *coupon.UsageLimitPerUser),
			}, nil
		}
	}

	// 6. Check minimum order amount
	if req.OrderAmount < coupon.MinOrderAmount {
		return &CouponValidationResult{
			Valid:     false,
			ErrorCode: "MIN_ORDER_NOT_MET",
			Message:   fmt.Sprintf("Minimum order amount of $%.2f required to use this coupon", coupon.MinOrderAmount),
		}, nil
	}

	// 7. Check category restrictions
	if coupon.ApplicableToCategories != "" && len(req.CartItems) > 0 {
		allowedCategories := parseIDList(coupon.ApplicableToCategories)
		hasApplicableItem := false

		for _, item := range req.CartItems {
			if contains(allowedCategories, item.CategoryID) {
				hasApplicableItem = true
				break
			}
		}

		if !hasApplicableItem {
			return &CouponValidationResult{
				Valid:     false,
				ErrorCode: "NOT_APPLICABLE",
				Message:   "This coupon is not applicable to items in your cart",
			}, nil
		}
	}

	// 8. Check product restrictions
	if coupon.ApplicableToProducts != "" && len(req.CartItems) > 0 {
		allowedProducts := parseIDList(coupon.ApplicableToProducts)
		hasApplicableItem := false

		for _, item := range req.CartItems {
			if contains(allowedProducts, item.ProductID) {
				hasApplicableItem = true
				break
			}
		}

		if !hasApplicableItem {
			return &CouponValidationResult{
				Valid:     false,
				ErrorCode: "NOT_APPLICABLE",
				Message:   "This coupon is not applicable to items in your cart",
			}, nil
		}
	}

	// 9. Calculate discount
	discountAmount := CalculateDiscount(&coupon, req.OrderAmount)
	finalAmount := req.OrderAmount - discountAmount

	return &CouponValidationResult{
		Valid:          true,
		Coupon:         &coupon,
		DiscountAmount: discountAmount,
		FinalAmount:    finalAmount,
		FreeShipping:   coupon.FreeShipping,
		Message:        fmt.Sprintf("Coupon applied successfully! You saved $%.2f", discountAmount),
	}, nil
}

// CalculateDiscount calculates the actual discount amount
func CalculateDiscount(coupon *models.Coupon, orderAmount float64) float64 {
	var discount float64

	switch coupon.Type {
	case "percentage":
		discount = orderAmount * (coupon.Discount / 100)
		// Apply max discount cap if set
		if coupon.MaxDiscount != nil && discount > *coupon.MaxDiscount {
			discount = *coupon.MaxDiscount
		}
	case "fixed":
		discount = coupon.Discount
		// Don't exceed order amount
		if discount > orderAmount {
			discount = orderAmount
		}
	case "free_shipping":
		discount = 0 // Discount is handled separately via FreeShipping flag
	default:
		discount = 0
	}

	return discount
}

// ApplyCouponToOrder applies a coupon to an order and tracks usage
func ApplyCouponToOrder(sellID uint, couponCode string, customerID *uint, orderAmount float64, discountApplied float64) error {
	var coupon models.Coupon

	// Get coupon
	if err := database.DB.Where("code = ?", couponCode).First(&coupon).Error; err != nil {
		return err
	}

	// Create usage record
	usage := models.CouponUsage{
		CouponID:        coupon.ID,
		CustomerID:      customerID,
		SellID:          sellID,
		CouponCode:      couponCode,
		DiscountApplied: discountApplied,
		OriginalAmount:  orderAmount + discountApplied,
		FinalAmount:     orderAmount,
	}

	if err := database.DB.Create(&usage).Error; err != nil {
		return err
	}

	// Increment times_used
	if err := database.DB.Model(&coupon).Update("times_used", gorm.Expr("times_used + 1")).Error; err != nil {
		return err
	}

	return nil
}

// GetCouponByCode retrieves a coupon by its code (public endpoint)
func GetCouponByCode(code string) (*models.Coupon, error) {
	var coupon models.Coupon

	if err := database.DB.Where("code = ? AND status = ?", code, true).First(&coupon).Error; err != nil {
		return nil, err
	}

	return &coupon, nil
}

// GetCouponUsageStats retrieves usage statistics for a coupon
func GetCouponUsageStats(couponID uint) (map[string]interface{}, error) {
	var coupon models.Coupon
	if err := database.DB.First(&coupon, couponID).Error; err != nil {
		return nil, err
	}

	var usages []models.CouponUsage
	if err := database.DB.
		Preload("Customer").
		Preload("Sell").
		Where("coupon_id = ?", couponID).
		Order("used_at DESC").
		Limit(10).
		Find(&usages).Error; err != nil {
		return nil, err
	}

	// Calculate analytics
	var totalDiscount, totalRevenue float64
	var uniqueCustomers int64

	database.DB.Model(&models.CouponUsage{}).
		Where("coupon_id = ?", couponID).
		Select("COALESCE(SUM(discount_applied), 0)").
		Scan(&totalDiscount)

	database.DB.Model(&models.CouponUsage{}).
		Where("coupon_id = ?", couponID).
		Select("COALESCE(SUM(final_amount), 0)").
		Scan(&totalRevenue)

	database.DB.Model(&models.CouponUsage{}).
		Where("coupon_id = ? AND customer_id IS NOT NULL", couponID).
		Distinct("customer_id").
		Count(&uniqueCustomers)

	avgOrderValue := 0.0
	if coupon.TimesUsed > 0 {
		avgOrderValue = totalRevenue / float64(coupon.TimesUsed)
	}

	stats := map[string]interface{}{
		"coupon": map[string]interface{}{
			"id":            coupon.ID,
			"code":          coupon.Code,
			"campaign_name": coupon.CampaignName,
		},
		"analytics": map[string]interface{}{
			"total_uses":          coupon.TimesUsed,
			"total_discount_given": totalDiscount,
			"total_revenue":       totalRevenue,
			"unique_customers":    uniqueCustomers,
			"avg_order_value":     avgOrderValue,
		},
		"recent_usage": usages,
	}

	return stats, nil
}

// Helper functions

func parseIDList(idList string) []uint {
	if idList == "" {
		return []uint{}
	}

	parts := strings.Split(idList, ",")
	ids := make([]uint, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if id, err := strconv.ParseUint(part, 10, 64); err == nil {
			ids = append(ids, uint(id))
		}
	}

	return ids
}

func contains(slice []uint, item uint) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
