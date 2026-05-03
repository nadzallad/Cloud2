package main

import "testing"

func TestCreateDelivery(t *testing.T) {
	d := CreateDelivery("DLV001", "JNE")

	if d.ID != "DLV001" {
		t.Error("ID salah")
	}

	if d.Courier != "JNE" {
		t.Error("Courier salah")
	}

	if d.Status != "On Process" {
		t.Error("Status salah")
	}
}

func TestUpdateDeliveryStatus(t *testing.T) {
	d := CreateDelivery("DLV001", "JNE")
	d = UpdateDeliveryStatus(d, "Delivered")

	if d.Status != "Delivered" {
		t.Error("Update status gagal")
	}
}