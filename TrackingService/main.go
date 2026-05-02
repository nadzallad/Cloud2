package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type TrackingRequest struct {
	ShipmentID    int    `json:"shipment_id"`
	TrackingNumber string `json:"tracking_number"`
	Status        string `json:"status"`
	Location      string `json:"location"`
	Note          string `json:"note"`
}

type TrackingResponse struct {
	Status string `json:"status"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/db_logistic")
	if err != nil {
		panic(err)
	}
}

func trackingHandler(w http.ResponseWriter, r *http.Request) {
	var req TrackingRequest
	json.NewDecoder(r.Body).Decode(&req)

	status := ValidateTracking(req.Status)

	// insert ke DB
	_, err := db.Exec(`
		INSERT INTO tracking_logs 
		(shipment_id, tracking_number, status, location, note, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		req.ShipmentID,
		req.TrackingNumber,
		status,
		req.Location,
		req.Note,
		time.Now(),
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := TrackingResponse{Status: status}
	json.NewEncoder(w).Encode(res)
}

func main() {
	initDB()

	http.HandleFunc("/tracking", trackingHandler)
	http.ListenAndServe(":8087", nil)
}
