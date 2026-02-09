package models

import (
	"time"

	"gorm.io/gorm"
)

type ShippingAddress struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID   *uint          `json:"customerId,omitempty" gorm:"column:customer_id"`
	Customer     *Customer      `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	FullName     string         `json:"fullName" gorm:"column:full_name;not null"`
	Phone        string         `json:"phone" gorm:"not null"`
	Email        string         `json:"email,omitempty"`
	AddressLine1 string         `json:"addressLine1" gorm:"column:address_line1;not null"`
	AddressLine2 string         `json:"addressLine2,omitempty" gorm:"column:address_line2"`
	City         string         `json:"city" gorm:"not null"`
	State        string         `json:"state" gorm:"not null"`
	PostalCode   string         `json:"postalCode" gorm:"column:postal_code;not null"`
	Country      string         `json:"country" gorm:"not null;default:'Bangladesh'"`
	IsDefault    bool           `json:"isDefault" gorm:"column:is_default;default:false"`
	AddressType  string         `json:"addressType" gorm:"column:address_type;type:enum('home','office','other');default:'home'"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ShippingAddress) TableName() string { return "shipping_addresses" }
