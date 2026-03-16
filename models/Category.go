package models

import (
	"time"
	"gorm.io/gorm"
)


type Category struct{
	ID          	uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID		uint 			`json:"companyId" gorm:"column:company_id;not null;index"`
	CategoryName    string          `json:"category_name" gorm:"type:varchar(100);not null"`
	ParentID     	*uint          	`json:"parent_id" gorm:"index"`
	Parent       	*Category      	`json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children     	[]Category     	`json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Status    		bool           	`json:"status" gorm:"default:true"`
    CreatedAt   	time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   	time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   	gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}