package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestOrderEndpoint(t *testing.T) {
	reqBody := map[string]interface{}{
		"user_id":          1,
		"sender_name":      "Gherry",
		"sender_phone":     "08123456789",
		"sender_address":   "Bandung",
		"receiver_name":    "Budi",
		"receiver_phone":   "08987654321",
		"receiver_address": "Jakarta Pusat",
		"item_name":        "Laptop",
		"item_type":        "elektronik",
		"weight_kg":        2,
		"origin_city":      "Bandung",
		"destination_city": "Jakarta Pusat",
		"service_type":     "express",
		"base_price":       10000,
	}

	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://host.docker.internal:8081/order",
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

func TestConfirmPaymentEndpoint(t *testing.T) {
	reqBody := map[string]interface{}{
		"order_id": 1,
	}

	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post("http://host.docker.internal:8081/order/confirm-payment",
		"application/json",
		bytes.NewBuffer(jsonBody))

	if err != nil {
		t.Skipf("Order service not running: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 200 or 400, got %d", resp.StatusCode)
	}
}