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

type NotificationRequest struct {
	Amount   int       `json:"amount" bson:"amount"`
	Paid     int       `json:"paid" bson:"paid"`
	Status   string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

type NotificationResponse struct {
	Status string `json:"status"`
}

var notificationCollection *mongo.Collection

func initMongo() {

	mongoURI := os.Getenv("MONGO_URI")

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

	db := client.Database("notification_db")
	notificationCollection = db.Collection("notifications")
}

func notificationHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req NotificationRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	req.Status = ValidateNotification(req.Amount, req.Paid)
	req.CreatedAt = time.Now()

	_, err = notificationCollection.InsertOne(
		context.TODO(),
		req,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(
		NotificationResponse{
			Status: req.Status,
		},
	)
}

func main() {

	initMongo()

	http.HandleFunc("/notification", notificationHandler)

	err := http.ListenAndServe(":8088", nil)

	if err != nil {
		panic(err)
	}
}
