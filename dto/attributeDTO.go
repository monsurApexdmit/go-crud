package dto

// CreateAttributeRequest represents the request body for creating an attribute
type CreateAttributeRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=150"`
	OptionType  string `json:"option_type" binding:"required,oneof=text dropdown radio checkbox color size"`
	Values      string `json:"values" binding:"omitempty"` // Required for dropdown, radio, checkbox
	Description string `json:"description" binding:"omitempty,max=500"`
	IsRequired  *bool  `json:"is_required" binding:"omitempty"`
	Status      *bool  `json:"status" binding:"omitempty"`
	SortOrder   *int   `json:"sort_order" binding:"omitempty,min=0"`
}

// UpdateAttributeRequest represents the request body for updating an attribute
type UpdateAttributeRequest struct {
	Name        string `json:"name" binding:"omitempty,min=2,max=100"`
	DisplayName string `json:"display_name" binding:"omitempty,min=2,max=150"`
	OptionType  string `json:"option_type" binding:"omitempty,oneof=text dropdown radio checkbox color size"`
	Values      string `json:"values" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty,max=500"`
	IsRequired  *bool  `json:"is_required" binding:"omitempty"`
	Status      *bool  `json:"status" binding:"omitempty"`
	SortOrder   *int   `json:"sort_order" binding:"omitempty,min=0"`
}
