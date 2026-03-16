package models

import (
	"time"

	"gorm.io/gorm"
)

type ShippingAddress struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID    uint           `json:"companyId" gorm:"column:company_id;not null;index"`
	CustomerID   *uint          `json:"customerId,omitempty" gorm:"column:customer_id"`
	Customer     *Customer      `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	FullName     string         `json:"fullName" gorm:"column:full_name;not null;type:varchar(255)"`
	Phone        string         `json:"phone" gorm:"not null;type:varchar(50)"`
	Email        string         `json:"email,omitempty" gorm:"type:varchar(255)"`
	AddressLine1 string         `json:"addressLine1" gorm:"type:varchar(255);column:address_line1;not null"`
	AddressLine2 string         `json:"addressLine2,omitempty" gorm:"type:varchar(255);column:address_line2"`
	City         string         `json:"city" gorm:"type:varchar(100);not null"`
	State        string         `json:"state" gorm:"type:varchar(100);not null"`
	PostalCode   string         `json:"postalCode" gorm:"type:varchar(20);column:postal_code;not null"`
	Country      string         `json:"country" gorm:"type:varchar(100);not null;default:'Bangladesh'"`
	IsDefault    bool           `json:"isDefault" gorm:"column:is_default;default:false"`
	AddressType  string         `json:"addressType" gorm:"column:address_type;type:enum('home','office','other');default:'home'"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ShippingAddress) TableName() string { return "shipping_addresses" }
