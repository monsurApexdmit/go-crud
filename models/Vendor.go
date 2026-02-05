package models

import (
	"time"

	"gorm.io/gorm"
)

type Vendor struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID        *uint          `json:"userId" gorm:"column:user_id"`
	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name          string         `json:"name" gorm:"not null"`
	Email         string         `json:"email" gorm:"uniqueIndex;not null"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	Logo          string         `json:"logo"`
	Status        string         `json:"status" gorm:"type:enum('Active','Inactive','Blocked');default:Active"`
	Description   string         `json:"description" gorm:"type:text"`
	TotalPaid     float64        `json:"totalPaid" gorm:"column:total_paid;default:0"`
	AmountPayable float64        `json:"amountPayable" gorm:"column:amount_payable;default:0"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
