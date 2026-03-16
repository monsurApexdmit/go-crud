package models

import (
	"time"

	"gorm.io/gorm"
)

type VendorReturnItem struct {
	ID          uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	ReturnID    uint            `json:"returnId" gorm:"column:return_id;not null"`
	ProductID   *uint           `json:"productId,omitempty" gorm:"column:product_id"`
	Product     *Product        `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	ProductName string          `json:"productName" gorm:"column:product_name;not null"`
	VariantID   *uint           `json:"variantId,omitempty" gorm:"column:variant_id"`
	Variant     *ProductVariant `json:"variant,omitempty" gorm:"foreignKey:VariantID"`
	VariantName string          `json:"variantName,omitempty" gorm:"column:variant_name"`
	Quantity    int             `json:"quantity" gorm:"not null;default:1"`
	UnitPrice   float64         `json:"unitPrice" gorm:"column:unit_price;not null;default:0"`
	TotalPrice  float64         `json:"totalPrice" gorm:"column:total_price;not null;default:0"`
	Reason      string          `json:"reason" gorm:"not null"`
	CreatedAt   time.Time       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (VendorReturnItem) TableName() string { return "vendor_return_items" }

type VendorReturn struct {
	ID            uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID     uint               `json:"companyId" gorm:"column:company_id;not null;index"`
	ReturnNumber  string             `json:"returnNumber" gorm:"type:varchar(100);column:return_number;uniqueIndex;not null"`
	VendorID      uint               `json:"vendorId" gorm:"column:vendor_id;not null"`
	Vendor        *Vendor            `json:"vendor,omitempty" gorm:"foreignKey:VendorID"`
	VendorName    string             `json:"vendorName" gorm:"type:varchar(255);column:vendor_name;not null"`
	Items         []VendorReturnItem `json:"items,omitempty" gorm:"foreignKey:ReturnID"`
	TotalAmount   float64            `json:"totalAmount" gorm:"column:total_amount;not null;default:0"`
	Status        string             `json:"status" gorm:"not null;default:'pending'"`
	ReturnDate    time.Time          `json:"returnDate" gorm:"column:return_date;not null"`
	CompletedDate *time.Time         `json:"completedDate,omitempty" gorm:"column:completed_date"`
	CreditType    string             `json:"creditType" gorm:"column:credit_type;not null"`
	Notes         string             `json:"notes,omitempty" gorm:"type:text"`
	CreatedBy     string             `json:"createdBy,omitempty" gorm:"column:created_by"`
	CreatedAt     time.Time          `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt     time.Time          `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt     `json:"-" gorm:"index"`
}

func (VendorReturn) TableName() string { return "vendor_returns" }
