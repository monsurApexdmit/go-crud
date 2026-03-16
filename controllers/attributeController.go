package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/dto"
	"go-crud/middlewares"
	"go-crud/models"
	"go-crud/utils"
)

// ==========================
// LIST ALL ATTRIBUTES (WITH PAGINATION & SEARCH)
// ==========================
func ListAttributes(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	search := c.Query("search")
	optionType := c.Query("option_type")
	includeInactive := c.DefaultQuery("include_inactive", "true")  // Changed to "true" to show all by default
	sortBy := c.DefaultQuery("sort_by", "sort_order") // sort_order, name, created_at
	sortOrder := c.DefaultQuery("sort_order", "asc")   // asc, desc

	// Calculate offset
	offset := (page - 1) * limit

	// Build query
	companyID, _ := middlewares.GetCompanyID(c)
	query := database.DB.Model(&models.Attribute{}).Where("company_id = ?", companyID)

	// Search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name LIKE ? OR display_name LIKE ? OR description LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Option type filter
	if optionType != "" {
		query = query.Where("option_type = ?", optionType)
	}

	// Status filter
	if includeInactive == "false" {
		query = query.Where("status = ?", true)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply sorting
	validSortFields := map[string]bool{
		"sort_order": true,
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}

	if validSortFields[sortBy] {
		orderClause := sortBy + " " + sortOrder
		query = query.Order(orderClause)
	} else {
		query = query.Order("sort_order ASC, name ASC")
	}

	// Get paginated results
	var attributes []models.Attribute
	if err := query.Limit(limit).Offset(offset).Find(&attributes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attributes"})
		return
	}

	// Calculate pagination info
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasNext := page < totalPages
	hasPrevious := page > 1

	c.JSON(http.StatusOK, gin.H{
		"message": "Attributes retrieved successfully",
		"data":    attributes,
		"pagination": gin.H{
			"total":        total,
			"page":         page,
			"limit":        limit,
			"total_pages":  totalPages,
			"has_next":     hasNext,
			"has_previous": hasPrevious,
		},
	})
}

// ==========================
// GET ALL ATTRIBUTES (SIMPLE - FOR DROPDOWNS)
// ==========================
func GetAllAttributesSimple(c *gin.Context) {
	var attributes []models.Attribute

	companyID, _ := middlewares.GetCompanyID(c)
	// Get only active attributes sorted by sort_order
	if err := database.DB.Where("company_id = ? AND status = ?", companyID, true).
		Order("sort_order ASC, name ASC").
		Select("id, name, display_name, option_type").
		Find(&attributes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attributes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attributes retrieved successfully",
		"data":    attributes,
	})
}

// ==========================
// GET ATTRIBUTE STATISTICS
// ==========================
func GetAttributeStats(c *gin.Context) {
	var stats struct {
		Total              int64            `json:"total"`
		Active             int64            `json:"active"`
		Inactive           int64            `json:"inactive"`
		Required           int64            `json:"required"`
		ByType             map[string]int64 `json:"by_type"`
	}

	stats.ByType = make(map[string]int64)

	companyID, _ := middlewares.GetCompanyID(c)

	// Total count
	database.DB.Model(&models.Attribute{}).Where("company_id = ?", companyID).Count(&stats.Total)

	// Active/Inactive count
	database.DB.Model(&models.Attribute{}).Where("company_id = ? AND status = ?", companyID, true).Count(&stats.Active)
	database.DB.Model(&models.Attribute{}).Where("company_id = ? AND status = ?", companyID, false).Count(&stats.Inactive)

	// Required count
	database.DB.Model(&models.Attribute{}).Where("company_id = ? AND is_required = ?", companyID, true).Count(&stats.Required)

	// Count by option type
	var typeCounts []struct {
		OptionType string
		Count      int64
	}
	database.DB.Model(&models.Attribute{}).
		Where("company_id = ?", companyID).
		Select("option_type, count(*) as count").
		Group("option_type").
		Scan(&typeCounts)

	for _, tc := range typeCounts {
		stats.ByType[tc.OptionType] = tc.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute statistics retrieved successfully",
		"data":    stats,
	})
}

// ==========================
// GET SINGLE ATTRIBUTE
// ==========================
func GetAttribute(c *gin.Context) {
	id := c.Param("id")
	var attribute models.Attribute

	companyID, _ := middlewares.GetCompanyID(c)
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&attribute).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute fetched successfully",
		"data":    attribute,
	})
}

// ==========================
// CREATE ATTRIBUTE
// ==========================
func CreateAttribute(c *gin.Context) {
	// Validate request using DTO
	var request dto.CreateAttributeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleValidationErrors(c, err)
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)

	// Check for duplicate name
	var existingAttribute models.Attribute
	if err := database.DB.Where("name = ? AND company_id = ?", strings.ToLower(strings.TrimSpace(request.Name)), companyID).
		First(&existingAttribute).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Attribute with this name already exists",
			"field": "name",
		})
		return
	}

	// Validate that values are provided for dropdown/radio/checkbox types
	requiresValues := map[string]bool{
		"dropdown": true,
		"radio":    true,
		"checkbox": true,
	}

	if requiresValues[request.OptionType] && strings.TrimSpace(request.Values) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Values are required for " + request.OptionType + " type",
			"field": "values",
		})
		return
	}

	// Create attribute
	attribute := models.Attribute{
		Name:        strings.ToLower(strings.TrimSpace(request.Name)),
		DisplayName: strings.TrimSpace(request.DisplayName),
		OptionType:  request.OptionType,
		Values:      strings.TrimSpace(request.Values),
		Description: strings.TrimSpace(request.Description),
		IsRequired:  false,
		Status:      true,
		SortOrder:   0,
		CompanyID:   companyID,
	}

	// Set optional fields
	if request.IsRequired != nil {
		attribute.IsRequired = *request.IsRequired
	}
	if request.Status != nil {
		attribute.Status = *request.Status
	}
	if request.SortOrder != nil {
		attribute.SortOrder = *request.SortOrder
	}

	if err := database.DB.Create(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attribute"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Attribute created successfully",
		"data":    attribute,
	})
}

// ==========================
// UPDATE ATTRIBUTE
// ==========================
func UpdateAttribute(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)

	// Find existing attribute
	var attribute models.Attribute
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&attribute).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	// Parse raw JSON to detect which fields were sent
	var rawRequest map[string]interface{}
	if err := c.ShouldBindJSON(&rawRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Extract fields from raw request
	if val, exists := rawRequest["name"]; exists && val != nil {
		if strVal, ok := val.(string); ok {
			name := strings.ToLower(strings.TrimSpace(strVal))

			// Check for duplicate name
			if name != attribute.Name {
				var existingAttribute models.Attribute
				if err := database.DB.Where("name = ? AND company_id = ? AND id != ?", name, companyID, id).
					First(&existingAttribute).Error; err == nil {
					c.JSON(http.StatusConflict, gin.H{
						"error": "Attribute with this name already exists",
						"field": "name",
					})
					return
				}
				attribute.Name = name
			}
		}
	}

	if val, exists := rawRequest["display_name"]; exists && val != nil {
		if strVal, ok := val.(string); ok {
			attribute.DisplayName = strings.TrimSpace(strVal)
		}
	}

	if val, exists := rawRequest["option_type"]; exists && val != nil {
		if strVal, ok := val.(string); ok {
			attribute.OptionType = strVal
		}
	}

	if val, exists := rawRequest["values"]; exists && val != nil {
		if strVal, ok := val.(string); ok {
			attribute.Values = strings.TrimSpace(strVal)
		}
	}

	if val, exists := rawRequest["description"]; exists && val != nil {
		if strVal, ok := val.(string); ok {
			attribute.Description = strings.TrimSpace(strVal)
		}
	}

	if val, exists := rawRequest["is_required"]; exists && val != nil {
		if boolVal, ok := val.(bool); ok {
			attribute.IsRequired = boolVal
		}
	}

	if val, exists := rawRequest["status"]; exists && val != nil {
		if boolVal, ok := val.(bool); ok {
			attribute.Status = boolVal
		}
	}

	if val, exists := rawRequest["sort_order"]; exists && val != nil {
		if floatVal, ok := val.(float64); ok {
			attribute.SortOrder = int(floatVal)
		}
	}

	// Validate that values are provided for dropdown/radio/checkbox types
	requiresValues := map[string]bool{
		"dropdown": true,
		"radio":    true,
		"checkbox": true,
	}

	if requiresValues[attribute.OptionType] && strings.TrimSpace(attribute.Values) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Values are required for " + attribute.OptionType + " type",
			"field": "values",
		})
		return
	}

	// Save updates
	if err := database.DB.Save(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attribute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute updated successfully",
		"data":    attribute,
	})
}

// ==========================
// TOGGLE ATTRIBUTE STATUS
// ==========================
func ToggleAttributeStatus(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)
	var attribute models.Attribute
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&attribute).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	// Toggle status
	attribute.Status = !attribute.Status

	if err := database.DB.Save(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attribute status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute status updated successfully",
		"data":    attribute,
	})
}

// ==========================
// DELETE ATTRIBUTE
// ==========================
func DeleteAttribute(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)

	// Check if attribute exists
	var attribute models.Attribute
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&attribute).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute not found"})
		return
	}

	// TODO: Check if attribute is being used by any products
	// For now, we'll allow deletion

	// Soft delete the attribute
	if err := database.DB.Delete(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attribute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute deleted successfully",
	})
}

// ==========================
// BULK DELETE ATTRIBUTES
// ==========================
func BulkDeleteAttributes(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)

	// TODO: Check if any attributes are being used by products
	// For now, we'll allow deletion

	// Delete attributes
	if err := database.DB.Where("company_id = ?", companyID).Delete(&models.Attribute{}, request.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attributes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attributes deleted successfully",
		"deleted": len(request.IDs),
	})
}
