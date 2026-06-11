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
	SenderPhone     string  `json:"sender_phone"`
	SenderAddress   string  `json:"sender_address"`
	ReceiverName    string  `json:"receiver_name"`
	ReceiverPhone   string  `json:"receiver_phone"`
	ReceiverAddress string  `json:"receiver_address"`
	ItemName        string  `json:"item_name"`
	ItemType        string  `json:"item_type"`
	WeightKg        float64 `json:"weight_kg"`
	DistanceKm      float64 `json:"distance_km"`
	OriginCity      string  `json:"origin_city"`
	DestinationCity string  `json:"destination_city"`
	ServiceType     string  `json:"service_type"`
	BasePrice       float64 `json:"base_price"`
}

type OrderResponse struct {
	OrderID      int     `json:"order_id"`
	UserID       int     `json:"user_id"`
	ShippingCost float64 `json:"shipping_cost"`
	TotalPrice   float64 `json:"total_price"`
	ServiceType  string  `json:"service_type"`
	Status       string  `json:"status"`
}

type ConfirmPaymentRequest struct {
	OrderID int `json:"order_id"`
}

type ConfirmPaymentResponse struct {
	OrderID int    `json:"order_id"`
	NoResi  string `json:"no_resi"`
	Status  string `json:"status"`
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	json.NewDecoder(r.Body).Decode(&req)

	if req.ServiceType == "" {
		req.ServiceType = "regular"
	}

	shippingCost := CalculateShippingCost(req.WeightKg, req.DistanceKm, req.ServiceType)
	totalPrice := CalculateTotalPrice(req.BasePrice, shippingCost)

	var orderID int64

	if db != nil {
		err := db.QueryRow(`
			INSERT INTO orders 
			(user_id, sender_name, sender_phone, sender_address, receiver_name, receiver_phone, receiver_address, item_name, item_type, weight_kg, distance_km, origin_city, destination_city, service_type, base_price, shipping_cost, total_price, status)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,'WAITING_PAYMENT')
			RETURNING order_id`,
			req.UserID, req.SenderName, req.SenderPhone, req.SenderAddress,
			req.ReceiverName, req.ReceiverPhone, req.ReceiverAddress,
			req.ItemName, req.ItemType, req.WeightKg, req.DistanceKm,
			req.OriginCity, req.DestinationCity, req.ServiceType,
			req.BasePrice, shippingCost, totalPrice,
		).Scan(&orderID)

		if err != nil {
			log.Println("DB error:", err)
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}
	}

	res := OrderResponse{
		OrderID:      int(orderID),
		UserID:       req.UserID,
		ShippingCost: shippingCost,
		TotalPrice:   totalPrice,
		ServiceType:  req.ServiceType,
		Status:       "WAITING_PAYMENT",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func confirmPaymentHandler(w http.ResponseWriter, r *http.Request) {
	var req ConfirmPaymentRequest
	json.NewDecoder(r.Body).Decode(&req)

	if db == nil {
		http.Error(w, "DB not connected", http.StatusInternalServerError)
		return
	}

	// cek order ada dan statusnya WAITING_PAYMENT
	var status string
	err := db.QueryRow(`SELECT status FROM orders WHERE order_id=$1`, req.OrderID).Scan(&status)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if status != "WAITING_PAYMENT" {
		http.Error(w, "Order already processed", http.StatusBadRequest)
		return
	}

	// generate no_resi dan insert ke tabel resi
	noResi := GenerateNoResi(int64(req.OrderID))
	_, err = db.Exec(`INSERT INTO resi (order_id, no_resi) VALUES ($1, $2)`, req.OrderID, noResi)
	if err != nil {
		log.Println("Resi error:", err)
		http.Error(w, "Failed to create resi", http.StatusInternalServerError)
		return
	}

	// update status order jadi PAID
	db.Exec(`UPDATE orders SET status='PAID' WHERE order_id=$1`, req.OrderID)

	res := ConfirmPaymentResponse{
		OrderID: req.OrderID,
		NoResi:  noResi,
		Status:  "PAID",
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
	http.HandleFunc("/order/confirm-payment", confirmPaymentHandler)
	fmt.Println("Order Service running on :8081")
	http.ListenAndServe(":8081", nil)
}