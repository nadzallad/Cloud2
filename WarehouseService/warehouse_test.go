package main

import "testing"

func TestCheckWarehouse_Available(t *testing.T) {
	result := CheckWarehouse(10)

	if result != "available" {
		t.Errorf("Expected available, got %s", result)
	}
}

func TestCheckWarehouse_OutOfStock(t *testing.T) {
	result := CheckWarehouse(0)

	if result != "out_of_stock" {
		t.Errorf("Expected out_of_stock, got %s", result)
	}
}
