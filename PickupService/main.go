package main

import (
	"encoding/json"
	"net/http"
)

type PickupRequest struct {
	UserID         int    `json:"user_id"`
	TrackingNumber string `json:"tracking_number"`
	PaymentStatus  string `json:"payment_status"`
	Weight         int    `json:"weight"`
}

type PickupResponse struct {
	Status string `json:"status"`
}

func pickupHandler(w http.ResponseWriter, r *http.Request) {

	var req PickupRequest

	json.NewDecoder(r.Body).Decode(&req)

	status := ProcessPickup(
		req.PaymentStatus,
		req.Weight,
	)

	_, err := db.Exec(
		`INSERT INTO pickups
		(user_id, tracking_number, payment_status, weight_kg, status)
		VALUES ($1,$2,$3,$4,$5)`,

		req.UserID,
		req.TrackingNumber,
		req.PaymentStatus,
		req.Weight,
		status,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	res := PickupResponse{
		Status: status,
	}

	json.NewEncoder(w).Encode(res)
}

func main() {
	InitDB()

	http.HandleFunc("/pickup", pickupHandler)
	http.ListenAndServe(":8089", nil)
}
