package main

import "testing"

func TestValidatePayment_Success(t *testing.T) {
	result := ValidatePayment(10000, 10000)

	if result != "paid" {
		t.Errorf("Expected paid, got %s", result)
	}
}

func TestValidatePayment_Pending(t *testing.T) {
	result := ValidatePayment(10000, 5000)

	if result != "pending" {
		t.Errorf("Expected pending, got %s", result)
	}
}