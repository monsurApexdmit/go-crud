package dto

// CreateCategoryRequest is the DTO for creating a new category
type CreateCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"required,min=2,max=100"`
	ParentID     *uint  `json:"parent_id" binding:"omitempty,gt=0"`
	Status       *bool  `json:"status" binding:"omitempty"`
}

// UpdateCategoryRequest is the DTO for updating a category
type UpdateCategoryRequest struct {
	CategoryName string `json:"category_name" binding:"omitempty,min=2,max=100"`
	ParentID     *uint  `json:"parent_id" binding:"omitempty,gt=0"`
	Status       *bool  `json:"status" binding:"omitempty"`
}

// Validation rules explained:
// - required: field must be present and not empty
// - min=2: minimum length of 2 characters
// - max=100: maximum length of 100 characters
// - omitempty: field is optional
// - gt=0: greater than 0 (for parent_id)
