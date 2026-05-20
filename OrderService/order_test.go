package main

import "testing"

func TestCalculateShippingCostRegular(t *testing.T) {
	result := CalculateShippingCost(2, 100, "regular")
	expected := float64((2 * 5000) + (100 * 1000))
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateShippingCostExpress(t *testing.T) {
	result := CalculateShippingCost(2, 100, "express")
	expected := float64((2*5000)+(100*1000)) * 1.5
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateShippingCostSameDay(t *testing.T) {
	result := CalculateShippingCost(2, 100, "same_day")
	expected := float64((2*5000)+(100*1000)) * 2
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateTotalPrice(t *testing.T) {
	result := CalculateTotalPrice(10000, 110000)
	expected := float64(120000)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}