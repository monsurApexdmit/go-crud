package models

import "time"

// CompanySettings represents company-level settings
type CompanySettings struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CompanyID       uint      `json:"companyId" gorm:"uniqueIndex"`
	CompanyName     string    `json:"companyName"`
	TaxID           string    `json:"taxId,omitempty"`
	TaxIDType       string    `json:"taxIdType,omitempty"` // vat, ein, gst, other
	TaxRate         float64   `json:"taxRate" gorm:"default:0"` // percentage
	Currency        string    `json:"currency" gorm:"default:'USD'"`
	Timezone        string    `json:"timezone" gorm:"default:'UTC'"`
	Language        string    `json:"language" gorm:"default:'en'"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relations
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name for the CompanySettings model
func (CompanySettings) TableName() string {
	return "company_settings"
}
