package main

import (
	"encoding/json"
	"net/http"
)

type WarehouseRequest struct {
	Stock int `json:"stock"`
}

type WarehouseResponse struct {
	Status string `json:"status"`
}

func warehouseHandler(w http.ResponseWriter, r *http.Request) {
	var req WarehouseRequest
	json.NewDecoder(r.Body).Decode(&req)

	status := CheckWarehouse(req.Stock)

	res := WarehouseResponse{Status: status}
	json.NewEncoder(w).Encode(res)
}

func main() {
	http.HandleFunc("/warehouse", warehouseHandler)
	http.ListenAndServe(":8083", nil)
}