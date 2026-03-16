package models
import (
    "time"
	"gorm.io/gorm"
)

type User struct {
    ID          uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	Username    string          `json:"username" gorm:"type:varchar(100)"`
	Email       string          `json:"email" gorm:"type:varchar(255);unique"`
    Password    string          `json:"-" gorm:"column:password;type:varchar(255)"`

	RoleID      uint            `json:"role_id"`
	Role        Role            `json:"role" gorm:"foreignKey:RoleID"`

    Address     string          `json:"address,omitempty" gorm:"type:text"`
    CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt   gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`

}