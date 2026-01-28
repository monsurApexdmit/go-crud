package models

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CampaignName string         `json:"campaign_name" form:"campaign_name"`
	Code         string         `json:"code" gorm:"unique" form:"code"`
	Discount     float64        `json:"discount" form:"discount"`
	Type         string         `json:"type" form:"type"` // percentage, fixed
	StartDate    time.Time      `json:"start_date" form:"start_date" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate      time.Time      `json:"end_date" form:"end_date" time_format:"2006-01-02T15:04:05Z07:00"`
	Status       bool           `json:"status" gorm:"default:false" form:"status"`
	Image        string         `json:"image" form:"image"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
