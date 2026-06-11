package main

import (
	"bytes"
	"net/http"
	"testing"
)

func TestWarehouseAPI(t *testing.T) {
	json := []byte(`{
		"user_id":1,
		"tracking_number":"TRK001",
		"item_name":"Laptop",
		"stock":10
	}`)

	resp, err := http.Post(
		"http://localhost:8084/warehouse",
		"application/json",
		bytes.NewBuffer(json),
	)

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}
