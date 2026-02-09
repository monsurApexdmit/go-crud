package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerReturnItem struct {
	ID          uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	ReturnID    uint            `json:"returnId" gorm:"column:return_id;not null"`
	ProductID   *uint           `json:"productId,omitempty" gorm:"column:product_id"`
	Product     *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductName string          `json:"productName" gorm:"column:product_name;not null"`
	VariantID   *uint           `json:"variantId,omitempty" gorm:"column:variant_id"`
	Variant     *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	VariantName string          `json:"variantName,omitempty" gorm:"column:variant_name"`
	Quantity    int             `json:"quantity" gorm:"not null;default:1"`
	Price       float64         `json:"price" gorm:"not null;default:0"`
	Reason      string          `json:"reason" gorm:"not null"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (CustomerReturnItem) TableName() string { return "customer_return_items" }

type CustomerReturn struct {
	ID            uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	ReturnNumber  string               `json:"returnNumber" gorm:"column:return_number;uniqueIndex;not null"`
	CustomerID    uint                 `json:"customerId" gorm:"column:customer_id;not null"`
	Customer      *Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CustomerName  string               `json:"customerName" gorm:"column:customer_name;not null"`
	OrderID       *uint                `json:"orderId,omitempty" gorm:"column:order_id"`
	Order         *Sell                `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	OrderNumber   string               `json:"orderNumber,omitempty" gorm:"column:order_number"`
	Items         []CustomerReturnItem `json:"items,omitempty" gorm:"foreignKey:ReturnID"`
	TotalAmount   float64              `json:"totalAmount" gorm:"column:total_amount;not null;default:0"`
	Status        string               `json:"status" gorm:"not null;default:'pending'"`
	RequestDate   time.Time            `json:"requestDate" gorm:"column:request_date;not null"`
	ProcessedDate *time.Time           `json:"processedDate,omitempty" gorm:"column:processed_date"`
	RefundMethod  string               `json:"refundMethod" gorm:"column:refund_method;not null"`
	Notes         string               `json:"notes,omitempty" gorm:"type:text"`
	ProcessedBy   string               `json:"processedBy,omitempty" gorm:"column:processed_by"`
	CreatedAt     time.Time            `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time            `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt       `json:"-" gorm:"index"`
}

func (CustomerReturn) TableName() string { return "customer_returns" }
