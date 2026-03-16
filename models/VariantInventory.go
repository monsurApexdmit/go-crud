package models

import (
	"time"
)

type VariantInventory struct {
	ID         uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	VariantID  uint     `json:"variant_id"`
	LocationID uint     `json:"location_id"`
	Location   Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	Quantity   int      `gorm:"not null;default:0" json:"quantity"`
	
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (VariantInventory) TableName() string {
	return "variant_inventory"
}
