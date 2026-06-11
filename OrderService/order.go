package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

const orsAPIKey = "eyJvcmciOiI1YjNjZTM1OTc4NTExMTAwMDFjZjYyNDgiLCJpZCI6IjZiOWFmZTFiYmVmYzRjMDlhMTkxNjUzMzQ3NDdkZDYyIiwiaCI6Im11cm11cjY0In0="

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
	CREATE TABLE IF NOT EXISTS provinces (
		province_id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS cities (
		city_id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		province_id INT NOT NULL REFERENCES provinces(province_id)
	);

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
		origin_city_id INT REFERENCES cities(city_id),
		destination_city_id INT REFERENCES cities(city_id),
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

func SeedData() error {
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM provinces`).Scan(&count)
	if count > 0 {
		return nil
	}

	_, err := db.Exec(`
	INSERT INTO provinces (name) VALUES
		('Jawa Barat'),
		('Jawa Tengah'),
		('Jawa Timur'),
		('DKI Jakarta'),
		('DI Yogyakarta'),
		('Banten'),
		('Sumatera Utara'),
		('Sumatera Selatan'),
		('Kalimantan Timur'),
		('Sulawesi Selatan');

	INSERT INTO cities (name, province_id) VALUES
		('Bandung', 1),
		('Bogor', 1),
		('Depok', 1),
		('Bekasi', 1),
		('Cimahi', 1),
		('Semarang', 2),
		('Solo', 2),
		('Magelang', 2),
		('Surabaya', 3),
		('Malang', 3),
		('Sidoarjo', 3),
		('Jakarta Pusat', 4),
		('Jakarta Selatan', 4),
		('Jakarta Utara', 4),
		('Jakarta Barat', 4),
		('Jakarta Timur', 4),
		('Yogyakarta', 5),
		('Sleman', 5),
		('Tangerang', 6),
		('Serang', 6),
		('Medan', 7),
		('Palembang', 8),
		('Samarinda', 9),
		('Balikpapan', 9),
		('Makassar', 10);
	`)
	return err
}

func GetCityID(cityName string) (int, error) {
	var cityID int
	err := db.QueryRow(`SELECT city_id FROM cities WHERE LOWER(name) = LOWER($1)`, cityName).Scan(&cityID)
	return cityID, err
}

func InsertCity(cityName string) (int, error) {
	// insert ke province "Lainnya" (id 99), buat dulu kalau belum ada
	db.Exec(`INSERT INTO provinces (province_id, name) VALUES (99, 'Lainnya') ON CONFLICT DO NOTHING`)
	var cityID int
	err := db.QueryRow(`
		INSERT INTO cities (name, province_id) VALUES ($1, 99) RETURNING city_id`,
		cityName,
	).Scan(&cityID)
	return cityID, err
}

type geocodeResult struct {
	Lat float64
	Lon float64
}

func geocodeCity(cityName string) (geocodeResult, error) {
	apiURL := fmt.Sprintf(
		"https://api.openrouteservice.org/geocode/search?api_key=%s&text=%s&size=1",
		orsAPIKey,
		url.QueryEscape(cityName+", Indonesia"),
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return geocodeResult{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Features []struct {
			Geometry struct {
				Coordinates []float64 `json:"coordinates"`
			} `json:"geometry"`
		} `json:"features"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return geocodeResult{}, err
	}

	if len(result.Features) == 0 {
		return geocodeResult{}, fmt.Errorf("city not found: %s", cityName)
	}

	coords := result.Features[0].Geometry.Coordinates
	return geocodeResult{Lon: coords[0], Lat: coords[1]}, nil
}

func GetDistanceORS(originCity string, destinationCity string) (float64, error) {
	origin, err := geocodeCity(originCity)
	if err != nil {
		return 0, fmt.Errorf("geocode origin failed: %v", err)
	}

	dest, err := geocodeCity(destinationCity)
	if err != nil {
		return 0, fmt.Errorf("geocode destination failed: %v", err)
	}

	apiURL := fmt.Sprintf(
		"https://api.openrouteservice.org/v2/directions/driving-car?api_key=%s&start=%f,%f&end=%f,%f",
		orsAPIKey,
		origin.Lon, origin.Lat,
		dest.Lon, dest.Lat,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Features []struct {
			Properties struct {
				Summary struct {
					Distance float64 `json:"distance"`
				} `json:"summary"`
			} `json:"properties"`
		} `json:"features"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	if len(result.Features) == 0 {
		return 0, fmt.Errorf("no route found")
	}

	distanceKm := result.Features[0].Properties.Summary.Distance / 1000
	return distanceKm, nil
}

func GenerateNoResi(orderID int64) string {
	return fmt.Sprintf("LOG-%d-%d", orderID, time.Now().Unix())
}

func CalculateShippingCost(weightKg float64, distanceKm float64, serviceType string) float64 {
	roundedWeight := math.Ceil(weightKg)

	var tarifPerKm float64
	switch serviceType {
	case "express":
		tarifPerKm = 150
	case "same_day":
		tarifPerKm = 200
	default:
		tarifPerKm = 100
	}

	ongkir := roundedWeight * tarifPerKm * distanceKm
	ongkir = math.Round(ongkir)

	if ongkir < 15000 {
		ongkir = 15000
	}

	return ongkir
}

func CalculateTotalPrice(basePrice float64, shippingCost float64) float64 {
	return basePrice + shippingCost
}