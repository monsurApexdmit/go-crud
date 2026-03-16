package models

import "time"

// Company represents a SaaS company/tenant
type Company struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"index"`
	Industry    string `json:"industry"`
	Website     string `json:"website"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city"`
	State       string `json:"state"`
	ZipCode     string `json:"zipCode"`
	Country     string `json:"country"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'trial'"` // trial, active, expired, suspended
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Relations
	Users           []SaasUser          `json:"users,omitempty" gorm:"foreignKey:CompanyID"`
	Subscriptions   []Subscription      `json:"subscriptions,omitempty" gorm:"foreignKey:CompanyID"`
	BillingContact  *BillingContact     `json:"billingContact,omitempty" gorm:"foreignKey:CompanyID"`
	CompanySettings *CompanySettings    `json:"settings,omitempty" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name for the Company model
func (Company) TableName() string {
	return "companies"
}
