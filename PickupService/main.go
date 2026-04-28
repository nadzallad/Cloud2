package main

import (
	"encoding/json"
	"net/http"
)

type PickupRequest struct {
	OrderID       string `json:"order_id"`
	PaymentStatus string `json:"payment_status"`
	Weight        int    `json:"weight"`
}

type PickupResponse struct {
	Status string `json:"status"`
}

func pickupHandler(w http.ResponseWriter, r *http.Request) {
	var req PickupRequest
	json.NewDecoder(r.Body).Decode(&req)

	status := ProcessPickup(req.PaymentStatus, req.Weight)

	res := PickupResponse{Status: status}
	json.NewEncoder(w).Encode(res)
}

// ✅ TAMBAH INI
func ProcessPickup(paymentStatus string, weight int) string {
	if paymentStatus == "paid" {
		return "scheduled"
	}
	return "pending"
}

func main() {
	http.HandleFunc("/pickup", pickupHandler)
	http.ListenAndServe(":8082", nil)
}