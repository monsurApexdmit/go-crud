package models

import (
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type GeneralSettings struct {
	StoreName       string `json:"storeName"`
	StoreEmail      string `json:"storeEmail"`
	StorePhone      string `json:"storePhone"`
	StoreAddress    string `json:"storeAddress"`
	StoreDescription string `json:"storeDescription"`
}

type TaxSettings struct {
	DefaultTaxRate      float64 `json:"defaultTaxRate"`
	TaxInclusivePrice   bool    `json:"taxInclusivePrice"`
	EnableGSTTracking   bool    `json:"enableGSTTracking"`
	GSTNumber           string  `json:"gstNumber"`
	EnableTaxExemption  bool    `json:"enableTaxExemption"`
	DefaultShippingTax  float64 `json:"defaultShippingTax"`
}

type ShippingMethod struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Cost            float64 `json:"cost"`
	EstimatedDays   int     `json:"estimatedDays"`
	IsActive        bool    `json:"isActive"`
}

type ShippingSettings struct {
	EnableShipping          bool              `json:"enableShipping"`
	DefaultShippingCost     float64           `json:"defaultShippingCost"`
	FreeShippingThreshold   float64           `json:"freeShippingThreshold"`
	ShippingMethods         []ShippingMethod  `json:"shippingMethods"`
}

type PaymentSettings struct {
	EnableCash          bool    `json:"enableCash"`
	EnableCard          bool    `json:"enableCard"`
	EnableOnline        bool    `json:"enableOnline"`
	CardProcessingFee   float64 `json:"cardProcessingFee"`
	StripeKey           string  `json:"stripeKey"`
	RazorpayKey         string  `json:"razorpayKey"`
}

type SocialLinks struct {
	Facebook string `json:"facebook"`
	Instagram string `json:"instagram"`
	Twitter  string `json:"twitter"`
}

type BusinessSettings struct {
	BusinessName       string      `json:"businessName"`
	BusinessType       string      `json:"businessType"`
	RegistrationNumber string      `json:"registrationNumber"`
	GSTNumber          string      `json:"gstNumber"`
	Website            string      `json:"website"`
	SocialLinks        SocialLinks `json:"socialLinks"`
}

type StoreHours struct {
	Open   string `json:"open"`
	Close  string `json:"close"`
	IsOpen bool   `json:"isOpen"`
}

type RegionalSettings struct {
	Language string `json:"language"`
	Currency string `json:"currency"`
	Timezone string `json:"timezone"`
}

type NotificationSettings struct {
	EmailNotifications    bool `json:"emailNotifications"`
	OrderNotifications    bool `json:"orderNotifications"`
	MarketingEmails       bool `json:"marketingEmails"`
}

type StoreHoursData struct {
	Monday    StoreHours `json:"monday"`
	Tuesday   StoreHours `json:"tuesday"`
	Wednesday StoreHours `json:"wednesday"`
	Thursday  StoreHours `json:"thursday"`
	Friday    StoreHours `json:"friday"`
	Saturday  StoreHours `json:"saturday"`
	Sunday    StoreHours `json:"sunday"`
}

// Settings main model
type Settings struct {
	ID                   uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	CompanyID            uint             `json:"companyId" gorm:"column:company_id;not null;index"`
	GeneralSettings      datatypes.JSON   `json:"generalSettings" gorm:"type:json;column:general_settings"`
	TaxSettings          datatypes.JSON   `json:"taxSettings" gorm:"type:json;column:tax_settings"`
	ShippingSettings     datatypes.JSON   `json:"shippingSettings" gorm:"type:json;column:shipping_settings"`
	PaymentSettings      datatypes.JSON   `json:"paymentSettings" gorm:"type:json;column:payment_settings"`
	BusinessSettings     datatypes.JSON   `json:"businessSettings" gorm:"type:json;column:business_settings"`
	RegionalSettings     datatypes.JSON   `json:"regionalSettings" gorm:"type:json;column:regional_settings"`
	NotificationSettings datatypes.JSON   `json:"notificationSettings" gorm:"type:json;column:notification_settings"`
	StoreHours           datatypes.JSON   `json:"storeHours" gorm:"type:json;column:store_hours"`
	LogoURL              string           `json:"logoUrl,omitempty" gorm:"column:logo_url"`
	BannerURL            string           `json:"bannerUrl,omitempty" gorm:"column:banner_url"`
	CreatedAt            time.Time        `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time        `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt            gorm.DeletedAt   `json:"-" gorm:"index"`
}

func (Settings) TableName() string { return "settings" }

// Helper methods to marshal/unmarshal JSON settings
func (s *Settings) GetGeneralSettings() GeneralSettings {
	var gs GeneralSettings
	json.Unmarshal(s.GeneralSettings, &gs)
	return gs
}

func (s *Settings) SetGeneralSettings(gs GeneralSettings) error {
	data, err := json.Marshal(gs)
	if err != nil {
		return err
	}
	s.GeneralSettings = data
	return nil
}

func (s *Settings) GetTaxSettings() TaxSettings {
	var ts TaxSettings
	json.Unmarshal(s.TaxSettings, &ts)
	return ts
}

func (s *Settings) SetTaxSettings(ts TaxSettings) error {
	data, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	s.TaxSettings = data
	return nil
}

func (s *Settings) GetShippingSettings() ShippingSettings {
	var ss ShippingSettings
	json.Unmarshal(s.ShippingSettings, &ss)
	return ss
}

func (s *Settings) SetShippingSettings(ss ShippingSettings) error {
	data, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	s.ShippingSettings = data
	return nil
}

func (s *Settings) GetPaymentSettings() PaymentSettings {
	var ps PaymentSettings
	json.Unmarshal(s.PaymentSettings, &ps)
	return ps
}

func (s *Settings) SetPaymentSettings(ps PaymentSettings) error {
	data, err := json.Marshal(ps)
	if err != nil {
		return err
	}
	s.PaymentSettings = data
	return nil
}

func (s *Settings) GetBusinessSettings() BusinessSettings {
	var bs BusinessSettings
	json.Unmarshal(s.BusinessSettings, &bs)
	return bs
}

func (s *Settings) SetBusinessSettings(bs BusinessSettings) error {
	data, err := json.Marshal(bs)
	if err != nil {
		return err
	}
	s.BusinessSettings = data
	return nil
}

func (s *Settings) GetRegionalSettings() RegionalSettings {
	var rs RegionalSettings
	json.Unmarshal(s.RegionalSettings, &rs)
	return rs
}

func (s *Settings) SetRegionalSettings(rs RegionalSettings) error {
	data, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	s.RegionalSettings = data
	return nil
}

func (s *Settings) GetNotificationSettings() NotificationSettings {
	var ns NotificationSettings
	json.Unmarshal(s.NotificationSettings, &ns)
	return ns
}

func (s *Settings) SetNotificationSettings(ns NotificationSettings) error {
	data, err := json.Marshal(ns)
	if err != nil {
		return err
	}
	s.NotificationSettings = data
	return nil
}

func (s *Settings) GetStoreHours() StoreHoursData {
	var sh StoreHoursData
	json.Unmarshal(s.StoreHours, &sh)
	return sh
}

func (s *Settings) SetStoreHours(sh StoreHoursData) error {
	data, err := json.Marshal(sh)
	if err != nil {
		return err
	}
	s.StoreHours = data
	return nil
}
