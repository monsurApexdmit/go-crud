package models

import (
	"time"

	"gorm.io/gorm"
)

type Sell struct {
	ID                 uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	InvoiceNo          string           `json:"invoiceNo" gorm:"column:invoice_no;uniqueIndex;not null"`
	OrderTime          time.Time        `json:"orderTime" gorm:"column:order_time;not null"`
	CustomerID         *uint            `json:"customerId,omitempty" gorm:"column:customer_id"`
	Customer           *Customer        `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CustomerName       string           `json:"customerName" gorm:"column:customer_name;not null"`
	ShippingAddressID  *uint            `json:"shippingAddressId,omitempty" gorm:"column:shipping_address_id"`
	ShippingAddress    *ShippingAddress `json:"shippingAddress,omitempty" gorm:"foreignKey:ShippingAddressID"`

	// Shipping Address Snapshot (stored directly with order)
	ShippingFullName      string `json:"shippingFullName,omitempty" gorm:"column:shipping_full_name"`
	ShippingPhone         string `json:"shippingPhone,omitempty" gorm:"column:shipping_phone"`
	ShippingEmail         string `json:"shippingEmail,omitempty" gorm:"column:shipping_email"`
	ShippingAddressLine1  string `json:"shippingAddressLine1,omitempty" gorm:"column:shipping_address_line1"`
	ShippingAddressLine2  string `json:"shippingAddressLine2,omitempty" gorm:"column:shipping_address_line2"`
	ShippingCity          string `json:"shippingCity,omitempty" gorm:"column:shipping_city"`
	ShippingState         string `json:"shippingState,omitempty" gorm:"column:shipping_state"`
	ShippingPostalCode    string `json:"shippingPostalCode,omitempty" gorm:"column:shipping_postal_code"`
	ShippingCountry       string `json:"shippingCountry,omitempty" gorm:"column:shipping_country"`
	ShippingAddressType   string `json:"shippingAddressType,omitempty" gorm:"column:shipping_address_type"`

	Method             string           `json:"method" gorm:"not null;default:'Cash'"`
	Amount             float64          `json:"amount" gorm:"not null;default:0"`
	ShippingCost       float64          `json:"shippingCost" gorm:"column:shipping_cost;default:0"`
	ShippingMethod     string           `json:"shippingMethod,omitempty" gorm:"column:shipping_method"`
	Discount           float64          `json:"discount" gorm:"default:0"`
	Status             string           `json:"status" gorm:"not null;default:'Pending'"`
	PaymentStatus      string           `json:"paymentStatus,omitempty" gorm:"column:payment_status;type:enum('pending','paid','partially_paid','refunded','failed');default:'pending'"`
	FulfillmentStatus  string           `json:"fulfillmentStatus,omitempty" gorm:"column:fulfillment_status;type:enum('unfulfilled','processing','shipped','delivered','cancelled');default:'unfulfilled'"`
	TrackingNumber     string           `json:"trackingNumber,omitempty" gorm:"column:tracking_number"`
	Carrier            string           `json:"carrier,omitempty"`
	ShippedAt          *time.Time       `json:"shippedAt,omitempty" gorm:"column:shipped_at"`
	DeliveredAt        *time.Time       `json:"deliveredAt,omitempty" gorm:"column:delivered_at"`
	Notes              string           `json:"notes,omitempty" gorm:"type:text"`
	Items              []OrderItem      `json:"items,omitempty" gorm:"foreignKey:SellID"`
	Shipments          []OrderShipment  `json:"shipments,omitempty" gorm:"foreignKey:SellID"`
	CreatedAt          time.Time        `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          time.Time        `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt   `json:"-" gorm:"index"`
}

func (Sell) TableName() string { return "sells" }
