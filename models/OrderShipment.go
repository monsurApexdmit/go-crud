package models

import (
	"time"
)

type ShipmentTrackingHistory struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ShipmentID  uint      `json:"shipmentId" gorm:"column:shipment_id;not null"`
	Status      string    `json:"status" gorm:"type:varchar(100);not null"`
	Location    string    `json:"location,omitempty" gorm:"type:varchar(255)"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	EventTime   time.Time `json:"eventTime" gorm:"column:event_time;not null"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
}

func (ShipmentTrackingHistory) TableName() string { return "shipment_tracking_history" }

type OrderShipment struct {
	ID                uint                      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID         uint                      `json:"companyId" gorm:"column:company_id;not null"`
	SellID            uint                      `json:"sellId" gorm:"column:sell_id;not null"`
	Sell              *Sell                     `json:"sell,omitempty" gorm:"foreignKey:SellID"`
	TrackingNumber    string                    `json:"trackingNumber" gorm:"type:varchar(100);column:tracking_number;not null"`
	Carrier           string                    `json:"carrier" gorm:"type:varchar(100);not null"`
	ShippingMethod    string                    `json:"shippingMethod,omitempty" gorm:"type:varchar(100);column:shipping_method"`
	Status            string                    `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	ShippedAt         *time.Time                `json:"shippedAt,omitempty" gorm:"column:shipped_at"`
	EstimatedDelivery *time.Time                `json:"estimatedDelivery,omitempty" gorm:"column:estimated_delivery"`
	DeliveredAt       *time.Time                `json:"deliveredAt,omitempty" gorm:"column:delivered_at"`
	ShippingCost      float64                   `json:"shippingCost" gorm:"column:shipping_cost;default:0"`
	Weight            float64                   `json:"weight,omitempty" gorm:"default:0"`
	Dimensions        string                    `json:"dimensions,omitempty" gorm:"type:varchar(100)"`
	Notes             string                    `json:"notes,omitempty" gorm:"type:text"`
	TrackingHistory   []ShipmentTrackingHistory `json:"trackingHistory,omitempty" gorm:"foreignKey:ShipmentID"`
	CreatedAt         time.Time                 `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time                 `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (OrderShipment) TableName() string { return "order_shipments" }
