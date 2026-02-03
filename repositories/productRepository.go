package repositories

import (
	"go-crud/database"
	"go-crud/models"

	"gorm.io/gorm"
)

// ProductFilters carries the optional WHERE clauses for the list query.
// Nil pointers are skipped.
type ProductFilters struct {
	CategoryID *uint
	VendorID   *uint
	LocationID *uint
}

// ProductRepository is the single owner of every GORM call that touches
// products, product_variants, product_images, and product_attributes.
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() ProductRepository {
	return ProductRepository{db: database.DB}
}

// --- Read ---

// FindAll returns the paginated product list and the total count before limit/offset.
func (r ProductRepository) FindAll(filters ProductFilters, page, limit, offset int) ([]models.Product, int64, error) {
	var products []models.Product

	query := r.db.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		Preload("Images")

	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	}
	if filters.VendorID != nil {
		query = query.Where("vendor_id = ?", *filters.VendorID)
	}
	if filters.LocationID != nil {
		query = query.Where("location_id = ?", *filters.LocationID)
	}

	var total int64
	if err := query.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindByID loads a single product with the full preload chain.
func (r ProductRepository) FindByID(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.
		Preload("Category").
		Preload("Vendor").
		Preload("Location").
		Preload("Attributes").
		Preload("Variants.Inventory.Location").
		Preload("Images").
		First(&product, id).Error
	return product, err
}

// FindByIDWithImages loads a product with only the Images preload (used by Delete).
func (r ProductRepository) FindByIDWithImages(id uint) (models.Product, error) {
	var product models.Product
	err := r.db.Preload("Images").First(&product, id).Error
	return product, err
}

// PluckImagePaths returns only the Path column for all ProductImage rows of a product.
func (r ProductRepository) PluckImagePaths(productID uint) ([]string, error) {
	var paths []string
	err := r.db.Model(&models.ProductImage{}).
		Where("product_id = ?", productID).
		Pluck("path", &paths).Error
	return paths, err
}

// --- Create ---

// Create inserts the product (with nested Attributes and Variants) using the provided transaction.
func (r ProductRepository) Create(tx *gorm.DB, product *models.Product) error {
	return tx.Create(product).Error
}

// --- Update ---

// UpdateFields applies a column-map to the product row inside the provided transaction.
func (r ProductRepository) UpdateFields(tx *gorm.DB, product *models.Product, updates map[string]interface{}) error {
	return tx.Model(product).Updates(updates).Error
}

// UpdateImagePath sets the single "image" column (called after the main image file is committed).
func (r ProductRepository) UpdateImagePath(productID uint, path string) error {
	return r.db.Model(&models.Product{}).Where("id = ?", productID).Update("image", path).Error
}

// --- Associations ---

// ReplaceAttributes clears the many2many join and appends the new set inside the provided transaction.
func (r ProductRepository) ReplaceAttributes(tx *gorm.DB, product *models.Product, attributeIDs []uint) error {
	if err := tx.Model(product).Association("Attributes").Clear(); err != nil {
		return err
	}
	if len(attributeIDs) == 0 {
		return nil
	}
	attributes := make([]models.Attribute, len(attributeIDs))
	for i, id := range attributeIDs {
		attributes[i] = models.Attribute{ID: id}
	}
	return tx.Model(product).Association("Attributes").Append(attributes)
}

// --- Variants ---

// DeleteVariants removes all ProductVariant rows for a product (CASCADE handles inventory).
func (r ProductRepository) DeleteVariants(tx *gorm.DB, productID uint) error {
	return tx.Where("product_id = ?", productID).Delete(&models.ProductVariant{}).Error
}

// CreateVariants bulk-inserts variant rows.
func (r ProductRepository) CreateVariants(tx *gorm.DB, variants []models.ProductVariant) error {
	return tx.Create(&variants).Error
}

// --- Gallery images ---

// DeleteProductImages removes all ProductImage rows for a product.
func (r ProductRepository) DeleteProductImages(tx *gorm.DB, productID uint) error {
	return tx.Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error
}

// CreateProductImages bulk-inserts gallery image rows.
func (r ProductRepository) CreateProductImages(tx *gorm.DB, images []models.ProductImage) error {
	return tx.Create(&images).Error
}

// --- Delete ---

// SoftDelete issues GORM's soft-delete on the product row.
func (r ProductRepository) SoftDelete(product *models.Product) error {
	return r.db.Delete(product).Error
}

// --- Transaction ---

// Transaction exposes db.Transaction so the service never imports gorm for transaction control.
func (r ProductRepository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
