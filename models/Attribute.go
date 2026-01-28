package models

import (
	"time"
	"gorm.io/gorm"
)

type Attribute struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	DisplayName string    `json:"display_name" gorm:"not null"`
	OptionType  string    `json:"option_type"` // e.g., dropdown, radio
	Values      string    `json:"values"`      // comma-separated values
    CreatedAt   	time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   	time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   	gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}
