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

	_, err = db.Exec("DELETE FROM payments")
	if err != nil {
		http.Error(w, "Cleanup failed", 500)
		return
	}
	
	// 🔥 INSERT KE DB
	_, err = db.Exec(
		"INSERT INTO payments (amount, paid, status) VALUES ($1, $2, $3)",
		req.Amount, req.Paid, status,
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

	// 🔥 INIT DB SEKALI
	db, err = sql.Open("postgres",
		"host=host.docker.internal port=5432 user=postgres password=admin123 dbname=payment_db sslmode=disable")
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// 🔥 PASTIKAN DB READY
	err = db.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}

	log.Println("✅ Connected to DB")

	http.HandleFunc("/payment", paymentHandler)

	log.Println("🚀 Server starting on :8082")

	err = http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatal("❌ Server failed:", err)
	}
}