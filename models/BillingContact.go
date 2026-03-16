package models

import "time"

// BillingContact represents billing address and tax information
type BillingContact struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CompanyID   uint      `json:"companyId" gorm:"uniqueIndex"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	ZipCode     string    `json:"zipCode"`
	Country     string    `json:"country"`
	TaxID       string    `json:"taxId,omitempty"`
	TaxIDType   string    `json:"taxIdType,omitempty"` // vat, ein, gst, other
	IsDefault   bool      `json:"isDefault" gorm:"default:true"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Relations
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name for the BillingContact model
func (BillingContact) TableName() string {
	return "billing_contacts"
}
