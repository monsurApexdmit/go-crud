package models

import "time"

// Subscription represents a company's subscription
type Subscription struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CompanyID       uint      `json:"companyId" gorm:"index"`
	PlanID          uint      `json:"planId" gorm:"index"`
	Status          string    `json:"status" gorm:"default:'active'"` // active, expired, cancelled, pending
	CurrentPeriodStart time.Time `json:"currentPeriodStart"`
	CurrentPeriodEnd time.Time `json:"currentPeriodEnd"`
	NextBillingDate time.Time `json:"nextBillingDate"`
	AutoRenew       bool      `json:"autoRenew" gorm:"default:true"`
	StripeSubscriptionID string `json:"stripeSubscriptionId,omitempty"`
	CancelledAt     *time.Time `json:"cancelledAt,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relations
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Plan    *SubscriptionPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Payments []Payment `json:"payments,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// TableName specifies the table name for the Subscription model
func (Subscription) TableName() string {
	return "subscriptions"
}
