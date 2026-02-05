package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// permissionItem is the JSON shape the frontend sends and receives.
type permissionItem struct {
	Name   string `json:"name"`
	Read   bool   `json:"read"`
	Write  bool   `json:"write"`
	Delete bool   `json:"delete"`
}

// roleRequest is the body for POST / PUT staff-roles.
type roleRequest struct {
	Name        string           `json:"name"`
	Permissions []permissionItem `json:"permissions"`
}

// roleResponse is what we return so the permissions array is always present.
type roleResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Permissions []permissionItem `json:"permissions"`
	CreatedAt   string           `json:"createdAt"`
	UpdatedAt   string           `json:"updatedAt"`
}

func buildRoleResponse(role models.StaffRole, perms []permissionItem) roleResponse {
	if perms == nil {
		perms = []permissionItem{}
	}
	return roleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: perms,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// loadPermissions fetches the permissions array for a role from the DB.
func loadPermissions(roleID uint) []permissionItem {
	var rps []models.RolePermission
	database.DB.
		Preload("Permission").
		Where("role_id = ?", roleID).
		Find(&rps)

	items := make([]permissionItem, 0, len(rps))
	for _, rp := range rps {
		items = append(items, permissionItem{
			Name:   rp.Permission.Name,
			Read:   rp.Read,
			Write:  rp.Write,
			Delete: rp.Delete,
		})
	}
	return items
}

// savePermissions replaces all role_permissions rows for a role inside a transaction.
func savePermissions(tx *gorm.DB, roleID uint, perms []permissionItem) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	if len(perms) == 0 {
		return nil
	}

	// Load permission name → id map
	var allPerms []models.Permission
	if err := tx.Find(&allPerms).Error; err != nil {
		return err
	}
	nameToID := make(map[string]uint, len(allPerms))
	for _, p := range allPerms {
		nameToID[p.Name] = p.ID
	}

	rows := make([]models.RolePermission, 0, len(perms))
	for _, p := range perms {
		pid, ok := nameToID[p.Name]
		if !ok {
			continue // skip unknown module names
		}
		rows = append(rows, models.RolePermission{
			RoleID:       roleID,
			PermissionID: pid,
			Read:         p.Read,
			Write:        p.Write,
			Delete:       p.Delete,
		})
	}

	if len(rows) == 0 {
		return nil
	}
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
}

func ListStaffRoles(c *gin.Context) {
	var roles []models.StaffRole
	if err := database.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve roles"})
		return
	}

	resp := make([]roleResponse, 0, len(roles))
	for _, role := range roles {
		resp = append(resp, buildRoleResponse(role, loadPermissions(role.ID)))
	}
	c.JSON(http.StatusOK, gin.H{"message": "Roles retrieved successfully", "data": resp})
}

func GetStaffRole(c *gin.Context) {
	var role models.StaffRole
	if err := database.DB.First(&role, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role fetched successfully", "data": buildRoleResponse(role, loadPermissions(role.ID))})
}

func CreateStaffRole(c *gin.Context) {
	var req roleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role name is required"})
		return
	}

	role := models.StaffRole{Name: req.Name}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return savePermissions(tx, role.ID, req.Permissions)
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	database.DB.First(&role, role.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Role created successfully", "data": buildRoleResponse(role, loadPermissions(role.ID))})
}

func UpdateStaffRole(c *gin.Context) {
	var role models.StaffRole
	if err := database.DB.First(&role, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var req roleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if req.Name != "" {
			role.Name = req.Name
		}
		if err := tx.Save(&role).Error; err != nil {
			return err
		}
		if req.Permissions != nil {
			return savePermissions(tx, role.ID, req.Permissions)
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	database.DB.First(&role, role.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "data": buildRoleResponse(role, loadPermissions(role.ID))})
}

func DeleteStaffRole(c *gin.Context) {
	var role models.StaffRole
	if err := database.DB.First(&role, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// role_permissions rows are removed by ON DELETE CASCADE;
	// soft-delete the role itself.
	database.DB.Delete(&role)
	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}
