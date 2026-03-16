package models

import "time"

// SaasUser represents a user in the SaaS platform
type SaasUser struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CompanyID   uint      `json:"companyId" gorm:"index"`
	Email       string    `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_company_email"`
	FullName    string    `json:"fullName"`
	Password    string    `json:"-" gorm:"column:password"`
	Role        string    `json:"role" gorm:"default:'staff'"` // owner, admin, manager, staff
	Status      string    `json:"status" gorm:"default:'active'"` // active, invited, inactive
	JoinedDate  time.Time `json:"joinedDate"`
	LastLogin   *time.Time `json:"lastLogin,omitempty"`
	Avatar      string    `json:"avatar,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Relations
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name for the SaasUser model
func (SaasUser) TableName() string {
	return "saas_users"
}
