package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type NotificationRequest struct {
	UserID  int    `json:"user_id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}

func TestNotificationAPI(t *testing.T) {

	reqBody := NotificationRequest{
		UserID:  1,
		Type:    "PAYMENT_SUCCESS",
		Message: "Pembayaran berhasil",
	}

	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://localhost:8083/notification",
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

	if res.Status != "sent" {
		t.Errorf("Expected sent, got %s", res.Status)
	}
}