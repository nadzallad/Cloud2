package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "testing"
)

type trackingResponse struct {
    Status string `json:"status"`
}

func TestTrackingAPI_Success(t *testing.T) {
    if testing.Short() {
        t.Skip("skip functional test")
    }

    host := os.Getenv("TRACKING_HOST")
    if host == "" {
        host = "localhost:8087"
    }

    jsonData := []byte(`{
        "shipment_id":1,
        "tracking_number":"LOG-1-123456",
        "status":"IN_TRANSIT",
        "location":"Jakarta",
        "note":"Paket sedang dikirim"
    }`)

    resp, err := http.Post(
        "http://test-tracking:8087/tracking",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected status 200, got %d", resp.StatusCode)
    }

    var result trackingResponse
    err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
        t.Fatal(err)
    }

    if result.Status != "IN_TRANSIT" {
        t.Errorf("Expected IN_TRANSIT, got %s", result.Status)
    }
}
