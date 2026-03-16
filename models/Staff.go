package models

import (
	"time"

	"gorm.io/gorm"
)

func (Staff) TableName() string { return "staff" }

type Staff struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID     uint           `json:"companyId" gorm:"column:company_id;not null;index"`
	UserID        *uint          `json:"userId" gorm:"column:user_id"`
	User          *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name          string         `json:"name" gorm:"not null;type:varchar(255)"`
	Email         string         `json:"email" gorm:"type:varchar(255);uniqueIndex:idx_staff_company_email,composite:company_id;not null"`
	Contact       string         `json:"contact" gorm:"type:varchar(50)"`
	JoiningDate   string         `json:"joiningDate" gorm:"column:joining_date"`
	Role          string         `json:"role"`
	Status        string         `json:"status" gorm:"type:enum('Active','Inactive');default:Active"`
	Published     bool           `json:"published" gorm:"default:false"`
	Avatar        string         `json:"avatar"`
	Salary        float64        `json:"salary" gorm:"default:0"`
	BankAccount   string         `json:"bankAccount" gorm:"column:bank_account"`
	PaymentMethod string         `json:"paymentMethod" gorm:"column:payment_method;type:enum('Bank Transfer','Cash','Check')"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
