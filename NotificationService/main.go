package main

import (
	"encoding/json"
	"net/http"
)

type NotificationRequest struct {
	Amount int `json:"amount"`
	Paid   int `json:"paid"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}

func notificationHandler(w http.ResponseWriter, r *http.Request) {
	var req NotificationRequest
	json.NewDecoder(r.Body).Decode(&req)

	// pakai fungsi dari notification.go
	status := ValidateNotification(req.Amount, req.Paid)

	res := NotificationResponse{Status: status}
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/notification", notificationHandler)
	http.ListenAndServe(":8081", nil)
}