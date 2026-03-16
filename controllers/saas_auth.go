package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go-crud/database"
	"go-crud/dto"
	"go-crud/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// SaasSignup registers a new company with owner
func SaasSignup(c *gin.Context) {
	var req dto.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.SaasUser
	if database.DB.Where("email = ?", req.Email).First(&existingUser).RowsAffected > 0 {
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Message: "Email already registered",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to process password",
		})
		return
	}

	// Create company
	company := models.Company{
		Name:   req.CompanyName,
		Status: "trial",
	}

	if err := database.DB.Create(&company).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create company",
		})
		return
	}

	// Create owner user
	user := models.SaasUser{
		CompanyID:  company.ID,
		Email:      req.Email,
		FullName:   req.OwnerFullName,
		Password:   string(hashedPassword),
		Role:       "owner",
		Status:     "active",
		JoinedDate: time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create user",
		})
		return
	}

	// Create subscription (trial)
	trialEndDate := time.Now().AddDate(0, 0, 10) // 10-day trial
	subscription := models.Subscription{
		CompanyID:          company.ID,
		PlanID:             1, // Trial plan
		Status:             "active",
		CurrentPeriodStart: time.Now(),
		CurrentPeriodEnd:   trialEndDate,
		NextBillingDate:    trialEndDate,
		AutoRenew:          true,
	}

	if err := database.DB.Create(&subscription).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create subscription",
		})
		return
	}

	// Create company settings
	settings := models.CompanySettings{
		CompanyID:   company.ID,
		CompanyName: company.Name,
		Currency:    "USD",
		Timezone:    "UTC",
		Language:    "en",
		TaxRate:     0,
	}

	if err := database.DB.Create(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create settings",
		})
		return
	}

	// Generate JWT token
	token := generateJWT(user.ID, company.ID, user.Email)

	trialDaysRemaining := 10
	trialStartDateStr := time.Now().Format(time.RFC3339)
	trialEndDateStr := time.Now().AddDate(0, 0, 10).Format(time.RFC3339)

	c.JSON(http.StatusCreated, dto.SignupResponse{
		Message: "Company created successfully. Trial activated for 10 days.",
		Data: dto.SignupData{
			UserID:             user.ID,
			CompanyID:          company.ID,
			Email:              user.Email,
			UserEmail:          user.Email,
			CompanyName:        company.Name,
			UserRole:           user.Role,
			CompanyStatus:      company.Status,
			Token:              token,
			LicenseKey:         "trial-" + fmt.Sprint(company.ID),
			LicenseType:        "trial",
			TrialStartDate:     trialStartDateStr,
			TrialEndDate:       trialEndDateStr,
			TrialDaysRemaining: trialDaysRemaining,
			Company: dto.CompanyDTO{
				ID:     company.ID,
				Name:   company.Name,
				Status: company.Status,
			},
		},
	})
}

// SaasLogin authenticates user and returns JWT token
func SaasLogin(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Find user
	var user models.SaasUser
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid email or password",
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid email or password",
		})
		return
	}

	// Get company
	var company models.Company
	if err := database.DB.First(&company, user.CompanyID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to fetch company",
		})
		return
	}

	// Update last login
	database.DB.Model(&user).Update("last_login", time.Now())

	// Generate JWT token
	token := generateJWT(user.ID, company.ID, user.Email)

	c.JSON(http.StatusOK, dto.LoginResponse{
		Message: "Login successful",
		Data: dto.LoginData{
			UserID:          user.ID,
			UserEmail:       user.Email,
			CompanyID:       company.ID,
			CompanyName:     company.Name,
			CompanyStatus:   company.Status,
			UserRole:        user.Role,
			Token:           token,
			LicenseKey:      "trial-" + fmt.Sprint(company.ID),
			LicenseType:     "trial",
			Company: dto.CompanyDTO{
				ID:     company.ID,
				Name:   company.Name,
				Status: company.Status,
			},
		},
	})
}

// ForgotPassword initiates password reset
func ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Find user
	var user models.SaasUser
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		// Don't reveal if email exists
		c.JSON(http.StatusOK, dto.APIResponse{
			Message: "If email exists, password reset link has been sent",
			Data:    nil,
		})
		return
	}

	// Generate reset token
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate reset token",
		})
		return
	}

	resetToken := hex.EncodeToString(token)

	// Create password reset record
	passwordReset := models.PasswordReset{
		UserID:     user.ID,
		Email:      user.Email,
		ResetToken: resetToken,
		ExpiresAt:  time.Now().Add(1 * time.Hour),
		Status:     "pending",
	}

	if err := database.DB.Create(&passwordReset).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create password reset",
		})
		return
	}

	// TODO: Send email with reset link
	resetLink := os.Getenv("FRONTEND_URL") + "/auth/reset-password?token=" + resetToken

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Password reset link has been sent to your email",
		Data: map[string]string{
			"resetLink": resetLink, // Remove in production
		},
	})
}

// ResetPassword resets user password with token
func ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Validate passwords match
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Passwords do not match",
		})
		return
	}

	// Find password reset record
	var passwordReset models.PasswordReset
	if err := database.DB.Where("reset_token = ? AND status = ? AND expires_at > ?", req.Token, "pending", time.Now()).First(&passwordReset).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or expired reset token",
		})
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to process password",
		})
		return
	}

	// Update user password
	if err := database.DB.Model(&models.SaasUser{}).Where("id = ?", passwordReset.UserID).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update password",
		})
		return
	}

	// Mark reset token as used
	now := time.Now()
	database.DB.Model(&passwordReset).Updates(map[string]interface{}{
		"status":  "used",
		"used_at": now,
	})

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Password reset successfully",
		Data:    nil,
	})
}

// SaasLogout logs out user
func SaasLogout(c *gin.Context) {
	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Logged out successfully",
		Data:    nil,
	})
}

// generateJWT generates JWT token for SaaS user
func generateJWT(userID, companyID uint, email string) string {
	claims := jwt.MapClaims{
		"user_id":    userID,
		"company_id": companyID,
		"email":      email,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}

	var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("your-secret-key-change-in-production")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

// GetCurrentUser returns the authenticated user and their company
func GetCurrentUser(c *gin.Context) {
	// Extract user ID from context (set by middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	userID := uint(userIDInterface.(float64))

	// Fetch user
	var user models.SaasUser
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	// Fetch company
	var company models.Company
	if err := database.DB.First(&company, user.CompanyID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Company not found",
		})
		return
	}

	// Fetch subscription info
	var subscription models.Subscription
	database.DB.Where("company_id = ? AND status = ?", company.ID, "active").First(&subscription)

	// Fetch subscription plan for trial days
	var plan models.SubscriptionPlan
	if subscription.ID > 0 {
		database.DB.First(&plan, subscription.PlanID)
	}

	// Calculate trial days remaining
	trialDaysRemaining := 0
	if company.Status == "trial" && !subscription.CurrentPeriodEnd.IsZero() {
		daysRemaining := int(subscription.CurrentPeriodEnd.Sub(time.Now()).Hours() / 24)
		if daysRemaining > 0 {
			trialDaysRemaining = daysRemaining
		}
	}

	response := dto.APIResponse{
		Message: "Current user fetched",
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"id":         user.ID,
				"companyId":  user.CompanyID,
				"email":      user.Email,
				"fullName":   user.FullName,
				"role":       user.Role,
				"status":     user.Status,
				"joinedDate": user.JoinedDate,
				"lastLogin":  user.LastLogin,
			},
			"company": map[string]interface{}{
				"id":                   company.ID,
				"name":                 company.Name,
				"status":               company.Status,
				"createdAt":            company.CreatedAt,
				"updatedAt":            company.UpdatedAt,
				"trialDaysRemaining":   trialDaysRemaining,
				"subscriptionEndDate":  subscription.CurrentPeriodEnd,
			},
		},
	}

	c.JSON(http.StatusOK, response)
}
