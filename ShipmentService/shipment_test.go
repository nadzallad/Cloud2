package main

import "testing"

func TestCreateShipment(t *testing.T) {
	s := CreateShipment("SHP001")

	if s.ID != "SHP001" {
		t.Error("ID salah")
	}

	if s.Status != "Pending" {
		t.Error("Status harus Pending")
	}
}

func TestUpdateShipmentStatus(t *testing.T) {
	s := CreateShipment("SHP001")
	s = UpdateShipmentStatus(s, "Delivered")

	if s.Status != "Delivered" {
		t.Error("Status gagal update")
	}
}
