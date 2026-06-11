package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	password := os.Getenv("DB_PASS")
	if password == "" {
		password = "admin123"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "order_db"
	}

	connStr := fmt.Sprintf("host=%s port=5432 user=postgres password=%s dbname=%s sslmode=disable", host, password, dbname)
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
		order_id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		sender_name VARCHAR(100),
		sender_phone VARCHAR(20),
		sender_address TEXT,
		receiver_name VARCHAR(100),
		receiver_phone VARCHAR(20),
		receiver_address TEXT,
		item_name VARCHAR(100),
		item_type VARCHAR(50),
		weight_kg DECIMAL(10,2),
		distance_km DECIMAL(10,2),
		origin_city VARCHAR(100),
		destination_city VARCHAR(100),
		service_type VARCHAR(50) DEFAULT 'regular',
		base_price DECIMAL(10,2),
		shipping_cost DECIMAL(10,2),
		total_price DECIMAL(10,2),
		status VARCHAR(50) DEFAULT 'WAITING_PAYMENT',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS resi (
		resi_id SERIAL PRIMARY KEY,
		order_id INT NOT NULL REFERENCES orders(order_id),
		no_resi VARCHAR(100) UNIQUE NOT NULL,
		issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(50) DEFAULT 'ACTIVE'
	);`
	_, err := db.Exec(query)
	return err
}

func GenerateNoResi(orderID int64) string {
	return fmt.Sprintf("LOG-%d-%d", orderID, time.Now().Unix())
}

func CalculateShippingCost(weightKg float64, distanceKm float64, serviceType string) float64 {
	base := (weightKg * 5000) + (distanceKm * 1000)

	switch serviceType {
	case "express":
		return base * 1.5
	case "same_day":
		return base * 2
	default:
		return base
	}
}

func CalculateTotalPrice(basePrice float64, shippingCost float64) float64 {
	return basePrice + shippingCost
}