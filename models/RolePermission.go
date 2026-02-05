package models

import "time"

type RolePermission struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       uint       `json:"roleId" gorm:"column:role_id;not null"`
	PermissionID uint       `json:"permissionId" gorm:"column:permission_id;not null"`
	Permission   Permission `json:"permission" gorm:"foreignKey:PermissionID"`
	Read         bool       `json:"read" gorm:"column:read;default:false"`
	Write        bool       `json:"write" gorm:"column:write;default:false"`
	Delete       bool       `json:"delete" gorm:"column:delete;default:false"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (RolePermission) TableName() string { return "role_permissions" }
