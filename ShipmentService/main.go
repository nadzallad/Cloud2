package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Shipment Service Running")
}

func main() {
	ConnectDB()

	http.HandleFunc("/", handler)
	fmt.Println("Shipment Service running on port 8085")
	http.ListenAndServe(":8085", nil)
}
