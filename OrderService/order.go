package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	connStr := "host=localhost port=5432 user=postgres password=admin123 dbname=order_db sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return db.Ping()
}

func CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		sender_name VARCHAR(100),
		sender_address TEXT,
		receiver_name VARCHAR(100),
		receiver_address TEXT,
		weight_kg DECIMAL(10,2),
		distance_km DECIMAL(10,2),
		base_price DECIMAL(10,2),
		shipping_cost DECIMAL(10,2),
		total_price DECIMAL(10,2),
		tracking_number VARCHAR(100) UNIQUE,
		status VARCHAR(50) DEFAULT 'CREATED',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

func GenerateTrackingNumber(orderID int64) string {
	return fmt.Sprintf("LOG-%d-%d", orderID, time.Now().Unix())
}

func CalculateShippingCost(weightKg float64, distanceKm float64) float64 {
	const pricePerKg = 5000
	const pricePerKm = 2000
	return (weightKg * pricePerKg) + (distanceKm * pricePerKm)
}

func CalculateTotalPrice(basePrice float64, shippingCost float64) float64 {
	return basePrice + shippingCost
}