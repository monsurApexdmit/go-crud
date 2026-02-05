package models

import (
	"time"

	"gorm.io/gorm"
)

type StaffRole struct {
	ID              uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string           `json:"name" gorm:"uniqueIndex;not null"`
	RolePermissions []RolePermission `json:"-" gorm:"foreignKey:RoleID"`
	CreatedAt       time.Time        `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time        `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt   `json:"-" gorm:"index"`
}

func (StaffRole) TableName() string { return "staff_roles" }
