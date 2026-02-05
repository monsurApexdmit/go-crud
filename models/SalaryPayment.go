package models

import "time"

type SalaryPayment struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StaffID       uint      `json:"staffId" gorm:"column:staff_id;not null"`
	Staff         Staff     `json:"staff" gorm:"foreignKey:StaffID"`
	Month         string    `json:"month" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"not null;default:0"`
	PaidAmount    float64   `json:"paidAmount" gorm:"column:paid_amount;default:0"`
	Status        string    `json:"status" gorm:"type:enum('Paid','Pending','Partial');default:Pending"`
	PaymentDate   string    `json:"paymentDate" gorm:"column:payment_date"`
	PaymentMethod string    `json:"paymentMethod" gorm:"column:payment_method"`
	Notes         string    `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (SalaryPayment) TableName() string { return "salary_payments" }
