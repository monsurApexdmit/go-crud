package controllers

import (
	"encoding/json"
	"fmt"
	"go-crud/database"
	"go-crud/models"
	"go-crud/utils"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// tempFile holds temporary and final file paths for uploads
type tempFile struct {
	tempPath  string
	finalPath string
}

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
		Preload("Variants.Inventory.Location").
		Preload("Images")

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

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var total int64
	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count products"})
		return
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products retrieved successfully",
		"page":    page,
		"limit":   limit,
		"total":   total,
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
		Preload("Images").
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
	type productForm struct {
		Name          string  `form:"name"`
		Description   string  `form:"description"`
		CategoryID    *uint   `form:"category_id"`
		VendorID      *uint   `form:"vendor_id"`
		LocationID    *uint   `form:"location_id"`
		Price         float64 `form:"price"`
		SalePrice     float64 `form:"sale_price"`
		Stock         int     `form:"stock"`
		SKU           string  `form:"sku"`
		Barcode       string  `form:"barcode"`
		Published     bool    `form:"published"`
		ReceiptNumber string  `form:"receipt_number"`
	}

	var form productForm
	var product models.Product
	var tempImagePath string
	var finalImagePath string
	var tempImages []tempFile

	// Parse form data
	if err := c.ShouldBind(&form); err != nil {
		log.Printf("CreateProduct: bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	product = models.Product{
		Name:          form.Name,
		Description:   form.Description,
		CategoryID:    form.CategoryID,
		VendorID:      form.VendorID,
		LocationID:    form.LocationID,
		Price:         form.Price,
		SalePrice:     form.SalePrice,
		Stock:         form.Stock,
		SKU:           form.SKU,
		Barcode:       form.Barcode,
		Published:     form.Published,
		ReceiptNumber: form.ReceiptNumber,
	}

	// Validation
	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product name is required"})
		return
	}

	// Handle image upload
	file, err := c.FormFile("image")
	if err == nil {
		tempPath, finalPath, err := utils.SaveUploadedFileTemp(c, file, "uploads/products")
		if err != nil {
			log.Printf("CreateProduct: save image error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "details": err.Error()})
			return
		}
		tempImagePath = tempPath
		finalImagePath = finalPath
	}

	// Handle multiple images upload (images, images[], image[0], images[0], etc.)
	if formData, err := c.MultipartForm(); err == nil && formData != nil {
		// Use a map to track unique files and avoid duplicates
		fileMap := make(map[string]*multipart.FileHeader)

		// Collect all files from different field names
		for key, fhs := range formData.File {
			// Match any field starting with "image" (includes image, images, image[0], image[1], etc.)
			if strings.HasPrefix(key, "image") {
				for _, fh := range fhs {
					// Use filename + size as unique key to avoid processing same file twice
					uniqueKey := fmt.Sprintf("%s_%d", fh.Filename, fh.Size)
					if _, exists := fileMap[uniqueKey]; !exists {
						fileMap[uniqueKey] = fh
					}
				}
			}
		}

		// Process unique files
		log.Printf("CreateProduct: Found %d unique image files to process", len(fileMap))
		for _, f := range fileMap {
			tp, fp, err := utils.SaveUploadedFileTemp(c, f, "uploads/products")
			if err != nil {
				log.Printf("CreateProduct: save multi image error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image", "details": err.Error()})
				return
			}
			log.Printf("CreateProduct: Saved temp image: %s -> %s", tp, fp)
			tempImages = append(tempImages, tempFile{tempPath: tp, finalPath: fp})
		}
	} else if err != nil {
		log.Printf("CreateProduct: multipart form error: %v", err)
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
		} else {
			log.Printf("CreateProduct: attributes JSON parse error: %v", err)
		}
	}

	// Parse variants (JSON array)
	if variantsJSON := c.PostForm("variants"); variantsJSON != "" {
		type variantInput struct {
			Name       string          `json:"name"`
			Attributes json.RawMessage `json:"attributes"`
			Price      float64         `json:"price"`
			SalePrice  float64         `json:"sale_price"`
			Stock      int             `json:"stock"`
			SKU        string          `json:"sku"`
			Barcode    string          `json:"barcode"`
		}

		var inputs []variantInput
		if err := json.Unmarshal([]byte(variantsJSON), &inputs); err == nil {
			var variants []models.ProductVariant
			for _, v := range inputs {
				variants = append(variants, models.ProductVariant{
					Name:       v.Name,
					Attributes: datatypes.JSON(v.Attributes),
					Price:      v.Price,
					SalePrice:  v.SalePrice,
					Stock:      v.Stock,
					SKU:        v.SKU,
					Barcode:    v.Barcode,
				})
			}
			product.Variants = variants
		} else {
			log.Printf("CreateProduct: variants JSON parse error: %v", err)
		}
	}

	// Create product with all nested relations in a transaction
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create the product
		if err := tx.Create(&product).Error; err != nil {
			log.Printf("CreateProduct: db create error: %v", err)
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	if tempImagePath != "" {
		if err := utils.CommitUploadedFile(tempImagePath, finalImagePath); err != nil {
			_ = os.Remove(tempImagePath)
			log.Printf("CreateProduct: commit image error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store image", "details": err.Error()})
			return
		}

		if err := database.DB.Model(&product).Update("image", finalImagePath).Error; err != nil {
			log.Printf("CreateProduct: image db update error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image", "details": err.Error()})
			return
		}
		product.Image = finalImagePath
	}

	if len(tempImages) > 0 {
		log.Printf("CreateProduct: Committing %d images to database", len(tempImages))
		var imageRows []models.ProductImage
		for i, tf := range tempImages {
			if err := utils.CommitUploadedFile(tf.tempPath, tf.finalPath); err != nil {
				_ = os.Remove(tf.tempPath)
				log.Printf("CreateProduct: commit multi image error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store image", "details": err.Error()})
				return
			}
			log.Printf("CreateProduct: Committed image %d: %s", i, tf.finalPath)
			imageRows = append(imageRows, models.ProductImage{
				ProductID: product.ID,
				Path:      tf.finalPath,
				Position:  i,
				IsPrimary: i == 0,
			})
		}
		if err := database.DB.Create(&imageRows).Error; err != nil {
			for _, r := range imageRows {
				_ = os.Remove(r.Path)
			}
			log.Printf("CreateProduct: images db create error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save images", "details": err.Error()})
			return
		}
		log.Printf("CreateProduct: Successfully saved %d images to database", len(imageRows))
		product.Images = imageRows
	}

	// Reload with all relations
	database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		Preload("Images").
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
	oldImagePath := product.Image
	var tempImagePath string
	var finalImagePath string
	var oldImagePaths []string
	var tempImages []tempFile

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
		tempPath, finalPath, err := utils.SaveUploadedFileTemp(c, file, "uploads/products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
			return
		}
		tempImagePath = tempPath
		finalImagePath = finalPath
	}

	if formData, err := c.MultipartForm(); err == nil && formData != nil {
		// Use a map to track unique files and avoid duplicates
		fileMap := make(map[string]*multipart.FileHeader)

		// Collect all files from different field names
		for key, fhs := range formData.File {
			// Match any field starting with "image" (includes image, images, image[0], image[1], etc.)
			if strings.HasPrefix(key, "image") {
				for _, fh := range fhs {
					// Use filename + size as unique key to avoid processing same file twice
					uniqueKey := fmt.Sprintf("%s_%d", fh.Filename, fh.Size)
					if _, exists := fileMap[uniqueKey]; !exists {
						fileMap[uniqueKey] = fh
					}
				}
			}
		}

		// Process unique files
		for _, f := range fileMap {
			tp, fp, err := utils.SaveUploadedFileTemp(c, f, "uploads/products")
			if err != nil {
				log.Printf("UpdateProduct: save multi image error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
				return
			}
			tempImages = append(tempImages, tempFile{tempPath: tp, finalPath: fp})
		}
	} else if err != nil {
		log.Printf("UpdateProduct: multipart form error: %v", err)
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

	if tempImagePath != "" {
		if err := utils.CommitUploadedFile(tempImagePath, finalImagePath); err != nil {
			_ = os.Remove(tempImagePath)
			log.Printf("UpdateProduct: commit image error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store image", "details": err.Error()})
			return
		}

		if err := database.DB.Model(&product).Update("image", finalImagePath).Error; err != nil {
			log.Printf("UpdateProduct: image db update error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image", "details": err.Error()})
			return
		}

		if oldImagePath != "" {
			_ = os.Remove(oldImagePath)
		}
		product.Image = finalImagePath
	}

	if len(tempImages) > 0 {
		if err := database.DB.Model(&models.ProductImage{}).
			Where("product_id = ?", product.ID).
			Pluck("path", &oldImagePaths).Error; err != nil {
			log.Printf("UpdateProduct: fetch old images error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update images", "details": err.Error()})
			return
		}

		var imageRows []models.ProductImage
		for i, tf := range tempImages {
			if err := utils.CommitUploadedFile(tf.tempPath, tf.finalPath); err != nil {
				_ = os.Remove(tf.tempPath)
				log.Printf("UpdateProduct: commit multi image error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store image", "details": err.Error()})
				return
			}
			imageRows = append(imageRows, models.ProductImage{
				ProductID: product.ID,
				Path:      tf.finalPath,
				Position:  i,
				IsPrimary: i == 0,
			})
		}

		if err := database.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductImage{}).Error; err != nil {
				return err
			}
			if err := tx.Create(&imageRows).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			for _, r := range imageRows {
				_ = os.Remove(r.Path)
			}
			log.Printf("UpdateProduct: images db update error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update images", "details": err.Error()})
			return
		}

		for _, p := range oldImagePaths {
			_ = os.Remove(p)
		}
		product.Images = imageRows
	}

	// Reload with all relations
	database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		Preload("Images").
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

	// First, fetch the product with all images
	var product models.Product
	if err := database.DB.Preload("Images").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Collect all image paths to delete
	var imagePaths []string

	// Add main product image if exists
	if product.Image != "" {
		imagePaths = append(imagePaths, product.Image)
	}

	// Add all additional images
	for _, img := range product.Images {
		if img.Path != "" {
			imagePaths = append(imagePaths, img.Path)
		}
	}

	// Soft delete the product (GORM will handle cascade for related records)
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	// Delete all physical image files
	for _, path := range imagePaths {
		if err := os.Remove(path); err != nil {
			// Log the error but don't fail the request
			log.Printf("DeleteProduct: failed to delete file %s: %v", path, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
