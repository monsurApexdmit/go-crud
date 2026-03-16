package models

import "time"

type Author struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name" gorm:"type:varchar(255)"`
	Email string `json:"email" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}