package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrderEndpoint(t *testing.T) {
	// bikin request body
	reqBody := map[string]interface{}{
		"user_id":     1,
		"weight_kg":   2,
		"distance_km": 5,
		"base_price":  10000,
	}

	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// panggil handler langsung
	orderHandler(rr, req)

	// cek status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// decode response
	var response map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&response)

	// cek field penting
	if response["status"] != "CREATED" {
		t.Errorf("Expected status CREATED, got %v", response["status"])
	}
}