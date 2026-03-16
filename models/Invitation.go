package models

import "time"

// Invitation represents a team member invitation
type Invitation struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	CompanyID       uint      `json:"companyId" gorm:"index"`
	Email           string    `json:"email" gorm:"type:varchar(255);index"`
	FullName        string    `json:"fullName"`
	Role            string    `json:"role"` // admin, manager, staff
	Status          string    `json:"status" gorm:"default:'pending'"` // pending, accepted, expired
	InvitationToken string    `json:"invitationToken" gorm:"type:varchar(255);uniqueIndex"`
	ExpiresAt       time.Time `json:"expiresAt"`
	AcceptedAt      *time.Time `json:"acceptedAt,omitempty"`
	InvitedAt       time.Time `json:"invitedAt"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	// Relations
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

// TableName specifies the table name for the Invitation model
func (Invitation) TableName() string {
	return "invitations"
}
