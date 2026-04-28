package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type trackingResponse struct {
	Status string `json:"status"`
}

func TestTrackingAPI_Success(t *testing.T) {

	// =========================
	// HIT API
	// =========================
	jsonData := []byte(`{
		"shipment_id":1,
		"tracking_number":"LOG-1-123456",
		"status":"IN_TRANSIT",
		"location":"Jakarta",
		"note":"Paket sedang dikirim"
	}`)

	resp, err := http.Post(
		"http://localhost:8083/tracking",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result trackingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != "IN_TRANSIT" {
		t.Errorf("Expected IN_TRANSIT, got %s", result.Status)
	}

	// =========================
	// CEK DATABASE
	// =========================
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_logistic")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var status string

	err = db.QueryRow(`
		SELECT status 
		FROM tracking_logs 
		ORDER BY id DESC 
		LIMIT 1
	`).Scan(&status)

	if err != nil {
		t.Fatal(err)
	}

	if status != "IN_TRANSIT" {
		t.Errorf("Expected DB status IN_TRANSIT, got %s", status)
	}
}