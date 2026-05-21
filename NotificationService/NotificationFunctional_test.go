package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestNotificationAPI(t *testing.T) {

	reqBody := NotificationRequest{
		Amount: 10000,
		Paid:   10000,
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://test-notification:8088/notification",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Fatalf("Failed to call API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200, got %d", resp.StatusCode)
	}

	var res NotificationResponse
	json.NewDecoder(resp.Body).Decode(&res)

	if res.Status != "paid" {
		t.Errorf("Expected paid, got %s", res.Status)
	}
}
