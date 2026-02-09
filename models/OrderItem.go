package models

import "time"

type OrderItem struct {
	ID          uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	SellID      uint            `json:"sellId" gorm:"column:sell_id;not null"`
	ProductID   *uint           `json:"productId,omitempty" gorm:"column:product_id"`
	Product     *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	VariantID   *uint           `json:"variantId,omitempty" gorm:"column:variant_id"`
	Variant     *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	ProductName string          `json:"productName" gorm:"column:product_name;not null"`
	VariantName string          `json:"variantName,omitempty" gorm:"column:variant_name"`
	Quantity    int             `json:"quantity" gorm:"not null;default:1"`
	UnitPrice   float64         `json:"unitPrice" gorm:"column:unit_price;not null;default:0"`
	TotalPrice  float64         `json:"totalPrice" gorm:"column:total_price;not null;default:0"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (OrderItem) TableName() string { return "order_items" }
