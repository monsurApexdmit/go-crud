package models
import (
    "time"
	"gorm.io/gorm"
)

type User struct {
    ID          uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string          `json:"username"`
	Email       string          `json:"email"`
    Password    string          `json:"-" gorm:"column:password"`

	RoleID      uint            `json:"role_id"`         
	Role        Role            `json:"role" gorm:"foreignKey:RoleID"`
    
    Address     string          `json:"address,omitempty"`
    CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`

}