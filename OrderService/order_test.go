package main

import "testing"

func TestCalculateShippingCost(t *testing.T) {
	result := CalculateShippingCost(2, 5)

	expected := float64((2 * 5000) + (5 * 2000))

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateTotalPrice(t *testing.T) {
	result := CalculateTotalPrice(10000, 20000)

	expected := float64(30000)

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}