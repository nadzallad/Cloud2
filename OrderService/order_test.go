package main

import (
	"math"
	"testing"
)

func TestCalculateShippingCostRegular(t *testing.T) {
	result := CalculateShippingCost(2, 100, "regular")
	expected := math.Round(2 * 100 * 100)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateShippingCostExpress(t *testing.T) {
	result := CalculateShippingCost(2, 100, "express")
	expected := math.Round(2 * 150 * 100)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateShippingCostSameDay(t *testing.T) {
	result := CalculateShippingCost(2, 100, "same_day")
	expected := math.Round(2 * 200 * 100)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateShippingCostMinimum(t *testing.T) {
	result := CalculateShippingCost(0.1, 1, "regular")
	expected := float64(15000)
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestCalculateShippingCostCeilWeight(t *testing.T) {
	// 1.3kg dibulatkan jadi 2kg
	result := CalculateShippingCost(1.3, 100, "regular")
	expected := math.Round(2 * 100 * 100)
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