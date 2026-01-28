package models

import (
	"time"
)

type Attribute struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	DisplayName string    `json:"display_name" gorm:"not null"`
	OptionType  string    `json:"option_type"` // e.g., dropdown, radio
	Values      string    `json:"values"`      // comma-separated values
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
