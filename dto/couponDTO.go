package dto

import "time"

// CreateCouponRequest validates coupon creation
type CreateCouponRequest struct {
	CampaignName string    `json:"campaign_name" binding:"required,min=3,max=200"`
	Code         string    `json:"code" binding:"required,min=3,max=50,alphanum"`
	Discount     float64   `json:"discount" binding:"required,gt=0"`
	Type         string    `json:"type" binding:"required,oneof=percentage fixed free_shipping"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required"`
	Status       bool      `json:"status"`

	// Optional fields
	UsageLimit        *int    `json:"usage_limit" binding:"omitempty,gt=0"`
	UsageLimitPerUser *int    `json:"usage_limit_per_user" binding:"omitempty,gt=0"`
	MinOrderAmount    float64 `json:"min_order_amount" binding:"omitempty,gte=0"`
	MaxDiscount       *float64 `json:"max_discount" binding:"omitempty,gt=0"`

	ApplicableToCategories string `json:"applicable_to_categories" binding:"omitempty"`
	ApplicableToProducts   string `json:"applicable_to_products" binding:"omitempty"`
	FreeShipping           bool   `json:"free_shipping"`

	Stackable  bool `json:"stackable"`
	AutoApply  bool `json:"auto_apply"`
	Priority   int  `json:"priority" binding:"omitempty,gte=0"`
}

// UpdateCouponRequest validates coupon updates
type UpdateCouponRequest struct {
	CampaignName *string    `json:"campaign_name" binding:"omitempty,min=3,max=200"`
	Code         *string    `json:"code" binding:"omitempty,min=3,max=50,alphanum"`
	Discount     *float64   `json:"discount" binding:"omitempty,gt=0"`
	Type         *string    `json:"type" binding:"omitempty,oneof=percentage fixed free_shipping"`
	StartDate    *time.Time `json:"start_date" binding:"omitempty"`
	EndDate      *time.Time `json:"end_date" binding:"omitempty"`
	Status       *bool      `json:"status"`

	UsageLimit        *int     `json:"usage_limit" binding:"omitempty,gt=0"`
	UsageLimitPerUser *int     `json:"usage_limit_per_user" binding:"omitempty,gt=0"`
	MinOrderAmount    *float64 `json:"min_order_amount" binding:"omitempty,gte=0"`
	MaxDiscount       *float64 `json:"max_discount" binding:"omitempty,gt=0"`

	ApplicableToCategories *string `json:"applicable_to_categories"`
	ApplicableToProducts   *string `json:"applicable_to_products"`
	FreeShipping           *bool   `json:"free_shipping"`

	Stackable *bool `json:"stackable"`
	AutoApply *bool `json:"auto_apply"`
	Priority  *int  `json:"priority" binding:"omitempty,gte=0"`
}

// ValidateCouponRequest validates a coupon code for checkout
type ValidateCouponRequest struct {
	Code        string     `json:"code" binding:"required"`
	CustomerID  *uint      `json:"customer_id"`
	OrderAmount float64    `json:"order_amount" binding:"required,gt=0"`
	CartItems   []CartItem `json:"cart_items"`
}

// CartItem represents a product in the cart
type CartItem struct {
	ProductID  uint    `json:"product_id" binding:"required"`
	CategoryID uint    `json:"category_id"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Quantity   int     `json:"quantity" binding:"required,gt=0"`
}
