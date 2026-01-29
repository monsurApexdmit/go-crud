package models

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `json:"name" form:"name"`
	Address       string         `json:"address" form:"address"`
	ContactPerson string         `json:"contact_person" form:"contact_person"`
	IsDefault     bool           `json:"is_default" form:"is_default" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
