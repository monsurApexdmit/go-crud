package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID     uint      `json:"companyId" gorm:"column:company_id;not null;index"`
	Name          string    `gorm:"not null;type:varchar(255)" json:"name"`
	Description   string    `json:"description" gorm:"type:text"`

	// Foreign Keys
	CategoryID    *uint     `json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	VendorID      *uint     `json:"vendor_id"`
	Vendor        User      `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`

	LocationID    *uint     `json:"location_id"`
	Location      Location  `gorm:"foreignKey:LocationID" json:"location,omitempty"`

	// Product Details
	Price         float64   `gorm:"not null;default:0" json:"price"`
	SalePrice     float64   `gorm:"not null;default:0" json:"sale_price"`
	Stock         int       `gorm:"not null;default:0" json:"stock"`
	SKU           string    `gorm:"type:varchar(100);uniqueIndex:idx_product_company_sku,composite:company_id" json:"sku"`
	Barcode       string    `gorm:"type:varchar(100);uniqueIndex:idx_product_company_barcode,composite:company_id" json:"barcode"`
	Published     bool      `gorm:"not null;default:false" json:"published"`
	ReceiptNumber string    `json:"receipt_number"`
	Image         string    `json:"image"`

	// Relationships
	Attributes []Attribute        `gorm:"many2many:product_attributes" json:"attributes,omitempty"`
	Variants   []ProductVariant   `gorm:"constraint:OnDelete:CASCADE" json:"variants,omitempty"`
	Images     []ProductImage     `gorm:"constraint:OnDelete:CASCADE" json:"images,omitempty"`

    CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`

}
