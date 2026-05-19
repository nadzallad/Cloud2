package main

import "testing"

func TestValidatePayment_Success(t *testing.T) {
	result := ValidatePayment(10000, 10000)

	if result != "PAID" {
		t.Errorf("Expected PAID, got %s", result)
	}
}

func TestValidatePayment_Pending(t *testing.T) {
	result := ValidatePayment(10000, 0)

	if result != "PENDING" {
		t.Errorf("Expected PENDING, got %s", result)
	}
}

func TestValidatePayment_Failed(t *testing.T) {
	result := ValidatePayment(10000, 5000)

	if result != "FAILED" {
		t.Errorf("Expected FAILED, got %s", result)
	}
}