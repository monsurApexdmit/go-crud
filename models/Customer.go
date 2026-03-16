package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID    uint           `json:"companyId" gorm:"column:company_id;not null;index"`
	UserID       *uint          `json:"userId" gorm:"column:user_id"`
	User         *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name         string         `json:"name" gorm:"not null;type:varchar(255)"`
	Email        string         `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_customer_company_email,composite:company_id;not null"`
	Phone        string         `json:"phone" gorm:"type:varchar(50)"`
	Address      string         `json:"address"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	ZipCode      string         `json:"zipCode" gorm:"column:zip_code"`
	Country      string         `json:"country"`
	CustomerType string         `json:"customerType" gorm:"column:customer_type;type:enum('retail','wholesale');default:retail"`
	Status       string         `json:"status" gorm:"type:enum('active','inactive');default:active"`
	Notes        string         `json:"notes" gorm:"type:text"`
	StoreCredit  float64        `json:"storeCredit" gorm:"column:store_credit;default:0"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
