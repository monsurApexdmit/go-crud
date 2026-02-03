package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"go-crud/models"
	"go-crud/repositories"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const productImageFolder = "uploads/products"

// ProductForm is the set of scalar fields that arrive on both POST and PUT.
// The controller binds this on Create; on Update it is unused (the controller
// builds a map instead, because partial-update semantics need to distinguish
// "not sent" from "sent as zero").
type ProductForm struct {
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

// ProductListParams carries parsed-and-clamped pagination values and filters.
type ProductListParams struct {
	Filters repositories.ProductFilters
	Page    int
	Limit   int
	Offset  int
}

// ProductService owns validation, file-upload orchestration, and the
// "save to temp → commit DB → commit files" sequencing.
type ProductService struct {
	repo repositories.ProductRepository
	fs   FileService
}

func NewProductService() ProductService {
	return ProductService{
		repo: repositories.NewProductRepository(),
		fs:   NewFileService(),
	}
}

// --- List & Get ---

// ListProducts clamps page/limit, computes offset, and delegates to the repository.
func (s ProductService) ListProducts(params ProductListParams) ([]models.Product, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	params.Offset = (params.Page - 1) * params.Limit

	return s.repo.FindAll(params.Filters, params.Page, params.Limit, params.Offset)
}

// GetProduct loads a single product by ID.
func (s ProductService) GetProduct(id uint) (models.Product, error) {
	return s.repo.FindByID(id)
}

// --- Create ---

// CreateProduct orchestrates the full create flow:
//  1. Validate required fields.
//  2. Save uploaded files to temp.
//  3. Parse attributes and variants from JSON strings.
//  4. DB transaction: insert product + nested relations.
//  5. Commit main image file, update DB with final path.
//  6. Commit gallery images, insert ProductImage rows.
//  7. Reload and return the product with all relations.
func (s ProductService) CreateProduct(
	form ProductForm,
	mainImage *multipart.FileHeader,
	galleryFiles []*multipart.FileHeader,
	attributesJSON string,
	variantsJSON string,
	c *gin.Context,
) (models.Product, error) {
	if form.Name == "" {
		return models.Product{}, errors.New("product name is required")
	}

	// --- 2. Save files to temp ---
	var mainTmp tempFile
	var hasMainImage bool
	if mainImage != nil {
		var err error
		mainTmp, err = s.fs.SaveTemp(c, mainImage, productImageFolder)
		if err != nil {
			log.Printf("CreateProduct: save main image to temp: %v", err)
			return models.Product{}, fmt.Errorf("failed to save image: %w", err)
		}
		hasMainImage = true
	}

	galleryTmps, err := s.savePending(c, galleryFiles)
	if err != nil {
		// clean up main temp if it was saved
		if hasMainImage {
			s.fs.RemoveFile(mainTmp.tempPath)
		}
		return models.Product{}, fmt.Errorf("failed to save images: %w", err)
	}

	// --- 3. Parse JSON sub-fields ---
	product := models.Product{
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

	if ids := parseAttributeIDs(attributesJSON); len(ids) > 0 {
		attrs := make([]models.Attribute, len(ids))
		for i, id := range ids {
			attrs[i] = models.Attribute{ID: id}
		}
		product.Attributes = attrs
	}

	product.Variants = parseVariants(0, variantsJSON) // ProductID set by GORM after insert

	// --- 4. DB transaction ---
	if err := s.repo.Transaction(func(tx *gorm.DB) error {
		return s.repo.Create(tx, &product)
	}); err != nil {
		s.cleanupTemps(mainTmp, hasMainImage, galleryTmps)
		log.Printf("CreateProduct: db create: %v", err)
		return models.Product{}, fmt.Errorf("failed to create product: %w", err)
	}

	// --- 5. Commit main image ---
	if hasMainImage {
		if err := s.fs.Commit(mainTmp); err != nil {
			s.fs.RemoveFile(mainTmp.tempPath)
			s.cleanupTemps(tempFile{}, false, galleryTmps)
			log.Printf("CreateProduct: commit main image: %v", err)
			return models.Product{}, fmt.Errorf("failed to store image: %w", err)
		}
		if err := s.repo.UpdateImagePath(product.ID, mainTmp.finalPath); err != nil {
			log.Printf("CreateProduct: update image path in DB: %v", err)
			return models.Product{}, fmt.Errorf("failed to update image: %w", err)
		}
		product.Image = mainTmp.finalPath
	}

	// --- 6. Commit gallery ---
	if len(galleryTmps) > 0 {
		finalPaths, err := s.fs.CommitAll(galleryTmps)
		if err != nil {
			log.Printf("CreateProduct: commit gallery: %v", err)
			return models.Product{}, fmt.Errorf("failed to store images: %w", err)
		}

		imageRows := make([]models.ProductImage, len(finalPaths))
		for i, path := range finalPaths {
			imageRows[i] = models.ProductImage{
				ProductID: product.ID,
				Path:      path,
				Position:  i,
				IsPrimary: i == 0,
			}
		}

		if err := s.repo.Transaction(func(tx *gorm.DB) error {
			return s.repo.CreateProductImages(tx, imageRows)
		}); err != nil {
			// roll back files
			for _, path := range finalPaths {
				s.fs.RemoveFile(path)
			}
			log.Printf("CreateProduct: insert gallery rows: %v", err)
			return models.Product{}, fmt.Errorf("failed to save images: %w", err)
		}
		product.Images = imageRows
	}

	// --- 7. Reload with all relations ---
	product, err = s.repo.FindByID(product.ID)
	if err != nil {
		log.Printf("CreateProduct: reload: %v", err)
	}

	return product, nil
}

// --- Update ---

// UpdateProduct orchestrates the full update flow:
//  1. Fetch existing product (returns error if not found).
//  2. Save any new images to temp.
//  3. DB transaction: apply field updates, replace attributes, replace variants.
//  4. Post-transaction: commit main image, update image column, remove old file.
//  5. Post-transaction: commit gallery, replace ProductImage rows, remove old gallery files.
//  6. Reload and return.
func (s ProductService) UpdateProduct(
	id uint,
	updates map[string]interface{},
	mainImage *multipart.FileHeader,
	galleryFiles []*multipart.FileHeader,
	attributesJSON string,
	variantsJSON string,
	c *gin.Context,
) (models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return models.Product{}, err // caller maps gorm.ErrRecordNotFound → 404
	}
	oldMainImagePath := product.Image

	// --- 2. Save new images to temp ---
	var mainTmp tempFile
	var hasMainImage bool
	if mainImage != nil {
		mainTmp, err = s.fs.SaveTemp(c, mainImage, productImageFolder)
		if err != nil {
			log.Printf("UpdateProduct: save main image to temp: %v", err)
			return models.Product{}, fmt.Errorf("failed to save image: %w", err)
		}
		hasMainImage = true
	}

	galleryTmps, err := s.savePending(c, galleryFiles)
	if err != nil {
		if hasMainImage {
			s.fs.RemoveFile(mainTmp.tempPath)
		}
		return models.Product{}, fmt.Errorf("failed to save images: %w", err)
	}

	// --- 3. DB transaction ---
	if err := s.repo.Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := s.repo.UpdateFields(tx, &product, updates); err != nil {
				return err
			}
		}

		if attributesJSON != "" {
			ids := parseAttributeIDs(attributesJSON)
			if err := s.repo.ReplaceAttributes(tx, &product, ids); err != nil {
				return err
			}
		}

		if variantsJSON != "" {
			variants := parseVariants(product.ID, variantsJSON)
			if len(variants) > 0 {
				if err := s.repo.DeleteVariants(tx, product.ID); err != nil {
					return err
				}
				if err := s.repo.CreateVariants(tx, variants); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		s.cleanupTemps(mainTmp, hasMainImage, galleryTmps)
		log.Printf("UpdateProduct: db transaction: %v", err)
		return models.Product{}, fmt.Errorf("failed to update product: %w", err)
	}

	// --- 4. Commit main image ---
	if hasMainImage {
		if err := s.fs.Commit(mainTmp); err != nil {
			s.fs.RemoveFile(mainTmp.tempPath)
			log.Printf("UpdateProduct: commit main image: %v", err)
			return models.Product{}, fmt.Errorf("failed to store image: %w", err)
		}
		if err := s.repo.UpdateImagePath(product.ID, mainTmp.finalPath); err != nil {
			log.Printf("UpdateProduct: update image path in DB: %v", err)
			return models.Product{}, fmt.Errorf("failed to update image: %w", err)
		}
		s.fs.RemoveFile(oldMainImagePath)
		product.Image = mainTmp.finalPath
	}

	// --- 5. Commit gallery ---
	if len(galleryTmps) > 0 {
		oldImagePaths, err := s.repo.PluckImagePaths(product.ID)
		if err != nil {
			log.Printf("UpdateProduct: pluck old image paths: %v", err)
			return models.Product{}, fmt.Errorf("failed to update images: %w", err)
		}

		finalPaths, err := s.fs.CommitAll(galleryTmps)
		if err != nil {
			log.Printf("UpdateProduct: commit gallery: %v", err)
			return models.Product{}, fmt.Errorf("failed to store images: %w", err)
		}

		imageRows := make([]models.ProductImage, len(finalPaths))
		for i, path := range finalPaths {
			imageRows[i] = models.ProductImage{
				ProductID: product.ID,
				Path:      path,
				Position:  i,
				IsPrimary: i == 0,
			}
		}

		if err := s.repo.Transaction(func(tx *gorm.DB) error {
			if err := s.repo.DeleteProductImages(tx, product.ID); err != nil {
				return err
			}
			return s.repo.CreateProductImages(tx, imageRows)
		}); err != nil {
			for _, path := range finalPaths {
				s.fs.RemoveFile(path)
			}
			log.Printf("UpdateProduct: replace gallery rows: %v", err)
			return models.Product{}, fmt.Errorf("failed to update images: %w", err)
		}

		// remove old gallery files after DB confirmed
		for _, p := range oldImagePaths {
			s.fs.RemoveFile(p)
		}
		product.Images = imageRows
	}

	// --- 6. Reload ---
	product, err = s.repo.FindByID(product.ID)
	if err != nil {
		log.Printf("UpdateProduct: reload: %v", err)
	}

	return product, nil
}

// --- Delete ---

// DeleteProduct fetches the product with images, soft-deletes it, then
// removes all image files from disk (best-effort).
func (s ProductService) DeleteProduct(id uint) error {
	product, err := s.repo.FindByIDWithImages(id)
	if err != nil {
		return err
	}

	if err := s.repo.SoftDelete(&product); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// best-effort file cleanup
	s.fs.RemoveFile(product.Image)
	for _, img := range product.Images {
		s.fs.RemoveFile(img.Path)
	}

	return nil
}

// --- unexported helpers ---

// CollectGalleryFiles extracts the multipart form from the request context,
// deduplicates file headers by filename+size, and returns the unique set.
// Called by the controller to hand off the file headers before invoking the service.
func CollectGalleryFiles(c *gin.Context) []*multipart.FileHeader {
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return nil
	}

	seen := make(map[string]struct{})
	var out []*multipart.FileHeader

	for key, fhs := range form.File {
		if !strings.HasPrefix(key, "image") {
			continue
		}
		for _, fh := range fhs {
			uniqueKey := fmt.Sprintf("%s_%d", fh.Filename, fh.Size)
			if _, exists := seen[uniqueKey]; exists {
				continue
			}
			seen[uniqueKey] = struct{}{}
			out = append(out, fh)
		}
	}
	return out
}

// savePending writes each file to a temp location via FileService and
// collects the pending tokens. On the first failure it cleans up all
// temp files that were already saved.
func (s ProductService) savePending(c *gin.Context, files []*multipart.FileHeader) ([]tempFile, error) {
	if len(files) == 0 {
		return nil, nil
	}
	pending := make([]tempFile, 0, len(files))
	for _, f := range files {
		tf, err := s.fs.SaveTemp(c, f, productImageFolder)
		if err != nil {
			// roll back previously saved temps
			for _, p := range pending {
				s.fs.RemoveFile(p.tempPath)
			}
			return nil, err
		}
		pending = append(pending, tf)
	}
	return pending, nil
}

// parseAttributeIDs unmarshals the "attributes" form value (a JSON array of uint)
// into a Go slice. Returns nil if the string is empty or malformed.
func parseAttributeIDs(raw string) []uint {
	if raw == "" {
		return nil
	}
	var ids []uint
	if err := json.Unmarshal([]byte(raw), &ids); err != nil {
		log.Printf("parseAttributeIDs: %v", err)
		return nil
	}
	return ids
}

// variantInput mirrors the shape of each element in the "variants" JSON array.
type variantInput struct {
	Name       string          `json:"name"`
	Attributes json.RawMessage `json:"attributes"`
	Price      float64         `json:"price"`
	SalePrice  float64         `json:"sale_price"`
	Stock      int             `json:"stock"`
	SKU        string          `json:"sku"`
	Barcode    string          `json:"barcode"`
}

// parseVariants unmarshals the "variants" form value into a typed slice with
// ProductID set. Returns nil if the string is empty or malformed.
func parseVariants(productID uint, raw string) []models.ProductVariant {
	if raw == "" {
		return nil
	}
	var inputs []variantInput
	if err := json.Unmarshal([]byte(raw), &inputs); err != nil {
		log.Printf("parseVariants: %v", err)
		return nil
	}
	variants := make([]models.ProductVariant, len(inputs))
	for i, v := range inputs {
		variants[i] = models.ProductVariant{
			ProductID:  productID,
			Name:       v.Name,
			Attributes: datatypes.JSON(v.Attributes),
			Price:      v.Price,
			SalePrice:  v.SalePrice,
			Stock:      v.Stock,
			SKU:        v.SKU,
			Barcode:    v.Barcode,
		}
	}
	return variants
}

// cleanupTemps removes all pending temp files after a DB failure prevents commit.
func (s ProductService) cleanupTemps(main tempFile, hasMain bool, gallery []tempFile) {
	if hasMain {
		s.fs.RemoveFile(main.tempPath)
	}
	for _, tf := range gallery {
		s.fs.RemoveFile(tf.tempPath)
	}
}
