package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ProductVariant struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	ProductID  uint           `gorm:"uniqueIndex:idx_product_sku" json:"product_id"`
	Name       string         `json:"name" gorm:"type:varchar(255)"` // e.g., "Small / Red"
	Attributes datatypes.JSON `json:"attributes"` // {"Size": "Small", "Color": "Red"}
	Price      float64        `gorm:"not null;default:0" json:"price"`
	SalePrice  float64        `gorm:"not null;default:0" json:"sale_price"`
	Stock      int            `gorm:"not null;default:0" json:"stock"`
	SKU        string         `gorm:"type:varchar(100);uniqueIndex:idx_product_sku" json:"sku"`
	Barcode    string         `gorm:"type:varchar(100);index" json:"barcode"`
	
	Inventory  []VariantInventory `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE" json:"inventory,omitempty"`
	
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
