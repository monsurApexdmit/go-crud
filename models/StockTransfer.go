package models

import "time"

type StockTransfer struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID      uint      `json:"companyId" gorm:"column:company_id;not null"`
	ProductID      uint      `json:"productId" gorm:"column:product_id;not null"`
	Product        Product   `json:"product" gorm:"foreignKey:ProductID"`
	VariantID      *uint     `json:"variantId" gorm:"column:variant_id"`
	Variant        *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	FromLocationID uint      `json:"fromLocationId" gorm:"column:from_location_id;not null"`
	FromLocation   Location  `json:"fromLocation" gorm:"foreignKey:FromLocationID"`
	ToLocationID   uint      `json:"toLocationId" gorm:"column:to_location_id;not null"`
	ToLocation     Location  `json:"toLocation" gorm:"foreignKey:ToLocationID"`
	Quantity       int       `json:"quantity" gorm:"not null"`
	Status         string    `json:"status" gorm:"type:enum('Pending','Completed','Cancelled');default:Pending"`
	Notes          string    `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (StockTransfer) TableName() string { return "stock_transfers" }
