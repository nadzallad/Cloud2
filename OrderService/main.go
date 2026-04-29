package main

import (
	"encoding/json"
	"net/http"
)

// Request dari client
type OrderRequest struct {
	UserID          int     `json:"user_id"`
	SenderName      string  `json:"sender_name"`
	SenderAddress   string  `json:"sender_address"`
	ReceiverName    string  `json:"receiver_name"`
	ReceiverAddress string  `json:"receiver_address"`
	WeightKg        float64 `json:"weight_kg"`
	DistanceKm      float64 `json:"distance_km"`
	BasePrice       float64 `json:"base_price"`
}

// Response ke client
type OrderResponse struct {
	UserID       int     `json:"user_id"`
	ShippingCost float64 `json:"shipping_cost"`
	TotalPrice   float64 `json:"total_price"`
	Status       string  `json:"status"`
}

// Handler
func orderHandler(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	json.NewDecoder(r.Body).Decode(&req)

	shippingCost := CalculateShippingCost(req.WeightKg, req.DistanceKm)
	totalPrice := CalculateTotalPrice(req.BasePrice, shippingCost)

	res := OrderResponse{
		UserID:       req.UserID,
		ShippingCost: shippingCost,
		TotalPrice:   totalPrice,
		Status:       "CREATED",
	}

	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/order", orderHandler)
	http.ListenAndServe(":8084", nil)
}