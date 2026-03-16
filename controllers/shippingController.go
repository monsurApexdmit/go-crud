package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/middlewares"
	"go-crud/models"
	"gorm.io/gorm"
)

// ============ Shipping Addresses ============

// ListShippingAddresses retrieves shipping addresses with filtering
func ListShippingAddresses(c *gin.Context) {
	var addresses []models.ShippingAddress
	query := database.DB.Model(&models.ShippingAddress{})

	// Filter by company_id
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	query = query.Where("company_id = ?", companyID)

	// Filter by customer ID
	if customerID := c.Query("customer_id"); customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	// Filter by address type
	if addressType := c.Query("address_type"); addressType != "" {
		query = query.Where("address_type = ?", addressType)
	}

	// Filter by default addresses only
	if isDefault := c.Query("is_default"); isDefault == "true" {
		query = query.Where("is_default = ?", true)
	}

	if err := query.Preload("Customer").Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shipping addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipping addresses retrieved successfully", "data": addresses})
}

// GetShippingAddress retrieves a single shipping address
func GetShippingAddress(c *gin.Context) {
	var address models.ShippingAddress
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Preload("Customer").
		Where("id = ? AND company_id = ?", c.Param("id"), companyID).
		First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shipping address fetched successfully", "data": address})
}

// CreateShippingAddress creates a new shipping address
func CreateShippingAddress(c *gin.Context) {
	var address models.ShippingAddress
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	address.CompanyID = companyID

	// Validate required fields
	if address.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Full name is required"})
		return
	}
	if address.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is required"})
		return
	}
	if address.AddressLine1 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address line 1 is required"})
		return
	}
	if address.City == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "City is required"})
		return
	}
	if address.State == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State is required"})
		return
	}
	if address.PostalCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Postal code is required"})
		return
	}

	// If this address is set as default, unset other default addresses for this customer
	if address.IsDefault && address.CustomerID != nil {
		database.DB.Model(&models.ShippingAddress{}).
			Where("customer_id = ? AND is_default = ?", *address.CustomerID, true).
			Update("is_default", false)
	}

	if err := database.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipping address"})
		return
	}

	database.DB.Preload("Customer").First(&address, address.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Shipping address created successfully", "data": address})
}

// UpdateShippingAddress updates an existing shipping address
func UpdateShippingAddress(c *gin.Context) {
	var address models.ShippingAddress
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}

	var input models.ShippingAddress
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// If updating to default, unset other default addresses
	if input.IsDefault && address.CustomerID != nil {
		database.DB.Model(&models.ShippingAddress{}).
			Where("customer_id = ? AND id != ? AND is_default = ?", *address.CustomerID, address.ID, true).
			Update("is_default", false)
	}

	// Update fields
	if input.FullName != "" {
		address.FullName = input.FullName
	}
	if input.Phone != "" {
		address.Phone = input.Phone
	}
	if input.Email != "" {
		address.Email = input.Email
	}
	if input.AddressLine1 != "" {
		address.AddressLine1 = input.AddressLine1
	}
	if input.AddressLine2 != "" {
		address.AddressLine2 = input.AddressLine2
	}
	if input.City != "" {
		address.City = input.City
	}
	if input.State != "" {
		address.State = input.State
	}
	if input.PostalCode != "" {
		address.PostalCode = input.PostalCode
	}
	if input.Country != "" {
		address.Country = input.Country
	}
	if input.AddressType != "" {
		address.AddressType = input.AddressType
	}
	address.IsDefault = input.IsDefault

	if err := database.DB.Save(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipping address"})
		return
	}

	database.DB.Preload("Customer").First(&address, address.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Shipping address updated successfully", "data": address})
}

// DeleteShippingAddress soft-deletes a shipping address
func DeleteShippingAddress(c *gin.Context) {
	var address models.ShippingAddress
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}

	database.DB.Delete(&address)
	c.JSON(http.StatusOK, gin.H{"message": "Shipping address deleted successfully"})
}

// SetDefaultShippingAddress sets an address as default for a customer
func SetDefaultShippingAddress(c *gin.Context) {
	var address models.ShippingAddress
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&address).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping address not found"})
		return
	}

	if address.CustomerID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot set default for address without customer"})
		return
	}

	// Unset other default addresses
	database.DB.Model(&models.ShippingAddress{}).
		Where("customer_id = ? AND is_default = ?", *address.CustomerID, true).
		Update("is_default", false)

	// Set this as default
	address.IsDefault = true
	database.DB.Save(&address)

	c.JSON(http.StatusOK, gin.H{"message": "Default shipping address updated successfully", "data": address})
}

// ============ Order Shipments ============

// ListOrderShipments retrieves order shipments with filtering
func ListOrderShipments(c *gin.Context) {
	var shipments []models.OrderShipment
	query := database.DB.Model(&models.OrderShipment{})

	// Filter by company_id
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	query = query.Where("company_id = ?", companyID)

	// Filter by sell/order ID
	if sellID := c.Query("sell_id"); sellID != "" {
		query = query.Where("sell_id = ?", sellID)
	}

	// Filter by status
	if status := c.Query("status"); status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Search by tracking number
	if trackingNumber := c.Query("tracking_number"); trackingNumber != "" {
		query = query.Where("tracking_number LIKE ?", "%"+trackingNumber+"%")
	}

	// Filter by carrier
	if carrier := c.Query("carrier"); carrier != "" {
		query = query.Where("carrier = ?", carrier)
	}

	// Order by most recent first
	query = query.Order("created_at DESC")

	// Pagination
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
		Preload("Sell").
		Preload("TrackingHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("event_time DESC")
		}).
		Find(&shipments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve shipments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Shipments retrieved successfully",
		"data":    shipments,
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
	})
}

// GetOrderShipment retrieves a single shipment with full tracking history
func GetOrderShipment(c *gin.Context) {
	var shipment models.OrderShipment
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.
		Preload("Sell").
		Preload("Sell.Customer").
		Preload("Sell.ShippingAddress").
		Preload("TrackingHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("event_time DESC")
		}).
		Where("id = ? AND company_id = ?", c.Param("id"), companyID).
		First(&shipment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shipment fetched successfully", "data": shipment})
}

// CreateOrderShipment creates a new shipment for an order
func CreateOrderShipment(c *gin.Context) {
	var shipment models.OrderShipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	shipment.CompanyID = companyID

	// Validate required fields
	if shipment.SellID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID (sellId) is required"})
		return
	}
	if shipment.TrackingNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tracking number is required"})
		return
	}
	if shipment.Carrier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Carrier is required"})
		return
	}

	// Validate order exists
	var sell models.Sell
	if err := database.DB.First(&sell, shipment.SellID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "picked_up", "in_transit", "out_for_delivery", "delivered", "failed", "returned"}
	if shipment.Status != "" && !contains(validStatuses, shipment.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shipment status"})
		return
	}
	if shipment.Status == "" {
		shipment.Status = "pending"
	}

	// Create shipment and update order in transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Create shipment
		if err := tx.Create(&shipment).Error; err != nil {
			return err
		}

		// Update sell record with tracking info
		updates := map[string]interface{}{
			"tracking_number":     shipment.TrackingNumber,
			"carrier":             shipment.Carrier,
			"fulfillment_status": "processing",
		}
		if shipment.ShippedAt != nil {
			updates["shipped_at"] = shipment.ShippedAt
			updates["fulfillment_status"] = "shipped"
		}
		if err := tx.Model(&sell).Updates(updates).Error; err != nil {
			return err
		}

		// Add initial tracking event
		tracking := models.ShipmentTrackingHistory{
			ShipmentID:  shipment.ID,
			Status:      shipment.Status,
			Description: fmt.Sprintf("Shipment created with %s", shipment.Carrier),
			EventTime:   time.Now(),
		}
		if err := tx.Create(&tracking).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shipment"})
		return
	}

	database.DB.Preload("Sell").Preload("TrackingHistory").First(&shipment, shipment.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Shipment created successfully", "data": shipment})
}

// UpdateShipmentStatus updates the status of a shipment
func UpdateShipmentStatus(c *gin.Context) {
	var shipment models.OrderShipment
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&shipment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
		return
	}

	var input struct {
		Status      string `json:"status" binding:"required"`
		Location    string `json:"location"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	validStatuses := []string{"pending", "picked_up", "in_transit", "out_for_delivery", "delivered", "failed", "returned"}
	if !contains(validStatuses, input.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Update shipment and order in transaction
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		oldStatus := shipment.Status
		shipment.Status = input.Status

		// Set timestamps based on status
		now := time.Now()
		if input.Status == "picked_up" && shipment.ShippedAt == nil {
			shipment.ShippedAt = &now
		}
		if input.Status == "delivered" && shipment.DeliveredAt == nil {
			shipment.DeliveredAt = &now
		}

		if err := tx.Save(&shipment).Error; err != nil {
			return err
		}

		// Update order fulfillment status
		var sell models.Sell
		if err := tx.First(&sell, shipment.SellID).Error; err != nil {
			return err
		}

		fulfillmentStatus := mapShipmentToFulfillmentStatus(input.Status)
		updates := map[string]interface{}{"fulfillment_status": fulfillmentStatus}
		if input.Status == "picked_up" || input.Status == "in_transit" {
			updates["shipped_at"] = now
		}
		if input.Status == "delivered" {
			updates["delivered_at"] = now
		}
		if err := tx.Model(&sell).Updates(updates).Error; err != nil {
			return err
		}

		// Add tracking history
		description := input.Description
		if description == "" {
			description = fmt.Sprintf("Status changed from %s to %s", oldStatus, input.Status)
		}
		tracking := models.ShipmentTrackingHistory{
			ShipmentID:  shipment.ID,
			Status:      input.Status,
			Location:    input.Location,
			Description: description,
			EventTime:   now,
		}
		if err := tx.Create(&tracking).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update shipment status"})
		return
	}

	database.DB.Preload("Sell").Preload("TrackingHistory", func(db *gorm.DB) *gorm.DB {
		return db.Order("event_time DESC")
	}).First(&shipment, shipment.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Shipment status updated successfully", "data": shipment})
}

// AddTrackingEvent adds a tracking event to shipment history
func AddTrackingEvent(c *gin.Context) {
	var shipment models.OrderShipment
	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}
	if err := database.DB.Where("id = ? AND company_id = ?", c.Param("id"), companyID).First(&shipment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
		return
	}

	var input models.ShipmentTrackingHistory
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	input.ShipmentID = shipment.ID
	if input.EventTime.IsZero() {
		input.EventTime = time.Now()
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add tracking event"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tracking event added successfully", "data": input})
}

// TrackShipment retrieves shipment info by tracking number (public endpoint)
func TrackShipment(c *gin.Context) {
	trackingNumber := c.Param("trackingNumber")

	var shipment models.OrderShipment
	if err := database.DB.
		Preload("TrackingHistory", func(db *gorm.DB) *gorm.DB {
			return db.Order("event_time DESC")
		}).
		Where("tracking_number = ?", trackingNumber).
		First(&shipment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipment not found"})
		return
	}

	// Return limited info for public tracking
	response := gin.H{
		"trackingNumber":    shipment.TrackingNumber,
		"carrier":           shipment.Carrier,
		"status":            shipment.Status,
		"shippedAt":         shipment.ShippedAt,
		"estimatedDelivery": shipment.EstimatedDelivery,
		"deliveredAt":       shipment.DeliveredAt,
		"trackingHistory":   shipment.TrackingHistory,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment tracking retrieved successfully", "data": response})
}

// GetShipmentStats returns statistics for shipments
func GetShipmentStats(c *gin.Context) {
	var stats struct {
		Total          int64   `json:"total"`
		Pending        int64   `json:"pending"`
		InTransit      int64   `json:"inTransit"`
		Delivered      int64   `json:"delivered"`
		Failed         int64   `json:"failed"`
		TotalShipCost  float64 `json:"totalShippingCost"`
		AvgDeliveryDays float64 `json:"avgDeliveryDays"`
	}

	companyID, ok := middlewares.GetCompanyID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in context"})
		return
	}

	database.DB.Model(&models.OrderShipment{}).Where("company_id = ?", companyID).Count(&stats.Total)
	database.DB.Model(&models.OrderShipment{}).Where("status = ? AND company_id = ?", "pending", companyID).Count(&stats.Pending)
	database.DB.Model(&models.OrderShipment{}).Where("status = ? AND company_id = ?", "in_transit", companyID).Count(&stats.InTransit)
	database.DB.Model(&models.OrderShipment{}).Where("status = ? AND company_id = ?", "delivered", companyID).Count(&stats.Delivered)
	database.DB.Model(&models.OrderShipment{}).Where("status = ? AND company_id = ?", "failed", companyID).Count(&stats.Failed)
	database.DB.Model(&models.OrderShipment{}).
		Where("company_id = ?", companyID).
		Select("COALESCE(SUM(shipping_cost), 0)").Scan(&stats.TotalShipCost)

	// Calculate average delivery time (in days) for delivered shipments
	var avgDays float64
	database.DB.Model(&models.OrderShipment{}).
		Where("status = ? AND company_id = ? AND shipped_at IS NOT NULL AND delivered_at IS NOT NULL", "delivered", companyID).
		Select("AVG(TIMESTAMPDIFF(DAY, shipped_at, delivered_at))").
		Scan(&avgDays)
	stats.AvgDeliveryDays = avgDays

	c.JSON(http.StatusOK, gin.H{"message": "Shipment statistics retrieved successfully", "data": stats})
}

// Helper function to map shipment status to fulfillment status
func mapShipmentToFulfillmentStatus(shipmentStatus string) string {
	switch shipmentStatus {
	case "pending":
		return "processing"
	case "picked_up", "in_transit", "out_for_delivery":
		return "shipped"
	case "delivered":
		return "delivered"
	case "failed", "returned":
		return "cancelled"
	default:
		return "unfulfilled"
	}
}
