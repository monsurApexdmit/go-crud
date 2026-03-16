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
	ProductName string          `json:"productName" gorm:"type:varchar(255);column:product_name;not null"`
	VariantID   *uint           `json:"variantId,omitempty" gorm:"column:variant_id"`
	Variant     *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	VariantName string          `json:"variantName,omitempty" gorm:"type:varchar(255);column:variant_name"`
	Quantity    int             `json:"quantity" gorm:"not null;default:1"`
	Price       float64         `json:"price" gorm:"not null;default:0"`
	Reason      string          `json:"reason" gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (CustomerReturnItem) TableName() string { return "customer_return_items" }

type CustomerReturn struct {
	ID            uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID     uint                 `json:"companyId" gorm:"column:company_id;not null;index"`
	ReturnNumber  string               `json:"returnNumber" gorm:"type:varchar(100);column:return_number;uniqueIndex;not null"`
	CustomerID    *uint                `json:"customerId,omitempty" gorm:"column:customer_id"`
	Customer      *Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	CustomerName  string               `json:"customerName,omitempty" gorm:"type:varchar(255);column:customer_name"`
	OrderID       *uint                `json:"orderId,omitempty" gorm:"column:order_id"`
	Order         *Sell                `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	OrderNumber   string               `json:"orderNumber,omitempty" gorm:"type:varchar(100);column:order_number"`
	Items         []CustomerReturnItem `json:"items,omitempty" gorm:"foreignKey:ReturnID"`
	TotalAmount   float64              `json:"totalAmount" gorm:"column:total_amount;not null;default:0"`
	Status        string               `json:"status" gorm:"type:varchar(50);not null;default:'pending'"`
	RequestDate   time.Time            `json:"requestDate" gorm:"column:request_date;not null"`
	ProcessedDate *time.Time           `json:"processedDate,omitempty" gorm:"column:processed_date"`
	RefundMethod  string               `json:"refundMethod" gorm:"type:varchar(100);column:refund_method;not null"`
	Notes         string               `json:"notes,omitempty" gorm:"type:text"`
	ProcessedBy   string               `json:"processedBy,omitempty" gorm:"type:varchar(255);column:processed_by"`
	CreatedAt     time.Time            `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time            `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt       `json:"-" gorm:"index"`
}

func (CustomerReturn) TableName() string { return "customer_returns" }
