package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/middlewares"
	"go-crud/models"
	"gorm.io/gorm"
)

var errInsufficientStock = errors.New("insufficient stock")

// GetProductsByLocation returns all products that have stock at a given warehouse/location.
// Uses the same preload logic as ListProducts.
// GET /transfers/products-by-location/:location_id
// Query params: search (name/sku), page, limit
func GetProductsByLocation(c *gin.Context) {
	locationID, err := strconv.ParseUint(c.Param("location_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid location_id"})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}
	offset := (page - 1) * limit

	// Products at this location = simple products with location_id match
	// OR products that have at least one variant with inventory at this location.
	// Variants are filtered to only those with stock at the requested location.
	query := database.DB.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants", func(db *gorm.DB) *gorm.DB {
			return db.Where(
				"id IN (SELECT variant_id FROM variant_inventory WHERE location_id = ? AND quantity > 0)",
				locationID,
			)
		}).
		Preload("Variants.Inventory", func(db *gorm.DB) *gorm.DB {
			return db.Where("location_id = ?", locationID)
		}).
		Preload("Variants.Inventory.Location").
		Preload("Images").
		Where(`
			products.deleted_at IS NULL AND products.company_id = ? AND (
				products.location_id = ?
				OR products.id IN (
					SELECT DISTINCT pv.product_id
					FROM product_variants pv
					JOIN variant_inventory vi ON vi.variant_id = pv.id AND vi.location_id = ? AND vi.quantity > 0
					WHERE pv.deleted_at IS NULL
				)
			)
		`, companyID, locationID, locationID)

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("products.name LIKE ? OR products.sku LIKE ?", like, like)
	}

	var total int64
	query.Model(&models.Product{}).Count(&total)

	var products []models.Product
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// Override each variant's stock with the location-specific quantity
	for i := range products {
		for j := range products[i].Variants {
			if len(products[i].Variants[j].Inventory) > 0 {
				products[i].Variants[j].Stock = products[i].Variants[j].Inventory[0].Quantity
			} else {
				products[i].Variants[j].Stock = 0
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Products at location retrieved successfully",
		"location_id": locationID,
		"data":        products,
		"total":       total,
		"page":        page,
		"limit":       limit,
	})
}

// ListTransfers returns stock transfer records with optional filters.
// Query params: status, product_id, from_location_id, to_location_id, page, limit
func ListTransfers(c *gin.Context) {
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	query := database.DB.
		Preload("Product").
		Preload("Variant").
		Preload("FromLocation").
		Preload("ToLocation").
		Model(&models.StockTransfer{}).
		Where("company_id = ?", companyID)

	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	if productID := c.Query("product_id"); productID != "" {
		query = query.Where("product_id = ?", productID)
	}
	if fromID := c.Query("from_location_id"); fromID != "" {
		query = query.Where("from_location_id = ?", fromID)
	}
	if toID := c.Query("to_location_id"); toID != "" {
		query = query.Where("to_location_id = ?", toID)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	var transfers []models.StockTransfer
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transfers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transfers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfers retrieved successfully",
		"data":    transfers,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// transferRequest is the body for POST /transfers
type transferRequest struct {
	ProductID      uint   `json:"productId" binding:"required"`
	VariantID      *uint  `json:"variantId"`
	FromLocationID uint   `json:"fromLocationId" binding:"required"`
	ToLocationID   uint   `json:"toLocationId" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	Notes          string `json:"notes"`
}

// CreateTransfer validates available stock, executes the move, and records the transfer.
func CreateTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	if req.FromLocationID == req.ToLocationID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source and destination warehouse must be different"})
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND company_id = ?", req.ProductID, companyID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var transfer models.StockTransfer
	transfer.CompanyID = companyID

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if req.VariantID != nil {
			return moveVariantStock(tx, req, &transfer)
		}
		return moveProductStock(tx, req, product, &transfer)
	})

	if err != nil {
		if errors.Is(err, errInsufficientStock) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock in source warehouse"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transfer", "details": err.Error()})
		return
	}

	database.DB.
		Preload("Product").
		Preload("Variant").
		Preload("FromLocation").
		Preload("ToLocation").
		First(&transfer, transfer.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Transfer completed successfully", "data": transfer})
}

// CancelTransfer reverses stock and marks the transfer as Cancelled.
func CancelTransfer(c *gin.Context) {
	var transfer models.StockTransfer
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&transfer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transfer not found"})
		return
	}

	if transfer.Status != "Completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only completed transfers can be cancelled"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if transfer.VariantID != nil {
			if err := reverseVariantStock(tx, transfer); err != nil {
				return err
			}
		} else {
			if err := reverseProductStock(tx, transfer); err != nil {
				return err
			}
		}
		return tx.Model(&transfer).Update("status", "Cancelled").Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel transfer", "details": err.Error()})
		return
	}

	database.DB.
		Preload("Product").
		Preload("Variant").
		Preload("FromLocation").
		Preload("ToLocation").
		First(&transfer, transfer.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Transfer cancelled successfully", "data": transfer})
}

// --- variant stock helpers ---

func moveVariantStock(tx *gorm.DB, req transferRequest, transfer *models.StockTransfer) error {
	var source models.VariantInventory
	err := tx.Where("variant_id = ? AND location_id = ?", *req.VariantID, req.FromLocationID).First(&source).Error
	if err != nil {
		// No variant_inventory row yet — fall back to variant.stock if the product
		// is at the requested from_location.
		var variant models.ProductVariant
		if err2 := tx.First(&variant, *req.VariantID); err2 != nil {
			return errInsufficientStock
		}
		var product models.Product
		if err2 := tx.First(&product, variant.ProductID); err2 != nil {
			return errInsufficientStock
		}
		if product.LocationID == nil || *product.LocationID != req.FromLocationID {
			return errInsufficientStock
		}
		if variant.Stock < req.Quantity {
			return errInsufficientStock
		}
		// Seed variant_inventory from variant.stock, then deduct
		source = models.VariantInventory{
			VariantID:  *req.VariantID,
			LocationID: req.FromLocationID,
			Quantity:   variant.Stock - req.Quantity,
		}
		if err2 := tx.Create(&source).Error; err2 != nil {
			return err2
		}
		// Also zero out the flat stock on the variant to keep it consistent
		if err2 := tx.Model(&variant).Update("stock", variant.Stock-req.Quantity).Error; err2 != nil {
			return err2
		}
	} else {
		if source.Quantity < req.Quantity {
			return errInsufficientStock
		}
		// deduct source
		if err := tx.Model(&source).Update("quantity", source.Quantity-req.Quantity).Error; err != nil {
			return err
		}
	}

	// add to destination (create row if missing)
	if err := upsertInventory(tx, *req.VariantID, req.ToLocationID, req.Quantity); err != nil {
		return err
	}

	// Sync product_variants.stock = sum of all variant_inventory quantities
	if err := syncVariantStock(tx, *req.VariantID); err != nil {
		return err
	}

	*transfer = models.StockTransfer{
		ProductID:      req.ProductID,
		VariantID:      req.VariantID,
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Quantity:       req.Quantity,
		Status:         "Completed",
		Notes:          req.Notes,
	}
	return tx.Create(transfer).Error
}

func reverseVariantStock(tx *gorm.DB, t models.StockTransfer) error {
	// deduct from destination
	var dest models.VariantInventory
	if err := tx.Where("variant_id = ? AND location_id = ?", *t.VariantID, t.ToLocationID).First(&dest).Error; err != nil {
		return err
	}
	if err := tx.Model(&dest).Update("quantity", dest.Quantity-t.Quantity).Error; err != nil {
		return err
	}

	// restore to source
	if err := upsertInventory(tx, *t.VariantID, t.FromLocationID, t.Quantity); err != nil {
		return err
	}

	return syncVariantStock(tx, *t.VariantID)
}

// syncVariantStock keeps product_variants.stock in sync with the sum of all
// variant_inventory quantities for that variant.
func syncVariantStock(tx *gorm.DB, variantID uint) error {
	var total int
	if err := tx.Model(&models.VariantInventory{}).
		Where("variant_id = ?", variantID).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&total).Error; err != nil {
		return err
	}
	return tx.Model(&models.ProductVariant{}).
		Where("id = ?", variantID).
		Update("stock", total).Error
}

// upsertInventory adds qty to an existing variant_inventory row or creates one.
func upsertInventory(tx *gorm.DB, variantID, locationID uint, qty int) error {
	var row models.VariantInventory
	if err := tx.Where("variant_id = ? AND location_id = ?", variantID, locationID).First(&row).Error; err != nil {
		return tx.Create(&models.VariantInventory{
			VariantID:  variantID,
			LocationID: locationID,
			Quantity:   qty,
		}).Error
	}
	return tx.Model(&row).Update("quantity", row.Quantity+qty).Error
}

// --- simple product stock helpers ---

func moveProductStock(tx *gorm.DB, req transferRequest, product models.Product, transfer *models.StockTransfer) error {
	if product.LocationID == nil || *product.LocationID != req.FromLocationID {
		return errInsufficientStock
	}
	if product.Stock < req.Quantity {
		return errInsufficientStock
	}

	newStock := product.Stock - req.Quantity
	updates := map[string]interface{}{"stock": newStock}
	if newStock == 0 {
		// all stock moved — point product to new location
		updates["location_id"] = req.ToLocationID
		updates["stock"] = req.Quantity
	}

	if err := tx.Model(&product).Updates(updates).Error; err != nil {
		return err
	}

	*transfer = models.StockTransfer{
		ProductID:      req.ProductID,
		FromLocationID: req.FromLocationID,
		ToLocationID:   req.ToLocationID,
		Quantity:       req.Quantity,
		Status:         "Completed",
		Notes:          req.Notes,
	}
	return tx.Create(transfer).Error
}

func reverseProductStock(tx *gorm.DB, t models.StockTransfer) error {
	var product models.Product
	if err := tx.First(&product, t.ProductID).Error; err != nil {
		return err
	}

	// restore stock at source, point location back to source
	return tx.Model(&product).Updates(map[string]interface{}{
		"stock":       product.Stock + t.Quantity,
		"location_id": t.FromLocationID,
	}).Error
}
