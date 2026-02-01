package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductVariant struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProductID  uint           `json:"product_id"`
	Name       string         `json:"name"` // e.g., "Small / Red"
	Attributes datatypes.JSON `json:"attributes"` // {"Size": "Small", "Color": "Red"}
	Price      float64        `gorm:"not null;default:0" json:"price"`
	SalePrice  float64        `gorm:"not null;default:0" json:"sale_price"`
	Stock      int            `gorm:"not null;default:0" json:"stock"`
	SKU        string         `gorm:"uniqueIndex" json:"sku"`
	Barcode    string         `gorm:"uniqueIndex" json:"barcode"`
	
	Inventory  []VariantInventory `gorm:"constraint:OnDelete:CASCADE" json:"inventory,omitempty"`
	
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
