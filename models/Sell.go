package models

import (
	"time"

	"gorm.io/gorm"
)

type Sell struct {
	ID                uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID         uint             `json:"companyId" gorm:"column:company_id;not null;index"`
	InvoiceNo         string           `json:"invoiceNo" gorm:"type:varchar(100);column:invoice_no;uniqueIndex;not null"`
	OrderTime         time.Time        `json:"orderTime" gorm:"column:order_time;not null"`
	CustomerID        *uint            `json:"customerId,omitempty" gorm:"column:customer_id"`
	Customer          *Customer        `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CustomerName      string           `json:"customerName" gorm:"type:varchar(255);column:customer_name;not null"`
	ShippingAddressID *uint            `json:"shippingAddressId,omitempty" gorm:"column:shipping_address_id"`
	ShippingAddress   *ShippingAddress `json:"shippingAddress,omitempty" gorm:"foreignKey:ShippingAddressID"`

	// Shipping Address Snapshot (stored directly with order)
	ShippingFullName     string `json:"shippingFullName,omitempty" gorm:"type:varchar(255);column:shipping_full_name"`
	ShippingPhone        string `json:"shippingPhone,omitempty" gorm:"type:varchar(50);column:shipping_phone"`
	ShippingEmail        string `json:"shippingEmail,omitempty" gorm:"type:varchar(255);column:shipping_email"`
	ShippingAddressLine1 string `json:"shippingAddressLine1,omitempty" gorm:"type:varchar(255);column:shipping_address_line1"`
	ShippingAddressLine2 string `json:"shippingAddressLine2,omitempty" gorm:"type:varchar(255);column:shipping_address_line2"`
	ShippingCity         string `json:"shippingCity,omitempty" gorm:"type:varchar(100);column:shipping_city"`
	ShippingState        string `json:"shippingState,omitempty" gorm:"type:varchar(100);column:shipping_state"`
	ShippingPostalCode   string `json:"shippingPostalCode,omitempty" gorm:"type:varchar(20);column:shipping_postal_code"`
	ShippingCountry      string `json:"shippingCountry,omitempty" gorm:"type:varchar(100);column:shipping_country"`
	ShippingAddressType  string `json:"shippingAddressType,omitempty" gorm:"type:varchar(50);column:shipping_address_type"`

	Method         string  `json:"method" gorm:"type:varchar(100);not null;default:'Cash'"`
	Amount         float64 `json:"amount" gorm:"not null;default:0"`
	ShippingCost   float64 `json:"shippingCost" gorm:"column:shipping_cost;default:0"`
	ShippingMethod string  `json:"shippingMethod,omitempty" gorm:"type:varchar(100);column:shipping_method"`

	// Coupon Integration
	CouponID   *uint   `json:"couponId,omitempty" gorm:"column:coupon_id;index"`
	Coupon     *Coupon `json:"coupon,omitempty" gorm:"foreignKey:CouponID"`
	CouponCode string  `json:"couponCode,omitempty" gorm:"column:coupon_code"`
	Discount   float64 `json:"discount" gorm:"default:0"`

	Status            string          `json:"status" gorm:"type:varchar(50);not null;default:'Pending'"`
	StockDeducted     bool            `json:"stockDeducted,omitempty" gorm:"column:stock_deducted;not null;default:false"`
	PaymentStatus     string          `json:"paymentStatus,omitempty" gorm:"column:payment_status;type:enum('pending','paid','partially_paid','refunded','failed');default:'pending'"`
	FulfillmentStatus string          `json:"fulfillmentStatus,omitempty" gorm:"column:fulfillment_status;type:enum('unfulfilled','processing','shipped','delivered','cancelled');default:'unfulfilled'"`
	TrackingNumber    string          `json:"trackingNumber,omitempty" gorm:"type:varchar(100);column:tracking_number"`
	Carrier           string          `json:"carrier,omitempty" gorm:"type:varchar(100)"`
	ShippedAt         *time.Time      `json:"shippedAt,omitempty" gorm:"column:shipped_at"`
	DeliveredAt       *time.Time      `json:"deliveredAt,omitempty" gorm:"column:delivered_at"`
	Notes             string          `json:"notes,omitempty" gorm:"type:text"`
	Items             []OrderItem     `json:"items,omitempty" gorm:"foreignKey:SellID"`
	Shipments         []OrderShipment `json:"shipments,omitempty" gorm:"foreignKey:SellID"`
	CreatedAt         time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt  `json:"-" gorm:"index"`
}

func (Sell) TableName() string { return "sells" }
