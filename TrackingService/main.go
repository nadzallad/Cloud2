package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrackingRequest struct {
	ShipmentID     int       `json:"shipment_id" bson:"shipment_id"`
	TrackingNumber string    `json:"tracking_number" bson:"tracking_number"`
	Status         string    `json:"status" bson:"status"`
	Location       string    `json:"location" bson:"location"`
	Note           string    `json:"note" bson:"note"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

type TrackingResponse struct {
	Status string `json:"status"`
}

var trackingCollection *mongo.Collection

func initMongo() {

	// ambil URI dari environment variable
	mongoURI := os.Getenv("MONGO_URI")

	// fallback default untuk local development
	if mongoURI == "" {
		mongoURI = "mongodb://admin:admin123@localhost:27017/?authSource=admin"
	}

	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(mongoURI),
	)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	db := client.Database("tracking_db")
	trackingCollection = db.Collection("tracking_logs")
}

func trackingHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TrackingRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.TrackingNumber == "" {
		http.Error(w, "tracking number required", http.StatusBadRequest)
		return
	}

	req.Status = ValidateTracking(req.Status)
	req.CreatedAt = time.Now()

	_, err = trackingCollection.InsertOne(
		context.TODO(),
		req,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		TrackingResponse{
			Status: req.Status,
		},
	)
}

func main() {

	initMongo()

	http.HandleFunc("/tracking", trackingHandler)

	err := http.ListenAndServe(":8087", nil)
	if err != nil {
		panic(err)
	}
}
