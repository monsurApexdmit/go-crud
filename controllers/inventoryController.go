package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-crud/database"
	"go-crud/models"
)

// ListInventory returns every product and its variants with per-warehouse stock.
// Query params: search (name/sku), location_id (filter by warehouse), page, limit
func ListInventory(c *gin.Context) {
	query := database.DB.
		Preload("Variants.Inventory.Location").
		Preload("Location").
		Model(&models.Product{})

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("products.name LIKE ? OR products.sku LIKE ?", like, like)
	}
	if locationID := c.Query("location_id"); locationID != "" {
		if id, err := strconv.Atoi(locationID); err == nil {
			query = query.Where("products.location_id = ?", id)
		}
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

	var products []models.Product
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve inventory"})
		return
	}

	rows := buildInventoryRows(products)

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory retrieved successfully",
		"data":    rows,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// inventoryRow is the flat structure the frontend renders — one row per product or variant.
type inventoryRow struct {
	Type        string                   `json:"type"`        // "product" | "variant"
	ID          uint                     `json:"id"`
	ProductID   uint                     `json:"productId"`
	ProductName string                   `json:"productName"`
	VariantName string                   `json:"variantName,omitempty"`
	SKU         string                   `json:"sku"`
	Barcode     string                   `json:"barcode"`
	Stock       int                      `json:"stock"`       // total across all warehouses
	Inventory   []inventoryWarehouseQty  `json:"inventory"`   // per-warehouse breakdown
}

type inventoryWarehouseQty struct {
	LocationID   uint   `json:"locationId"`
	LocationName string `json:"locationName"`
	Quantity     int    `json:"quantity"`
}

// buildInventoryRows flattens products+variants into the table rows the frontend expects.
func buildInventoryRows(products []models.Product) []inventoryRow {
	var rows []inventoryRow

	for _, p := range products {
		if len(p.Variants) == 0 {
			// Simple product — stock lives on the product row at its location_id
			inv := []inventoryWarehouseQty{}
			if p.LocationID != nil && p.Stock > 0 {
				inv = append(inv, inventoryWarehouseQty{
					LocationID:   *p.LocationID,
					LocationName: p.Location.Name,
					Quantity:     p.Stock,
				})
			}
			rows = append(rows, inventoryRow{
				Type:        "product",
				ID:          p.ID,
				ProductID:   p.ID,
				ProductName: p.Name,
				SKU:         p.SKU,
				Barcode:     p.Barcode,
				Stock:       p.Stock,
				Inventory:   inv,
			})
			continue
		}

		// Product with variants — each variant is its own row
		for _, v := range p.Variants {
			inv := make([]inventoryWarehouseQty, 0, len(v.Inventory))
			totalStock := 0
			for _, vi := range v.Inventory {
				inv = append(inv, inventoryWarehouseQty{
					LocationID:   vi.LocationID,
					LocationName: vi.Location.Name,
					Quantity:     vi.Quantity,
				})
				totalStock += vi.Quantity
			}
			rows = append(rows, inventoryRow{
				Type:        "variant",
				ID:          v.ID,
				ProductID:   p.ID,
				ProductName: p.Name,
				VariantName: v.Name,
				SKU:         v.SKU,
				Barcode:     v.Barcode,
				Stock:       totalStock,
				Inventory:   inv,
			})
		}
	}

	return rows
}
