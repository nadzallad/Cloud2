package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Status string `json:"status"`
}

func TestPaymentAPI_Success(t *testing.T) {

	// =========================
	// HIT API
	// =========================
	jsonData := []byte(`{"amount":10000,"paid":10000}`)

	resp, err := http.Post(
		"http://localhost:8081/payment",
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var result Response
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status != "paid" {
		t.Errorf("Expected paid, got %s", result.Status)
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