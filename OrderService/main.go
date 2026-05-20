package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

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

type OrderResponse struct {
	UserID          int     `json:"user_id"`
	TrackingNumber  string  `json:"tracking_number"`
	ShippingCost    float64 `json:"shipping_cost"`
	TotalPrice      float64 `json:"total_price"`
	Status          string  `json:"status"`
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	json.NewDecoder(r.Body).Decode(&req)

	shippingCost := CalculateShippingCost(req.WeightKg, req.DistanceKm)
	totalPrice := CalculateTotalPrice(req.BasePrice, shippingCost)

	var orderID int64
	var trackingNumber string

	if db != nil {
		err := db.QueryRow(`
			INSERT INTO orders 
			(user_id, sender_name, sender_address, receiver_name, receiver_address, weight_kg, distance_km, base_price, shipping_cost, total_price, status)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,'CREATED')
			RETURNING id`,
			req.UserID, req.SenderName, req.SenderAddress,
			req.ReceiverName, req.ReceiverAddress,
			req.WeightKg, req.DistanceKm,
			req.BasePrice, shippingCost, totalPrice,
		).Scan(&orderID)

		if err == nil {
			trackingNumber = GenerateTrackingNumber(orderID)
			db.Exec(`UPDATE orders SET tracking_number=$1 WHERE id=$2`, trackingNumber, orderID)
		} else {
			trackingNumber = GenerateTrackingNumber(0)
		}
	} else {
		trackingNumber = GenerateTrackingNumber(0)
	}

	res := OrderResponse{
		UserID:         req.UserID,
		TrackingNumber: trackingNumber,
		ShippingCost:   shippingCost,
		TotalPrice:     totalPrice,
		Status:         "CREATED",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	err := InitDB()
	if err != nil {
		log.Println("DB not connected:", err)
	} else {
		log.Println("DB connected!")
		CreateTable()
	}

	http.HandleFunc("/order", orderHandler)
	fmt.Println("Order Service running on :8081")
	http.ListenAndServe(":8081", nil)
}