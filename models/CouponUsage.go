package models

import "time"

type CouponUsage struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CouponID        uint      `json:"coupon_id" gorm:"not null;index"`
	Coupon          *Coupon   `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	CustomerID      *uint     `json:"customer_id" gorm:"index"`
	Customer        *Customer `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	SellID          uint      `json:"sell_id" gorm:"not null;index"`
	Sell            *Sell     `json:"sell,omitempty" gorm:"foreignKey:SellID"`

	// Snapshot of values at time of use
	CouponCode      string  `json:"coupon_code" gorm:"not null"`
	DiscountApplied float64 `json:"discount_applied" gorm:"not null;comment:Actual discount amount given"`
	OriginalAmount  float64 `json:"original_amount" gorm:"not null;comment:Order total before discount"`
	FinalAmount     float64 `json:"final_amount" gorm:"not null;comment:Order total after discount"`

	UsedAt    time.Time `json:"used_at" gorm:"autoCreateTime;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
