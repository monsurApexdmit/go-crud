package dto

import "time"

// ============= AUTH DTOs =============

type SignupRequest struct {
	CompanyName    string `json:"companyName" binding:"required"`
	OwnerFullName  string `json:"ownerFullName" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Phone          string `json:"phone" binding:"required"`
	Password       string `json:"password" binding:"required,min=8"`
	BusinessType   string `json:"businessType"`
	Website        string `json:"website"`
	Country        string `json:"country"`
}

type SignupResponse struct {
	Message string                 `json:"message"`
	Data    SignupData             `json:"data"`
}

type SignupData struct {
	UserID               uint                 `json:"userId"`
	CompanyID            uint                 `json:"companyId"`
	Email                string               `json:"email"`
	UserEmail            string               `json:"userEmail"`
	CompanyName          string               `json:"companyName"`
	UserRole             string               `json:"userRole"`
	CompanyStatus        string               `json:"companyStatus"`
	Token                string               `json:"token"`
	LicenseKey           string               `json:"licenseKey"`
	LicenseType          string               `json:"licenseType"`
	TrialStartDate       string               `json:"trialStartDate"`
	TrialEndDate         string               `json:"trialEndDate"`
	TrialDaysRemaining   int                  `json:"trialDaysRemaining"`
	Company              CompanyDTO           `json:"company"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Message string               `json:"message"`
	Data    LoginData            `json:"data"`
}

type LoginData struct {
	UserID          uint                 `json:"userId"`
	UserEmail       string               `json:"userEmail"`
	CompanyID       uint                 `json:"companyId"`
	CompanyName     string               `json:"companyName"`
	CompanyStatus   string               `json:"companyStatus"`
	UserRole        string               `json:"userRole"`
	Token           string               `json:"token"`
	LicenseKey      string               `json:"licenseKey"`
	LicenseType     string               `json:"licenseType"`
	Email           string               `json:"email"` // Keep for backward compatibility
	Company         CompanyDTO           `json:"company"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

// ============= COMPANY DTOs =============

type CompanyDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Industry    string    `json:"industry,omitempty"`
	Website     string    `json:"website,omitempty"`
	Phone       string    `json:"phone,omitempty"`
	Address     string    `json:"address,omitempty"`
	City        string    `json:"city,omitempty"`
	State       string    `json:"state,omitempty"`
	ZipCode     string    `json:"zipCode,omitempty"`
	Country     string    `json:"country,omitempty"`
	Logo        string    `json:"logo,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateCompanyProfileRequest struct {
	CompanyProfile CompanyProfileUpdateDTO `json:"companyProfile"`
}

type CompanyProfileUpdateDTO struct {
	Name        string `json:"name,omitempty"`
	Industry    string `json:"industry,omitempty"`
	Website     string `json:"website,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	State       string `json:"state,omitempty"`
	ZipCode     string `json:"zipCode,omitempty"`
	Country     string `json:"country,omitempty"`
	Description string `json:"description,omitempty"`
}

type CompanyStatusDTO struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Status               string    `json:"status"`
	UserCount            int       `json:"userCount"`
	MaxUsers             int       `json:"maxUsers"`
	CreatedAt            time.Time `json:"createdAt"`
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ============= COMPANY SETTINGS DTOs =============

type CompanySettingsDTO struct {
	ID          uint      `json:"id"`
	CompanyID   uint      `json:"companyId"`
	CompanyName string    `json:"companyName"`
	TaxID       string    `json:"taxId,omitempty"`
	TaxIDType   string    `json:"taxIdType,omitempty"`
	TaxRate     float64   `json:"taxRate"`
	Currency    string    `json:"currency"`
	Timezone    string    `json:"timezone"`
	Language    string    `json:"language"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UpdateCompanySettingsRequest struct {
	CompanyName string  `json:"companyName,omitempty"`
	TaxID       string  `json:"taxId,omitempty"`
	TaxIDType   string  `json:"taxIdType,omitempty"`
	TaxRate     float64 `json:"taxRate,omitempty"`
	Currency    string  `json:"currency,omitempty"`
	Timezone    string  `json:"timezone,omitempty"`
	Language    string  `json:"language,omitempty"`
}

// ============= BILLING CONTACT DTOs =============

type BillingContactDTO struct {
	ID          uint      `json:"id"`
	CompanyID   uint      `json:"companyId"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	ZipCode     string    `json:"zipCode"`
	Country     string    `json:"country"`
	TaxID       string    `json:"taxId,omitempty"`
	TaxIDType   string    `json:"taxIdType,omitempty"`
	IsDefault   bool      `json:"isDefault"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type BillingContactRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Address   string `json:"address" binding:"required"`
	City      string `json:"city" binding:"required"`
	State     string `json:"state"`
	ZipCode   string `json:"zipCode"`
	Country   string `json:"country"`
	TaxID     string `json:"taxId,omitempty"`
	TaxIDType string `json:"taxIdType,omitempty"`
}

// ============= SUBSCRIPTION DTOs =============

type SubscriptionDTO struct {
	ID                uint               `json:"id"`
	PlanID            uint               `json:"planId"`
	PlanName          string             `json:"planName"`
	Price             int64              `json:"price"`
	BillingPeriod     string             `json:"billingPeriod"`
	Status            string             `json:"status"`
	CurrentPeriodStart time.Time         `json:"currentPeriodStart"`
	CurrentPeriodEnd  time.Time          `json:"currentPeriodEnd"`
	NextBillingDate   time.Time          `json:"nextBillingDate"`
	AutoRenew         bool               `json:"autoRenew"`
	CreatedAt         time.Time          `json:"createdAt"`
}

type SubscriptionPlanDTO struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         int64     `json:"price"`
	BillingPeriod string    `json:"billingPeriod"`
	MaxUsers      int       `json:"maxUsers"`
	MaxProducts   int       `json:"maxProducts"`
	MaxBranches   int       `json:"maxBranches"`
	Features      string    `json:"features"`
	IsActive      bool      `json:"isActive"`
	IsFeatured    bool      `json:"isFeatured"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PaymentRecordDTO struct {
	ID            uint       `json:"id"`
	PaymentDate   *time.Time `json:"paymentDate,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	InvoiceNumber string     `json:"invoiceNumber"`
	Amount        int64      `json:"amount"`
	Status        string     `json:"status"`
	InvoiceUrl    string     `json:"invoiceUrl,omitempty"`
}

// ============= TEAM USER DTOs =============

type TeamUserDTO struct {
	ID        uint       `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"fullName"`
	Role      string     `json:"role"`
	Status    string     `json:"status"`
	JoinedDate time.Time `json:"joinedDate"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	Avatar    string     `json:"avatar,omitempty"`
}

type TeamUsersResponse struct {
	Message   string          `json:"message"`
	Data      TeamUsersData   `json:"data"`
}

type TeamUsersData struct {
	Users       []TeamUserDTO `json:"users"`
	TotalUsers  int          `json:"totalUsers"`
	MaxUsers    int          `json:"maxUsers"`
	CanAddMore  bool         `json:"canAddMore"`
}

type InviteUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"fullName" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin manager staff"`
}

type InviteUserResponse struct {
	Message string        `json:"message"`
	Data    InviteData    `json:"data"`
}

type InviteData struct {
	UserID            uint      `json:"userId"`
	Email             string    `json:"email"`
	Status            string    `json:"status"`
	InvitationToken   string    `json:"invitationToken"`
	ExpiresAt         time.Time `json:"expiresAt"`
	InvitedAt         time.Time `json:"invitedAt"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin manager staff"`
}

type UpdateUserRoleResponse struct {
	Message string           `json:"message"`
	Data    UpdateRoleData   `json:"data"`
}

type UpdateRoleData struct {
	UserID   uint      `json:"userId"`
	Role     string    `json:"role"`
	UpdatedAt time.Time `json:"updatedAt"`
}
