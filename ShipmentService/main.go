package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Shipment Service Running")
	})

	http.ListenAndServe(":8080", nil)
}
