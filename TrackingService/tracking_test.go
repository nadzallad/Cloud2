package main

import "testing"

func TestValidateTracking(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PICKED_UP", "PICKED_UP"},
		{"IN_TRANSIT", "IN_TRANSIT"},
		{"DELIVERED", "DELIVERED"},
		{"INVALID", "UNKNOWN"},
		{"", "UNKNOWN"},
	}

	for _, tt := range tests {
		result := ValidateTracking(tt.input)
		if result != tt.expected {
			t.Errorf("input %s: expected %s, got %s", tt.input, tt.expected, result)
		}
	}
}