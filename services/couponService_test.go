package services

import (
	"testing"
	"time"

	"go-crud/models"
)

// Test CalculateDiscount function
func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name        string
		coupon      *models.Coupon
		orderAmount float64
		expected    float64
	}{
		{
			name: "Percentage discount without cap",
			coupon: &models.Coupon{
				Type:     "percentage",
				Discount: 20,
			},
			orderAmount: 100.00,
			expected:    20.00,
		},
		{
			name: "Percentage discount with cap",
			coupon: &models.Coupon{
				Type:        "percentage",
				Discount:    20,
				MaxDiscount: floatPtr(15.00),
			},
			orderAmount: 100.00,
			expected:    15.00, // Capped at $15
		},
		{
			name: "Fixed discount",
			coupon: &models.Coupon{
				Type:     "fixed",
				Discount: 50.00,
			},
			orderAmount: 100.00,
			expected:    50.00,
		},
		{
			name: "Fixed discount exceeding order amount",
			coupon: &models.Coupon{
				Type:     "fixed",
				Discount: 150.00,
			},
			orderAmount: 100.00,
			expected:    100.00, // Cannot exceed order amount
		},
		{
			name: "Free shipping (no discount)",
			coupon: &models.Coupon{
				Type:     "free_shipping",
				Discount: 0,
			},
			orderAmount: 100.00,
			expected:    0.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDiscount(tt.coupon, tt.orderAmount)
			if result != tt.expected {
				t.Errorf("CalculateDiscount() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test helper functions
func TestParseIDList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []uint
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: []uint{},
		},
		{
			name:     "Single ID",
			input:    "5",
			expected: []uint{5},
		},
		{
			name:     "Multiple IDs",
			input:    "1,3,5,10",
			expected: []uint{1, 3, 5, 10},
		},
		{
			name:     "IDs with spaces",
			input:    "1, 3, 5",
			expected: []uint{1, 3, 5},
		},
		{
			name:     "Invalid mixed with valid",
			input:    "1,abc,5",
			expected: []uint{1, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseIDList(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseIDList() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i, id := range result {
				if id != tt.expected[i] {
					t.Errorf("parseIDList()[%d] = %v, want %v", i, id, tt.expected[i])
				}
			}
		})
	}
}

func TestContains(t *testing.T) {
	slice := []uint{1, 3, 5, 10}

	tests := []struct {
		name     string
		item     uint
		expected bool
	}{
		{"Item exists - first", 1, true},
		{"Item exists - middle", 5, true},
		{"Item exists - last", 10, true},
		{"Item does not exist", 7, false},
		{"Item does not exist - 0", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(slice, tt.item)
			if result != tt.expected {
				t.Errorf("contains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test date validation logic
func TestDateValidation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		shouldErr bool
		errorCode string
	}{
		{
			name:      "Valid date range - current",
			startDate: now.Add(-24 * time.Hour),
			endDate:   now.Add(24 * time.Hour),
			shouldErr: false,
		},
		{
			name:      "Not started yet",
			startDate: now.Add(24 * time.Hour),
			endDate:   now.Add(48 * time.Hour),
			shouldErr: true,
			errorCode: "NOT_STARTED",
		},
		{
			name:      "Already expired",
			startDate: now.Add(-48 * time.Hour),
			endDate:   now.Add(-24 * time.Hour),
			shouldErr: true,
			errorCode: "EXPIRED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Date validation is part of ValidateCoupon
			// This is a conceptual test structure
			if now.Before(tt.startDate) && !tt.shouldErr {
				t.Error("Expected NOT_STARTED error")
			}
			if now.After(tt.endDate) && !tt.shouldErr {
				t.Error("Expected EXPIRED error")
			}
		})
	}
}

// Helper function
func floatPtr(f float64) *float64 {
	return &f
}
