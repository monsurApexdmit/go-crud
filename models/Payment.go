package models

import "time"

// Payment represents a payment record
type Payment struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	SubscriptionID  uint      `json:"subscriptionId" gorm:"index"`
	CompanyID       uint      `json:"companyId" gorm:"index"`
	Amount          int64     `json:"amount"` // in cents
	Status          string    `json:"status"` // paid, pending, failed
	PaymentMethod   string    `json:"paymentMethod"` // credit_card, bank_transfer, etc.
	PaymentDate     *time.Time `json:"paymentDate,omitempty"`
	InvoiceNumber   string    `json:"invoiceNumber"`
	InvoiceUrl      string    `json:"invoiceUrl,omitempty"`
	StripePaymentID string    `json:"stripePaymentId,omitempty"`
	Description     string    `json:"description,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relations
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	Company      *Company      `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name for the Payment model
func (Payment) TableName() string {
	return "payments"
}
