package models

import "time"

type ProductImage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Path      string    `gorm:"not null" json:"path"`
	Position  int       `gorm:"not null;default:0" json:"position"`
	IsPrimary bool      `gorm:"not null;default:false" json:"is_primary"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
