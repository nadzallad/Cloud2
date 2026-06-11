package main

import (
	"encoding/json"
	"net/http"
)

type WarehouseRequest struct {
	UserID         int    `json:"user_id"`
	TrackingNumber string `json:"tracking_number"`
	ItemName       string `json:"item_name"`
	Stock          int    `json:"stock"`
}

type WarehouseResponse struct {
	Status string `json:"status"`
}

func warehouseHandler(w http.ResponseWriter, r *http.Request) {

	var req WarehouseRequest

	json.NewDecoder(r.Body).Decode(&req)

	status := CheckWarehouse(req.Stock)

	_, err := db.Exec(
		`INSERT INTO warehouse_logs
		(user_id, tracking_number, item_name, stock, status)
		VALUES ($1,$2,$3,$4,$5)`,

		req.UserID,
		req.TrackingNumber,
		req.ItemName,
		req.Stock,
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

	res := WarehouseResponse{
		Status: status,
	}

	json.NewEncoder(w).Encode(res)
}

func main() {

	InitDB()

	http.HandleFunc("/warehouse", warehouseHandler)

	http.ListenAndServe(":8084", nil)
}