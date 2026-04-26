package main

import (
	"encoding/json"
	"net/http"
)

type PaymentRequest struct {
	Amount int `json:"amount"`
	Paid   int `json:"paid"`
}

type PaymentResponse struct {
	Status string `json:"status"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	status := "pending"
	if req.Paid >= req.Amount {
		status = "paid"
	}

	res := PaymentResponse{Status: status}
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/payment", paymentHandler)
	http.ListenAndServe(":8081", nil)
}