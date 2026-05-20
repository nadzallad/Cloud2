package main

import (
	"encoding/json"
	"log"
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

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	status := ValidatePayment(req.Amount, req.Paid)

	res := PaymentResponse{
		Status: status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/payment", paymentHandler)

	log.Println("🚀 Server starting on :8082")

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}