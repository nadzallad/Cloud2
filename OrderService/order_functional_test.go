package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestOrderEndpoint(t *testing.T) {
	reqBody := map[string]interface{}{
		"user_id":     1,
		"weight_kg":   2,
		"distance_km": 5,
		"base_price":  10000,
	}

	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8084/order",
		"application/json",
		bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Skipf("Order service not running: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}