package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID    uint           `json:"companyId" gorm:"column:company_id;not null;index"`
	CampaignName string         `json:"campaign_name" form:"campaign_name" gorm:"type:varchar(255)"`
	Code         string         `json:"code" gorm:"type:varchar(100);uniqueIndex:idx_coupon_company_code,composite:company_id" form:"code"`
	Discount     float64        `json:"discount" form:"discount"`
	Type         string         `json:"type" form:"type" gorm:"type:varchar(50)"` // percentage, fixed, free_shipping
	StartDate    time.Time      `json:"start_date" gorm:"type:timestamp"`
	EndDate      time.Time      `json:"end_date" gorm:"type:timestamp"`
	Status       bool           `json:"status" gorm:"default:false" form:"status"`
	Image        string         `json:"image"`

	// Usage Limits
	UsageLimit        *int    `json:"usage_limit" gorm:"comment:Total times coupon can be used (null = unlimited)" form:"usage_limit"`
	UsageLimitPerUser *int    `json:"usage_limit_per_user" gorm:"comment:Times per customer (null = unlimited)" form:"usage_limit_per_user"`
	TimesUsed         int     `json:"times_used" gorm:"default:0;comment:Total redemptions"`

	// Order Requirements
	MinOrderAmount float64  `json:"min_order_amount" gorm:"default:0;comment:Minimum order value required" form:"min_order_amount"`
	MaxDiscount    *float64 `json:"max_discount" gorm:"comment:Max discount cap for percentage type" form:"max_discount"`

	// Applicability Rules
	ApplicableToCategories string `json:"applicable_to_categories" gorm:"type:text;comment:Comma-separated category IDs (empty = all)" form:"applicable_to_categories"`
	ApplicableToProducts   string `json:"applicable_to_products" gorm:"type:text;comment:Comma-separated product IDs (empty = all)" form:"applicable_to_products"`
	FreeShipping           bool   `json:"free_shipping" gorm:"default:false;comment:Whether this coupon grants free shipping" form:"free_shipping"`

	// Advanced Features
	Stackable  bool `json:"stackable" gorm:"default:false;comment:Can be combined with other coupons" form:"stackable"`
	AutoApply  bool `json:"auto_apply" gorm:"default:false;comment:Automatically apply if conditions met" form:"auto_apply"`
	Priority   int  `json:"priority" gorm:"default:0;comment:Higher priority applies first" form:"priority"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
