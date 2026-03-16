package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/middlewares"
	"go-crud/models"
	"gorm.io/gorm"
)

// ListSells retrieves sells with filtering, search, and pagination
func ListSells(c *gin.Context) {
	var sells []models.Sell
	query := database.DB.Model(&models.Sell{})

	// Filter by company_id
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	query = query.Where("company_id = ?", companyID)

	// Search by customer name
	if search := c.Query("search"); search != "" {
		query = query.Where("customer_name LIKE ?", "%"+search+"%")
	}

	// Filter by status
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Filter by payment method
	if method := c.Query("method"); method != "" && method != "all" {
		query = query.Where("method = ?", method)
	}

	// Filter by customer ID
	if customerID := c.Query("customer_id"); customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("order_time >= ?", t)
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add 23:59:59 to include the entire end date
			endOfDay := t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			query = query.Where("order_time <= ?", endOfDay)
		}
	}

	// Order by most recent first
	query = query.Order("order_time DESC")

	// Limit (for "Last 5", "Last 10", etc.) — if set, skip pagination
	limitParam := c.Query("limit")
	if limitParam != "" && limitParam != "all" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			if err := query.Limit(l).
				Preload("Customer").
				Preload("ShippingAddress").
				Preload("Items").
				Preload("Shipments").
				Find(&sells).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sells"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Sells retrieved successfully",
				"data":    sells,
			})
			return
		}
	}

	// Pagination (only if limit is not set)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * perPage
	if err := query.Offset(offset).Limit(perPage).
		Preload("Customer").
		Preload("ShippingAddress").
		Preload("Items").
		Preload("Shipments").
		Find(&sells).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sells"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sells retrieved successfully",
		"data":    sells,
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

// GetSell retrieves a single sell by ID
func GetSell(c *gin.Context) {
	var sell models.Sell
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.
		Preload("Customer").
		Preload("ShippingAddress").
		Preload("Items").
		Preload("Shipments").
		Preload("Shipments.TrackingHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("event_time DESC")
		}).
		Where("id = ? AND company_id = ?", c.Param("id"), companyID).
		First(&sell).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sell not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sell fetched successfully", "data": sell})
}

// GetSellByInvoice retrieves a single sell by invoice number
func GetSellByInvoice(c *gin.Context) {
	var sell models.Sell
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.
		Preload("Customer").
		Preload("ShippingAddress").
		Preload("Items").
		Preload("Shipments").
		Preload("Shipments.TrackingHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("event_time DESC")
		}).
		Where("invoice_no = ? AND company_id = ?", c.Param("invoiceNo"), companyID).
		First(&sell).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sell not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sell fetched successfully", "data": sell})
}

// sellItemInput accepts both "unitPrice" and "price" from the frontend
type sellItemInput struct {
	ProductID        *uint   `json:"productId"`
	ProductIDSnake   *uint   `json:"product_id"`
	VariantID        *uint   `json:"variantId"`
	VariantIDSnake   *uint   `json:"variant_id"`
	InventoryID      *uint   `json:"inventoryId"`
	InventoryIDSnake *uint   `json:"inventory_id"`
	ProductName      string  `json:"productName"`
	VariantName      string  `json:"variantName"`
	Quantity         int     `json:"quantity"`
	UnitPrice        float64 `json:"unitPrice"`
	Price            float64 `json:"price"` // alias accepted from frontend
}

type createSellInput struct {
	InvoiceNo            string          `json:"invoiceNo"`
	OrderTime            *time.Time      `json:"orderTime"`
	CustomerID           *uint           `json:"customerId"`
	CustomerName         string          `json:"customerName"`
	ShippingAddressID    *uint           `json:"shippingAddressId"`
	ShippingFullName     string          `json:"shippingFullName"`
	ShippingPhone        string          `json:"shippingPhone"`
	ShippingEmail        string          `json:"shippingEmail"`
	ShippingAddressLine1 string          `json:"shippingAddressLine1"`
	ShippingAddressLine2 string          `json:"shippingAddressLine2"`
	ShippingCity         string          `json:"shippingCity"`
	ShippingState        string          `json:"shippingState"`
	ShippingPostalCode   string          `json:"shippingPostalCode"`
	ShippingCountry      string          `json:"shippingCountry"`
	ShippingAddressType  string          `json:"shippingAddressType"`
	Method               string          `json:"method"`
	Amount               float64         `json:"amount"`
	ShippingCost         float64         `json:"shippingCost"`
	ShippingMethod       string          `json:"shippingMethod"`
	CouponID             *uint           `json:"couponId"`
	CouponCode           string          `json:"couponCode"`
	Discount             float64         `json:"discount"`
	Status               string          `json:"status"`
	Notes                string          `json:"notes"`
	Items                []sellItemInput `json:"items"`
}

// CreateSell creates a new sell record
func CreateSell(c *gin.Context) {
	var input createSellInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	// Map input to sell model
	sell := models.Sell{
		CompanyID:            companyID,
		InvoiceNo:            input.InvoiceNo,
		CustomerID:           input.CustomerID,
		CustomerName:         input.CustomerName,
		ShippingAddressID:    input.ShippingAddressID,
		ShippingFullName:     input.ShippingFullName,
		ShippingPhone:        input.ShippingPhone,
		ShippingEmail:        input.ShippingEmail,
		ShippingAddressLine1: input.ShippingAddressLine1,
		ShippingAddressLine2: input.ShippingAddressLine2,
		ShippingCity:         input.ShippingCity,
		ShippingState:        input.ShippingState,
		ShippingPostalCode:   input.ShippingPostalCode,
		ShippingCountry:      input.ShippingCountry,
		ShippingAddressType:  input.ShippingAddressType,
		Method:               input.Method,
		Amount:               input.Amount,
		ShippingCost:         input.ShippingCost,
		ShippingMethod:       input.ShippingMethod,
		CouponID:             input.CouponID,
		CouponCode:           input.CouponCode,
		Discount:             input.Discount,
		Status:               input.Status,
		Notes:                input.Notes,
	}
	if input.OrderTime != nil {
		sell.OrderTime = *input.OrderTime
	}

	// Map items — resolve unitPrice from either unitPrice or price field
	sell.Items = make([]models.OrderItem, len(input.Items))
	for i, it := range input.Items {
		productID := it.ProductID
		if productID == nil {
			productID = it.ProductIDSnake
		}
		variantID := it.VariantID
		if variantID == nil {
			variantID = it.VariantIDSnake
		}
		inventoryID := it.InventoryID
		if inventoryID == nil {
			inventoryID = it.InventoryIDSnake
		}
		unitPrice := it.UnitPrice
		if unitPrice == 0 && it.Price > 0 {
			unitPrice = it.Price
		}
		qty := it.Quantity
		if qty == 0 {
			qty = 1
		}
		sell.Items[i] = models.OrderItem{
			ProductID:   productID,
			VariantID:   variantID,
			InventoryID: inventoryID,
			ProductName: it.ProductName,
			VariantName: it.VariantName,
			Quantity:    qty,
			UnitPrice:   unitPrice,
			TotalPrice:  unitPrice * float64(qty),
		}
	}

	// Generate invoice number if not provided
	if sell.InvoiceNo == "" {
		sell.InvoiceNo = generateInvoiceNo()
	}

	// Set order time to now if not provided
	if sell.OrderTime.IsZero() {
		sell.OrderTime = time.Now()
	}

	// Validate required fields
	if sell.CustomerName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer name is required"})
		return
	}

	// Validate status
	validStatuses := []string{"Pending", "Processing", "Delivered"}
	if sell.Status != "" && !contains(validStatuses, sell.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: Pending, Processing, Delivered"})
		return
	}
	if sell.Status == "" {
		sell.Status = "Pending"
	}

	// Validate method
	validMethods := []string{"Cash", "Card", "Online"}
	if sell.Method != "" && !contains(validMethods, sell.Method) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method. Must be one of: Cash, Card, Online"})
		return
	}
	if sell.Method == "" {
		sell.Method = "Cash"
	}

	// Validate and set shipping address
	var addressToUse *models.ShippingAddress

	// Check if custom shipping address fields are provided directly in payload
	hasCustomAddress := sell.ShippingFullName != "" || sell.ShippingAddressLine1 != ""

	if hasCustomAddress {
		// Option 1: Custom inline address provided (no need for shippingAddressId)
		// Address fields already in sell object, just validate required fields
		if sell.ShippingFullName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping full name is required when providing custom address"})
			return
		}
		if sell.ShippingAddressLine1 == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping address line 1 is required when providing custom address"})
			return
		}
		if sell.ShippingCity == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping city is required when providing custom address"})
			return
		}
		if sell.ShippingCountry == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping country is required when providing custom address"})
			return
		}
		// Set default address type if not provided (enum field requires valid value)
		if sell.ShippingAddressType == "" {
			sell.ShippingAddressType = "other"
		}
		// Custom address already in sell object, no need to copy
	} else if sell.ShippingAddressID != nil {
		// Option 2: Use saved shipping address by ID
		var address models.ShippingAddress
		if err := database.DB.First(&address, *sell.ShippingAddressID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipping address ID"})
			return
		}
		// Optionally validate it belongs to the customer
		if sell.CustomerID != nil && address.CustomerID != nil && *address.CustomerID != *sell.CustomerID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping address does not belong to this customer"})
			return
		}
		addressToUse = &address
	} else if sell.CustomerID != nil {
		// Option 3: Auto-select customer's default address
		var defaultAddress models.ShippingAddress
		err := database.DB.Where("customer_id = ? AND is_default = ?", *sell.CustomerID, true).First(&defaultAddress).Error
		if err == nil {
			// Default address found
			sell.ShippingAddressID = &defaultAddress.ID
			addressToUse = &defaultAddress
		}
		// If no default found, continue without address (guest checkout or pickup)
	}

	// Copy shipping address fields to order (snapshot for historical data)
	if addressToUse != nil {
		sell.ShippingFullName = addressToUse.FullName
		sell.ShippingPhone = addressToUse.Phone
		sell.ShippingEmail = addressToUse.Email
		sell.ShippingAddressLine1 = addressToUse.AddressLine1
		sell.ShippingAddressLine2 = addressToUse.AddressLine2
		sell.ShippingCity = addressToUse.City
		sell.ShippingState = addressToUse.State
		sell.ShippingPostalCode = addressToUse.PostalCode
		sell.ShippingCountry = addressToUse.Country
		sell.ShippingAddressType = addressToUse.AddressType
	}

	// Check for duplicate invoice number within company
	var existingCount int64
	database.DB.Model(&models.Sell{}).Where("invoice_no = ? AND company_id = ?", sell.InvoiceNo, companyID).Count(&existingCount)
	if existingCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Invoice number already exists"})
		return
	}

	// Create sell and items in a transaction; deduct stock immediately for active orders
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&sell).Error; err != nil {
			return err
		}

		if err := deductOrderItemsStock(tx, sell.Items); err != nil {
			return err
		}
		if err := tx.Model(&sell).Update("stock_deducted", true).Error; err != nil {
			return err
		}
		sell.StockDeducted = true

		return nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "insufficient stock") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sell"})
		return
	}

	// Reload with customer data, shipping address, and items
	database.DB.Preload("Customer").Preload("ShippingAddress").Preload("Items").First(&sell, sell.ID)

	c.JSON(http.StatusCreated, gin.H{"message": "Sell created successfully", "data": sell})
}

// UpdateSell updates an existing sell record
func UpdateSell(c *gin.Context) {
	var sell models.Sell
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&sell).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sell not found"})
		return
	}

	var input models.Sell
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update allowed fields
	if input.InvoiceNo != "" {
		sell.InvoiceNo = input.InvoiceNo
	}
	if !input.OrderTime.IsZero() {
		sell.OrderTime = input.OrderTime
	}
	if input.CustomerID != nil {
		sell.CustomerID = input.CustomerID
	}
	if input.CustomerName != "" {
		sell.CustomerName = input.CustomerName
	}
	if input.Method != "" {
		validMethods := []string{"Cash", "Card", "Online"}
		if !contains(validMethods, input.Method) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
			return
		}
		sell.Method = input.Method
	}
	if input.Amount > 0 {
		sell.Amount = input.Amount
	}
	if input.ShippingCost >= 0 {
		sell.ShippingCost = input.ShippingCost
	}
	if input.Discount >= 0 {
		sell.Discount = input.Discount
	}
	if input.Status != "" {
		validStatuses := []string{"Pending", "Processing", "Delivered"}
		if !contains(validStatuses, input.Status) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
			return
		}
		sell.Status = input.Status
	}
	if input.Notes != "" {
		sell.Notes = input.Notes
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Save(&sell).Error
	}); err != nil {
		if strings.Contains(err.Error(), "insufficient stock") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sell"})
		return
	}

	database.DB.Preload("Customer").Preload("ShippingAddress").Preload("Items").First(&sell, sell.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Sell updated successfully", "data": sell})
}

// UpdateSellStatus updates only the status of a sell
func UpdateSellStatus(c *gin.Context) {
	var sell models.Sell
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&sell).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sell not found"})
		return
	}

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	validStatuses := []string{"Pending", "Processing", "Delivered"}
	if !contains(validStatuses, input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status. Must be one of: Pending, Processing, Delivered"})
		return
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		sell.Status = input.Status
		return tx.Save(&sell).Error
	}); err != nil {
		if strings.Contains(err.Error(), "insufficient stock") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	database.DB.Preload("Customer").Preload("ShippingAddress").Preload("Items").First(&sell, sell.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully", "data": sell})
}

// DeleteSell soft-deletes a sell record
func DeleteSell(c *gin.Context) {
	var sell models.Sell
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).Preload("Items").First(&sell).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sell not found"})
		return
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if sell.StockDeducted {
			if err := restoreOrderItemsStock(tx, sell.Items); err != nil {
				return err
			}
			sell.StockDeducted = false
		}
		if err := tx.Save(&sell).Error; err != nil {
			return err
		}
		return tx.Delete(&sell).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sell and restore stock"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sell deleted successfully"})
}

func restoreOrderItemsStock(tx *gorm.DB, items []models.OrderItem) error {
	for _, item := range items {
		qty := item.Quantity
		if qty <= 0 {
			continue
		}

		if item.VariantID != nil {
			if item.InventoryID != nil {
				result := tx.Model(&models.VariantInventory{}).
					Where("id = ? AND variant_id = ?", *item.InventoryID, *item.VariantID).
					UpdateColumn("quantity", gorm.Expr("quantity + ?", qty))
				if result.Error != nil {
					return result.Error
				}
			} else {
				var inv models.VariantInventory
				if err := tx.Where("variant_id = ?", *item.VariantID).First(&inv).Error; err == nil {
					if err := tx.Model(&inv).UpdateColumn("quantity", gorm.Expr("quantity + ?", qty)).Error; err != nil {
						return err
					}
				}
			}

			if err := syncVariantStock(tx, *item.VariantID); err != nil {
				return err
			}
			
			var variant models.ProductVariant
			if err := tx.First(&variant, *item.VariantID).Error; err != nil {
				return fmt.Errorf("variant %d not found", *item.VariantID)
			}
			if err := syncProductStockFromVariants(tx, variant.ProductID); err != nil {
				return err
			}
			continue
		}

		if item.ProductID == nil {
			continue
		}
		if err := tx.Model(&models.Product{}).
			Where("id = ?", *item.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", qty)).Error; err != nil {
			return err
		}
	}
	return nil
}

func deductOrderItemsStock(tx *gorm.DB, items []models.OrderItem) error {
	for _, item := range items {
		qty := item.Quantity
		if qty <= 0 {
			continue
		}

		if item.VariantID != nil {
			var variant models.ProductVariant
			if err := tx.First(&variant, *item.VariantID).Error; err != nil {
				return fmt.Errorf("variant %d not found", *item.VariantID)
			}
			if variant.Stock < qty {
				return fmt.Errorf("insufficient stock for variant '%s' (available: %d, requested: %d)", variant.Name, variant.Stock, qty)
			}

			if item.InventoryID != nil {
				result := tx.Model(&models.VariantInventory{}).
					Where("id = ? AND variant_id = ? AND quantity >= ?", *item.InventoryID, *item.VariantID, qty).
					UpdateColumn("quantity", gorm.Expr("quantity - ?", qty))
				if result.Error != nil {
					return result.Error
				}
				if result.RowsAffected == 0 {
					return fmt.Errorf("insufficient stock for inventory %d", *item.InventoryID)
				}
			} else {
				remaining := qty
				var invRows []models.VariantInventory
				if err := tx.Where("variant_id = ? AND quantity > 0", *item.VariantID).
					Order("quantity DESC").
					Find(&invRows).Error; err != nil {
					return err
				}
				for i := range invRows {
					if remaining <= 0 {
						break
					}
					deduct := remaining
					if invRows[i].Quantity < deduct {
						deduct = invRows[i].Quantity
					}
					if err := tx.Model(&invRows[i]).UpdateColumn("quantity", invRows[i].Quantity-deduct).Error; err != nil {
						return err
					}
					remaining -= deduct
				}
				if remaining > 0 {
					return fmt.Errorf("insufficient stock for variant '%s' (available: %d, requested: %d)", variant.Name, variant.Stock, qty)
				}
			}

			if err := syncVariantStock(tx, *item.VariantID); err != nil {
				return err
			}
			if err := syncProductStockFromVariants(tx, variant.ProductID); err != nil {
				return err
			}
			continue
		}

		if item.ProductID == nil {
			continue
		}
		result := tx.Model(&models.Product{}).
			Where("id = ? AND stock >= ?", *item.ProductID, qty).
			UpdateColumn("stock", gorm.Expr("stock - ?", qty))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("insufficient stock for product %d", *item.ProductID)
		}
	}
	return nil
}

func syncProductStockFromVariants(tx *gorm.DB, productID uint) error {
	var total int
	if err := tx.Model(&models.ProductVariant{}).
		Where("product_id = ?", productID).
		Select("COALESCE(SUM(stock), 0)").
		Scan(&total).Error; err != nil {
		return err
	}
	return tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock", total).Error
}

// generateInvoiceNo generates a unique invoice number
func generateInvoiceNo() string {
	return fmt.Sprintf("INV-%d", time.Now().Unix())
}

// GetSellsStats returns summary statistics for sells/orders
func GetSellsStats(c *gin.Context) {
	var stats struct {
		TotalSells       int64   `json:"totalSells"`
		TotalRevenue     float64 `json:"totalRevenue"`
		PendingOrders    int64   `json:"pendingOrders"`
		ProcessingOrders int64   `json:"processingOrders"`
		DeliveredOrders  int64   `json:"deliveredOrders"`
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	// Total sells
	database.DB.Model(&models.Sell{}).Where("company_id = ?", companyID).Count(&stats.TotalSells)

	// Total revenue
	database.DB.Model(&models.Sell{}).Where("company_id = ?", companyID).Select("COALESCE(SUM(amount), 0)").Scan(&stats.TotalRevenue)

	// Status counts
	database.DB.Model(&models.Sell{}).Where("status = ? AND company_id = ?", "Pending", companyID).Count(&stats.PendingOrders)
	database.DB.Model(&models.Sell{}).Where("status = ? AND company_id = ?", "Processing", companyID).Count(&stats.ProcessingOrders)
	database.DB.Model(&models.Sell{}).Where("status = ? AND company_id = ?", "Delivered", companyID).Count(&stats.DeliveredOrders)

	c.JSON(http.StatusOK, gin.H{"message": "Statistics retrieved successfully", "data": stats})
}

// contains checks if a string slice contains a value
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
