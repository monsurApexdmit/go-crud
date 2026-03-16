package models

import (
	"time"
	"gorm.io/gorm"
)

type Attribute struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID   uint           `json:"companyId" gorm:"column:company_id;not null;index"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;uniqueIndex:idx_attribute_company_name,composite:company_id"`
	DisplayName string         `json:"display_name" gorm:"type:varchar(150);not null"`
	OptionType  string         `json:"option_type" gorm:"type:varchar(50);not null;default:'text'"` // text, dropdown, radio, checkbox, color, size
	Values      string         `json:"values" gorm:"type:text"` // JSON array of possible values for dropdown/radio/checkbox
	Description string         `json:"description" gorm:"type:text"`
	IsRequired  bool           `json:"is_required" gorm:"default:false"`
	Status      bool           `json:"status" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
