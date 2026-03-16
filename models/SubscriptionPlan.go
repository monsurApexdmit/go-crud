package models

import "time"

// SubscriptionPlan represents a subscription plan
type SubscriptionPlan struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"index"`
	Description     string    `json:"description"`
	Price           int64     `json:"price"` // in cents (e.g., 7900 = $79.00)
	BillingPeriod   string    `json:"billingPeriod"` // monthly, yearly
	MaxUsers        int       `json:"maxUsers"`
	MaxProducts     int       `json:"maxProducts"`
	MaxBranches     int       `json:"maxBranches"`
	Features        string    `json:"features"` // JSON encoded features list
	IsActive        bool      `json:"isActive" gorm:"default:true"`
	IsFeatured      bool      `json:"isFeatured" gorm:"default:false"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relations
	Subscriptions []Subscription `json:"subscriptions,omitempty" gorm:"foreignKey:PlanID"`
}

// TableName specifies the table name for the SubscriptionPlan model
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}
