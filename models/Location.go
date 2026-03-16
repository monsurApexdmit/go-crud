package models

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CompanyID     uint           `json:"companyId" gorm:"column:company_id;not null;index"`
	Name          string         `json:"name" form:"name"`
	Address       string         `json:"address" form:"address"`
	ContactPerson string         `json:"contact_person" form:"contact_person"`
	IsDefault     bool           `json:"is_default" form:"is_default" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
