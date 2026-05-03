package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestPickupEndpoint(t *testing.T) {
	reqBody := map[string]interface{}{
		"order_id":       "ORD1",
		"payment_status": "paid",
		"weight":         2,
	}

	body, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://localhost:8082/pickup", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	defer resp.Body.Close()

	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)

	if res["status"] != "scheduled" {
		t.Errorf("Expected scheduled, got %v", res["status"])
	}
}