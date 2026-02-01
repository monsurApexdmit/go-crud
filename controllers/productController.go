package controllers

import (
	"encoding/json"
	"go-crud/database"
	"go-crud/models"
	"go-crud/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ==========================
// LIST ALL PRODUCTS
// ==========================
func ListProducts(c *gin.Context) {
	var products []models.Product

	// Build query with preloads
	query := database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location")

	// Optional filters
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if vendorID := c.Query("vendor_id"); vendorID != "" {
		query = query.Where("vendor_id = ?", vendorID)
	}
	if locationID := c.Query("location_id"); locationID != "" {
		query = query.Where("location_id = ?", locationID)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products retrieved successfully",
		"data":    products,
	})
}

// ==========================
// GET SINGLE PRODUCT
// ==========================
func GetProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		First(&product, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product fetched successfully",
		"data":    product,
	})
}

// ==========================
// CREATE PRODUCT
// ==========================
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Parse form data
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Validation
	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		path, err := utils.SaveUploadedFile(c, file, "uploads/products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "details": err.Error()})
			return
		}
		product.Image = path
	}

	// Parse attributes (JSON array of IDs)
	if attributesJSON := c.PostForm("attributes"); attributesJSON != "" {
		var attributeIDs []uint
		if err := json.Unmarshal([]byte(attributesJSON), &attributeIDs); err == nil {
			var attributes []models.Attribute
			for _, id := range attributeIDs {
				attributes = append(attributes, models.Attribute{ID: id})
			}
			product.Attributes = attributes
		}
	}

	// Parse variants (JSON array)
	if variantsJSON := c.PostForm("variants"); variantsJSON != "" {
		var variants []models.ProductVariant
		if err := json.Unmarshal([]byte(variantsJSON), &variants); err == nil {
			product.Variants = variants
		}
	}

	// Create product with all nested relations in a transaction
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create the product
		if err := tx.Create(&product).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	// Reload with all relations
	database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		First(&product, product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    product,
	})
}

// ==========================
// UPDATE PRODUCT
// ==========================
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	updates := make(map[string]interface{})

	// Parse basic fields from form
	if v := c.PostForm("name"); v != "" {
		updates["name"] = v
	}
	if v := c.PostForm("description"); v != "" {
		updates["description"] = v
	}
	if v := c.PostForm("category_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			updates["category_id"] = uint(id)
		}
	}
	if v := c.PostForm("vendor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			updates["vendor_id"] = uint(id)
		}
	}
	if v := c.PostForm("location_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			updates["location_id"] = uint(id)
		}
	}
	if v := c.PostForm("price"); v != "" {
		if price, err := strconv.ParseFloat(v, 64); err == nil {
			updates["price"] = price
		}
	}
	if v := c.PostForm("sale_price"); v != "" {
		if price, err := strconv.ParseFloat(v, 64); err == nil {
			updates["sale_price"] = price
		}
	}
	if v := c.PostForm("stock"); v != "" {
		if stock, err := strconv.Atoi(v); err == nil {
			updates["stock"] = stock
		}
	}
	if v := c.PostForm("sku"); v != "" {
		updates["sku"] = v
	}
	if v := c.PostForm("barcode"); v != "" {
		updates["barcode"] = v
	}
	if v := c.PostForm("published"); v != "" {
		updates["published"] = v == "true"
	}
	if v := c.PostForm("receipt_number"); v != "" {
		updates["receipt_number"] = v
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		// Delete old image if exists
		if product.Image != "" {
			_ = os.Remove(product.Image)
		}

		path, err := utils.SaveUploadedFile(c, file, "uploads/products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
			return
		}
		updates["image"] = path
	}

	// Handle nested updates in a transaction
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Update basic fields
		if len(updates) > 0 {
			if err := tx.Model(&product).Updates(updates).Error; err != nil {
				return err
			}
		}

		// Handle attributes update if provided
		if attributesJSON := c.PostForm("attributes"); attributesJSON != "" {
			var attributeIDs []uint
			if err := json.Unmarshal([]byte(attributesJSON), &attributeIDs); err == nil {
				// Clear existing associations
				if err := tx.Model(&product).Association("Attributes").Clear(); err != nil {
					return err
				}
				// Add new associations
				if len(attributeIDs) > 0 {
					var attributes []models.Attribute
					for _, id := range attributeIDs {
						attributes = append(attributes, models.Attribute{ID: id})
					}
					if err := tx.Model(&product).Association("Attributes").Append(attributes); err != nil {
						return err
					}
				}
			}
		}

		// Handle variants update if provided
		if variantsJSON := c.PostForm("variants"); variantsJSON != "" {
			var variants []models.ProductVariant
			if err := json.Unmarshal([]byte(variantsJSON), &variants); err == nil && len(variants) > 0 {
				// Delete existing variants (cascade will handle inventory)
				if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductVariant{}).Error; err != nil {
					return err
				}
				// Create new variants
				for i := range variants {
					variants[i].ProductID = product.ID
				}
				if err := tx.Create(&variants).Error; err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	// Reload with all relations
	database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		First(&product, product.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    product,
	})
}

// ==========================
// DELETE PRODUCT
// ==========================
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Soft delete (GORM will handle cascade via database constraints)
	if err := database.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
