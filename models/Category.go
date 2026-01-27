package models

import (
	"time"
	"gorm.io/gorm"
)


type Category struct{
	ID          	uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryName    string          `json:"category_name"`
	ParentID     	*uint          	`json:"parent_id"`
	Parent       	*Category      	`gorm:"foreignKey:ParentID"`
	Children     	[]Category     	`gorm:"foreignKey:ParentID"`
	Status    		bool           	`json:"status" gorm:"default:false"`
    CreatedAt   	time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   	time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   	gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}