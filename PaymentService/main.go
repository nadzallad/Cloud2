package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

type PaymentRequest struct {
	Amount int `json:"amount"`
	Paid   int `json:"paid"`
}

type PaymentResponse struct {
	Status string `json:"status"`
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	status := ValidatePayment(req.Amount, req.Paid)

	// 🔥 CLEAN TABLE (ANTI ERROR TEST)
	_, err = db.Exec("TRUNCATE TABLE payments RESTART IDENTITY")
	if err != nil {
		http.Error(w, "Cleanup failed", 500)
		return
	}

	// 🔥 INSERT KE DB
	_, err = db.Exec(
		`INSERT INTO payments 
		(order_id, amount, paid_amount, payment_method, status, transaction_id, created_at, paid_at) 
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		1, req.Amount, req.Paid, "BANK_TRANSFER", status, "TRX123",
	)

	if err != nil {
		http.Error(w, "Insert failed", 500)
		return
	}

	res := PaymentResponse{
		Status: status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func main() {
	var err error

	// 🔥 INIT DB
	db, err = sql.Open("postgres",
		"host=host.docker.internal port=5432 user=postgres password=admin123 dbname=payment_db sslmode=disable")
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// 🔥 CEK KONEKSI
	err = db.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("✅ Connected to DB")

	// 🔥 AUTO CREATE TABLE (BIAR GA ERROR)
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS payments (
		id SERIAL PRIMARY KEY,
		amount INT,
		paid INT,
		status VARCHAR(20)
	)
	`)
	if err != nil {
		log.Fatal("Create table failed:", err)
	}

	http.HandleFunc("/payment", paymentHandler)

	log.Println("🚀 Server starting on :8082")

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}