package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delivery Service Running")
}

func createDelivery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	d := CreateDelivery("DLV001", "Courier A")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func main() {
	ConnectDB()

	http.HandleFunc("/", home)
	http.HandleFunc("/delivery", createDelivery)

	fmt.Println("Delivery Service running on port 8086")
	http.ListenAndServe(":8086", nil)
}
