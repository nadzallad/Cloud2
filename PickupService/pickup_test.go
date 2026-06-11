package main

import "testing"

func TestPickup_Success(t *testing.T) {
	result := ProcessPickup("paid", 2)

	if result != "scheduled" {
		t.Errorf("Expected scheduled, got %s", result)
	}
}

func TestPickup_Failed(t *testing.T) {
	result := ProcessPickup("unpaid", 0)

	if result != "pending" {
		t.Errorf("Expected pending, got %s", result)
	}
}
