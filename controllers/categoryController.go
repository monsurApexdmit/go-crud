package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"go-crud/database"
	"go-crud/dto"
	"go-crud/middlewares"
	"go-crud/models"
	"go-crud/utils"
)


// ==========================
// GET ALL CATEGORIES (WITH PAGINATION & SEARCH)
// ==========================
func ListCategories(c *gin.Context) {
	// Get query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	search := c.Query("search")
	view := c.DefaultQuery("view", "tree") // tree, flat, all
	includeInactive := c.DefaultQuery("include_inactive", "false")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	// Get company_id from context
	companyID, _ := middlewares.GetCompanyID(c)
	query := database.DB.Model(&models.Category{}).Where("company_id = ?", companyID)

	// Search filter
	if search != "" {
		query = query.Where("category_name LIKE ?", "%"+search+"%")
	}

	// Status filter
	if includeInactive != "true" {
		query = query.Where("status = ?", true)
	}

	// Different view modes
	var categories []models.Category
	var total int64

	switch view {
	case "flat":
		// Flat list with pagination
		query.Count(&total)
		err := query.
			Preload("Parent").
			Offset(offset).
			Limit(limit).
			Order("created_at DESC").
			Find(&categories).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

	case "all":
		// All categories without tree structure
		query.Count(&total)
		err := query.
			Offset(offset).
			Limit(limit).
			Order("category_name ASC").
			Find(&categories).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}

	default: // tree
		// Tree structure (root categories with children)
		rootQuery := query.Where("parent_id IS NULL")
		rootQuery.Count(&total)

		err := rootQuery.
			Preload("Children", func(db *gorm.DB) *gorm.DB {
				if includeInactive != "true" {
					return db.Where("status = ?", true)
				}
				return db
			}).
			Offset(offset).
			Limit(limit).
			Order("category_name ASC").
			Find(&categories).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
			return
		}
	}

	// Calculate pagination metadata
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Page-Count", strconv.Itoa(totalPages))
	c.Header("X-Current-Page", strconv.Itoa(page))

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"data":    categories,
		"pagination": gin.H{
			"total":        total,
			"page":         page,
			"limit":        limit,
			"total_pages":  totalPages,
			"has_next":     page < totalPages,
			"has_previous": page > 1,
		},
	})
}


// ==========================
// GET ALL CATEGORIES (SIMPLE)
// ==========================
func GetAllCategoriesSimple(c *gin.Context) {
	var categories []models.Category

	companyID, _ := middlewares.GetCompanyID(c)
	err := database.DB.
		Select("id, category_name, parent_id, status").
		Where("company_id = ? AND status = ?", companyID, true).
		Order("category_name ASC").
		Find(&categories).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories retrieved successfully",
		"data":    categories,
	})
}


// ==========================
// GET SINGLE CATEGORY
// ==========================
func GetCategory(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)
	var category models.Category
	if err := database.DB.
		Where("id = ? AND company_id = ?", id, companyID).
		Preload("Parent").
		Preload("Children").
		First(&category).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category fetched successfully",
		"data":    category,
	})
}


// ==========================
// CREATE CATEGORY
// ==========================
func CreateCategory(c *gin.Context) {
	// Validate request using DTO
	var request dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleValidationErrors(c, err)
		return
	}

	// Get company_id from context
	companyID, _ := middlewares.GetCompanyID(c)

	// Validate parent category exists if provided
	if request.ParentID != nil {
		var parentCategory models.Category
		if err := database.DB.Where("id = ? AND company_id = ?", *request.ParentID, companyID).First(&parentCategory).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Parent category not found",
				"field": "parent_id",
			})
			return
		}
	}

	// Check for duplicate category name
	var existingCategory models.Category
	if err := database.DB.Where("category_name = ? AND company_id = ?", request.CategoryName, companyID).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Category with this name already exists",
			"field": "category_name",
		})
		return
	}

	// Map DTO to Model
	category := models.Category{
		CategoryName: strings.TrimSpace(request.CategoryName),
		ParentID:     request.ParentID,
		CompanyID:    companyID,
	}

	// Set status (default to true if not provided)
	if request.Status != nil {
		category.Status = *request.Status
	} else {
		category.Status = true
	}

	// Create category in database
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	// Reload with parent data for response
	database.DB.Preload("Parent").First(&category, category.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Category created successfully",
		"data":    category,
	})
}


// ==========================
// UPDATE CATEGORY
// ==========================
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	// Get company_id from context
	companyID, _ := middlewares.GetCompanyID(c)

	// Find existing category
	var category models.Category
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Parse raw JSON to detect which fields were actually sent
	var rawRequest map[string]interface{}
	if err := c.ShouldBindJSON(&rawRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Check if parent_id field was explicitly sent (even if null)
	parentIDProvided := false
	var parentIDValue *uint
	if val, exists := rawRequest["parent_id"]; exists {
		parentIDProvided = true
		// If value is not nil, try to convert it
		if val != nil {
			if floatVal, ok := val.(float64); ok {
				uintVal := uint(floatVal)
				parentIDValue = &uintVal
			}
		}
		// If val is nil, parentIDValue stays nil (which is what we want)
	}

	// Validate parent category exists if provided and not null
	if parentIDProvided && parentIDValue != nil {
		categoryID, _ := strconv.ParseUint(id, 10, 32)

		// Prevent category from being its own parent
		if *parentIDValue == uint(categoryID) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Category cannot be its own parent",
				"field": "parent_id",
			})
			return
		}

		var parentCategory models.Category
		if err := database.DB.Where("id = ? AND company_id = ?", *parentIDValue, companyID).First(&parentCategory).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Parent category not found",
				"field": "parent_id",
			})
			return
		}

		// Prevent circular reference (if parent is a child of current category)
		if isCircularReference(c, uint(categoryID), *parentIDValue) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot create circular category reference",
				"field": "parent_id",
			})
			return
		}
	}

	// Extract other fields from raw request
	categoryName := ""
	if val, exists := rawRequest["category_name"]; exists && val != nil {
		if strVal, ok := val.(string); ok {
			categoryName = strVal
		}
	}

	statusValue := (*bool)(nil)
	if val, exists := rawRequest["status"]; exists && val != nil {
		if boolVal, ok := val.(bool); ok {
			statusValue = &boolVal
		}
	}

	// Check for duplicate category name if updating name
	if categoryName != "" && categoryName != category.CategoryName {
		var existingCategory models.Category
		if err := database.DB.Where("category_name = ? AND company_id = ? AND id != ?", categoryName, companyID, id).First(&existingCategory).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "Category with this name already exists",
				"field": "category_name",
			})
			return
		}
	}

	// Update only provided fields
	if categoryName != "" {
		category.CategoryName = strings.TrimSpace(categoryName)
	}

	// Handle parent_id update - set to null if provided as null, or set to value if provided
	if parentIDProvided {
		category.ParentID = parentIDValue
	}

	if statusValue != nil {
		category.Status = *statusValue
	}

	// Save updates
	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	// Reload with parent data for response
	database.DB.Preload("Parent").Preload("Children").First(&category, category.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    category,
	})
}


// ==========================
// TOGGLE CATEGORY STATUS
// ==========================
func ToggleCategoryStatus(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)
	var category models.Category
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Toggle status
	category.Status = !category.Status

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category status updated successfully",
		"data":    category,
	})
}


// ==========================
// DELETE CATEGORY
// ==========================
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	companyID, _ := middlewares.GetCompanyID(c)

	// Check if category exists
	var category models.Category
	if err := database.DB.Where("id = ? AND company_id = ?", id, companyID).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Check if category has children
	var childCount int64
	database.DB.Model(&models.Category{}).Where("parent_id = ? AND company_id = ?", id, companyID).Count(&childCount)

	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete category with subcategories. Delete or reassign subcategories first.",
		})
		return
	}

	// Soft delete the category
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Category deleted successfully",
	})
}


// ==========================
// BULK DELETE CATEGORIES
// ==========================
func BulkDeleteCategories(c *gin.Context) {
	var request struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)

	// Check if any category has children
	var childCount int64
	database.DB.Model(&models.Category{}).Where("parent_id IN ? AND company_id = ?", request.IDs, companyID).Count(&childCount)

	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete categories with subcategories",
		})
		return
	}

	// Delete categories
	if err := database.DB.Where("company_id = ?", companyID).Delete(&models.Category{}, request.IDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categories deleted successfully",
		"deleted": len(request.IDs),
	})
}


// ==========================
// GET CATEGORY STATISTICS
// ==========================
func GetCategoryStats(c *gin.Context) {
	var stats struct {
		Total      int64 `json:"total"`
		Active     int64 `json:"active"`
		Inactive   int64 `json:"inactive"`
		Root       int64 `json:"root_categories"`
		WithParent int64 `json:"subcategories"`
	}

	companyID, _ := middlewares.GetCompanyID(c)

	database.DB.Model(&models.Category{}).Where("company_id = ?", companyID).Count(&stats.Total)
	database.DB.Model(&models.Category{}).Where("company_id = ? AND status = ?", companyID, true).Count(&stats.Active)
	database.DB.Model(&models.Category{}).Where("company_id = ? AND status = ?", companyID, false).Count(&stats.Inactive)
	database.DB.Model(&models.Category{}).Where("company_id = ? AND parent_id IS NULL", companyID).Count(&stats.Root)
	database.DB.Model(&models.Category{}).Where("company_id = ? AND parent_id IS NOT NULL", companyID).Count(&stats.WithParent)

	c.JSON(http.StatusOK, gin.H{
		"message": "Category statistics retrieved successfully",
		"data":    stats,
	})
}


// ==========================
// HELPER: Check Circular Reference
// ==========================
func isCircularReference(c *gin.Context, categoryID, parentID uint) bool {
	var parent models.Category

	// Traverse up the parent chain
	currentID := parentID
	companyID, _ := middlewares.GetCompanyID(c)
	for i := 0; i < 10; i++ { // Limit depth to prevent infinite loop
		if err := database.DB.Where("id = ? AND company_id = ?", currentID, companyID).First(&parent).Error; err != nil {
			return false // Parent not found, no circular reference
		}

		if parent.ParentID == nil {
			return false // Reached root, no circular reference
		}

		if *parent.ParentID == categoryID {
			return true // Found circular reference
		}

		currentID = *parent.ParentID
	}

	return false
}
