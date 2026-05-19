package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/lib/pq"
)

type Response struct {
	Status string `json:"status"`
}

func TestPaymentAPI_Success(t *testing.T) {

	// =========================
	// HIT API
	// =========================
	jsonData := []byte(`{
		"order_id":1,
		"amount":10000,
		"paid":10000,
		"payment_method":"BANK_TRANSFER"
	}`)

	resp, err := http.Post(
		"http://localhost:8082/payment",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != "PAID" {
		t.Errorf("Expected PAID, got %s", result.Status)
	}

	// =========================
	// CEK DATABASE (POSTGRES)
	// =========================
	db, err := sql.Open("postgres",
		"host=localhost port=5432 user=postgres password=1234 dbname=payment_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var status string

	err = db.QueryRow(`
		SELECT status 
		FROM payments 
		ORDER BY id DESC 
		LIMIT 1
	`).Scan(&status)

	if err != nil {
		t.Fatal(err)
	}

	if status != "PAID" {
		t.Errorf("Expected DB status PAID, got %s", status)
	}
}