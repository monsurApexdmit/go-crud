package models

import "time"

// PasswordReset represents a password reset request
type PasswordReset struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          uint      `json:"userId" gorm:"index"`
	Email           string    `json:"email" gorm:"type:varchar(255);index"`
	ResetToken      string    `json:"resetToken" gorm:"type:varchar(255);uniqueIndex"`
	ExpiresAt       time.Time `json:"expiresAt"`
	Status          string    `json:"status" gorm:"default:'pending'"` // pending, used, expired
	UsedAt          *time.Time `json:"usedAt,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relations
	User *SaasUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for the PasswordReset model
func (PasswordReset) TableName() string {
	return "password_resets"
}
