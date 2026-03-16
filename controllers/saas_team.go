package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"go-crud/database"
	"go-crud/dto"
	"go-crud/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// GetTeamCompanyIDFromRequest extracts company_id from context or query param
func GetTeamCompanyIDFromRequest(c *gin.Context) (uint, error) {
	// Try to get from middleware context first
	if id, exists := c.Get("company_id"); exists {
		return uint(id.(float64)), nil
	}

	// Try to get from query parameter
	companyIDStr := c.Query("company_id")
	if companyIDStr == "" {
		return 0, nil // Return 0 if not provided
	}

	parsedID, err := strconv.ParseUint(companyIDStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(parsedID), nil
}

// GetTeamUsers gets all team members
func GetTeamUsers(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetTeamCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var users []models.SaasUser
	if err := database.DB.Where("company_id = ?", companyID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to fetch team members",
		})
		return
	}

	// Get subscription plan for max users
	var subscription models.Subscription
	var maxUsers int = 5
	if err := database.DB.Where("company_id = ?", companyID).First(&subscription).Error; err == nil {
		var plan models.SubscriptionPlan
		if err := database.DB.First(&plan, subscription.PlanID).Error; err == nil {
			maxUsers = plan.MaxUsers
		}
	}

	userCount := len(users)
	canAddMore := userCount < maxUsers

	// Convert to DTOs
	var userDTOs []dto.TeamUserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dto.TeamUserDTO{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Role:      user.Role,
			Status:    user.Status,
			JoinedDate: user.JoinedDate,
			LastLogin: user.LastLogin,
			Avatar:    user.Avatar,
		})
	}

	c.JSON(http.StatusOK, dto.TeamUsersResponse{
		Message: "Team members retrieved",
		Data: dto.TeamUsersData{
			Users:      userDTOs,
			TotalUsers: userCount,
			MaxUsers:   maxUsers,
			CanAddMore: canAddMore,
		},
	})
}

// InviteUser invites a new team member
func InviteUser(c *gin.Context) {
	// Get company_id from context or query parameter
	companyID, err := GetTeamCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	var req dto.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.SaasUser
	if database.DB.Where("email = ? AND company_id = ?", req.Email, uint(companyID)).First(&existingUser).RowsAffected > 0 {
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Message: "User with this email already exists in the company",
		})
		return
	}

	// Check team member limit
	var userCount int64
	database.DB.Model(&models.SaasUser{}).Where("company_id = ?", uint(companyID)).Count(&userCount)

	var subscription models.Subscription
	var maxUsers int = 5
	if err := database.DB.Where("company_id = ?", uint(companyID)).First(&subscription).Error; err == nil {
		var plan models.SubscriptionPlan
		if err := database.DB.First(&plan, subscription.PlanID).Error; err == nil {
			maxUsers = plan.MaxUsers
		}
	}

	if int(userCount) >= maxUsers {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Message: "Team member limit reached",
		})
		return
	}

	// Generate invitation token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate invitation token",
		})
		return
	}
	invitationToken := hex.EncodeToString(tokenBytes)

	// Create invitation
	invitation := models.Invitation{
		CompanyID:         uint(companyID),
		Email:             req.Email,
		FullName:          req.FullName,
		Role:              req.Role,
		Status:            "pending",
		InvitationToken:   invitationToken,
		ExpiresAt:         time.Now().Add(7 * 24 * time.Hour),
		InvitedAt:         time.Now(),
	}

	if err := database.DB.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create invitation",
		})
		return
	}

	// TODO: Send invitation email with acceptance link

	c.JSON(http.StatusCreated, dto.InviteUserResponse{
		Message: "Invitation sent successfully",
		Data: dto.InviteData{
			UserID:          invitation.ID,
			Email:           invitation.Email,
			Status:          invitation.Status,
			InvitationToken: invitation.InvitationToken,
			ExpiresAt:       invitation.ExpiresAt,
			InvitedAt:       invitation.InvitedAt,
		},
	})
}

// UpdateUserRole updates user role
func UpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	// Get company_id from context or query parameter
	companyID, err := GetTeamCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	// Get current user (from token)
	currentUserID, err := strconv.ParseUint(c.GetString("userID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid current user",
		})
		return
	}

	// Check if user is trying to change their own role
	if uint(userID) == uint(currentUserID) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Message: "Cannot change your own role",
		})
		return
	}

	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Find user and verify they belong to the company
	var user models.SaasUser
	if err := database.DB.Where("id = ? AND company_id = ?", uint(userID), uint(companyID)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	// Don't allow changing owner role
	if user.Role == "owner" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Message: "Cannot change owner role",
		})
		return
	}

	// Update role
	if err := database.DB.Model(&user).Update("role", req.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update user role",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UpdateUserRoleResponse{
		Message: "User role updated successfully",
		Data: dto.UpdateRoleData{
			UserID:   user.ID,
			Role:     user.Role,
			UpdatedAt: time.Now(),
		},
	})
}

// RemoveUser removes a team member
func RemoveUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	// Get company_id from context or query parameter
	companyID, err := GetTeamCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	// Get current user (from token)
	currentUserID, err := strconv.ParseUint(c.GetString("userID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid current user",
		})
		return
	}

	// Check if user is trying to remove themselves
	if uint(userID) == uint(currentUserID) {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Message: "Cannot remove yourself",
		})
		return
	}

	// Find user and verify they belong to the company
	var user models.SaasUser
	if err := database.DB.Where("id = ? AND company_id = ?", uint(userID), uint(companyID)).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	// Don't allow removing owner
	if user.Role == "owner" {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Message: "Cannot remove owner",
		})
		return
	}

	// Delete user
	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to remove user",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "User removed successfully",
		Data: map[string]interface{}{
			"userId":   user.ID,
			"success":  true,
			"removedAt": time.Now(),
		},
	})
}

// ResendInvitation resends an invitation
func ResendInvitation(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	// Get company_id from context or query parameter
	companyID, err := GetTeamCompanyIDFromRequest(c)
	if err != nil || companyID == 0 {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or missing company ID",
		})
		return
	}

	// Find invitation
	var invitation models.Invitation
	if err := database.DB.Where("id = ? AND company_id = ?", uint(userID), uint(companyID)).First(&invitation).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Message: "Invitation not found",
		})
		return
	}

	// Generate new invitation token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate new token",
		})
		return
	}
	newToken := hex.EncodeToString(tokenBytes)

	// Update invitation
	updateData := map[string]interface{}{
		"invitation_token": newToken,
		"expires_at":       time.Now().Add(7 * 24 * time.Hour),
		"invited_at":       time.Now(),
	}

	if err := database.DB.Model(&invitation).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to resend invitation",
		})
		return
	}

	// TODO: Send invitation email again

	c.JSON(http.StatusOK, dto.APIResponse{
		Message: "Invitation resent successfully",
		Data: nil,
	})
}

// AcceptInvitation accepts an invitation and creates user
func AcceptInvitation(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request",
		})
		return
	}

	// Find valid invitation
	var invitation models.Invitation
	if err := database.DB.Where("invitation_token = ? AND status = ? AND expires_at > ?", req.Token, "pending", time.Now()).First(&invitation).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid or expired invitation",
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

	// Create user
	user := models.SaasUser{
		CompanyID: invitation.CompanyID,
		Email:     invitation.Email,
		FullName:  invitation.FullName,
		Password:  string(hashedPassword),
		Role:      invitation.Role,
		Status:    "active",
		JoinedDate: time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create user",
		})
		return
	}

	// Mark invitation as accepted
	now := time.Now()
	database.DB.Model(&invitation).Updates(map[string]interface{}{
		"status":       "accepted",
		"accepted_at":  now,
	})

	// Generate token
	token := generateJWT(user.ID, user.CompanyID, user.Email)

	c.JSON(http.StatusCreated, dto.APIResponse{
		Message: "Invitation accepted successfully",
		Data: map[string]interface{}{
			"userId":    user.ID,
			"token":     token,
			"email":     user.Email,
			"role":      user.Role,
			"joinedAt":  user.JoinedDate,
		},
	})
}
