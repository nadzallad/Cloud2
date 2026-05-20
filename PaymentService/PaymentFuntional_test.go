package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

type Response struct {
	Status string `json:"status"`
}

func TestPaymentAPI_Success(t *testing.T) {

	// WAIT API READY
	ready := false

	
	for i := 0; i < 5; i++ {
		resp, err := http.Post(
			"http://host.docker.internal:8082/payment",
			"application/json",
			bytes.NewBuffer([]byte(`{"amount":1,"paid":1}`)),
		)

		if err == nil && resp.StatusCode == 200 {
			ready = true 
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	if !ready {
		t.Fatal("API NOT READY")
	}

	// HIT API
	jsonData := []byte(`{
		"amount":10000,
		"paid":10000
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

	// DB CHECK
	db, err := sql.Open("postgres",
		"host=host.docker.internal port=5432 user=postgres password=admin123 dbname=payment_db sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	var status string
	err = db.QueryRow(`
		SELECT status FROM payments ORDER BY id DESC LIMIT 1
	`).Scan(&status)

	if err != nil {
		t.Fatal(err)
	}

	if status != "PAID" {
		t.Errorf("Expected DB status PAID, got %s", status)
	}
}