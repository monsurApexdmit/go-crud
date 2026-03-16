package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go-crud/middlewares"
	"go-crud/services"
	"gorm.io/gorm"
)

// ListProducts parses pagination and filters from query params, delegates to the service.
func ListProducts(c *gin.Context) {
	service := services.NewProductService()
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	companyID, _ := middlewares.GetCompanyID(c)
	params := services.ProductListParams{
		Page:  page,
		Limit: limit,
	}

	// Add company_id to filters
	params.Filters.CompanyID = &companyID

	if v := c.Query("category_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid := uint(id)
			params.Filters.CategoryID = &uid
		}
	}
	if v := c.Query("vendor_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid := uint(id)
			params.Filters.VendorID = &uid
		}
	}
	if v := c.Query("location_id"); v != "" {
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			uid := uint(id)
			params.Filters.LocationID = &uid
		}
	}

	products, total, err := service.ListProducts(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Products retrieved successfully",
		"page":    params.Page,
		"limit":   params.Limit,
		"total":   total,
		"data":    products,
	})
}

// GetProduct loads a single product by :id.
func GetProduct(c *gin.Context) {
	service := services.NewProductService()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)
	product, err := service.GetProductByCompany(uint(id), companyID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product fetched successfully",
		"data":    product,
	})
}

// CreateProduct binds the product form, extracts uploaded files, and delegates
// the full create flow to the service.
func CreateProduct(c *gin.Context) {
	service := services.NewProductService()
	var form services.ProductForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Extract single main image (may be nil — that's fine)
	mainImage, _ := c.FormFile("image")

	// Extract gallery images via the service's exported helper
	var galleryFiles = services.CollectGalleryFiles(c)

	companyID, _ := middlewares.GetCompanyID(c)
	product, err := service.CreateProduct(
		form,
		mainImage,
		galleryFiles,
		c.PostForm("attributes"),
		c.PostForm("variants"),
		companyID,
		c,
	)
	if err != nil {
		if err.Error() == "product name is required" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "a foreign key constraint fails") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid reference: vendor, category, or location ID does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"data":    product,
	})
}

// UpdateProduct builds the partial-update map from PostForm fields (only
// fields the client actually sends are included), extracts files, and
// delegates to the service.
func UpdateProduct(c *gin.Context) {
	service := services.NewProductService()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	updates := make(map[string]interface{})
	if v := c.PostForm("name"); v != "" {
		updates["name"] = v
	}
	if v := c.PostForm("description"); v != "" {
		updates["description"] = v
	}
	if v := c.PostForm("category_id"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 32); err == nil {
			updates["category_id"] = uint(parsed)
		}
	}
	if v := c.PostForm("vendor_id"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 32); err == nil {
			updates["vendor_id"] = uint(parsed)
		}
	}
	if v := c.PostForm("location_id"); v != "" {
		if parsed, err := strconv.ParseUint(v, 10, 32); err == nil {
			updates["location_id"] = uint(parsed)
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

	mainImage, _ := c.FormFile("image")
	galleryFiles := services.CollectGalleryFiles(c)

	companyID, _ := middlewares.GetCompanyID(c)
	product, err := service.UpdateProduct(
		uint(id),
		companyID,
		updates,
		mainImage,
		galleryFiles,
		c.PostForm("attributes"),
		c.PostForm("variants"),
		c,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		if strings.Contains(err.Error(), "a foreign key constraint fails") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid reference: vendor, category, or location ID does not exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"data":    product,
	})
}

// UpdateProductStatus toggles the published flag via JSON body { "published": bool }.
func UpdateProductStatus(c *gin.Context) {
	service := services.NewProductService()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var body struct {
		Published bool `json:"published"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)
	product, err := service.TogglePublished(uint(id), companyID, body.Published)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product status updated successfully",
		"data":    product,
	})
}

// DeleteProduct soft-deletes the product and cleans up image files.
func DeleteProduct(c *gin.Context) {
	service := services.NewProductService()
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	companyID, _ := middlewares.GetCompanyID(c)
	if err := service.DeleteProductByCompany(uint(id), companyID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
