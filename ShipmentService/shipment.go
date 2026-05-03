package main

type Shipment struct {
	ID     string
	Status string
}

func CreateShipment(id string) Shipment {
	return Shipment{
		ID: id,
		Status: "Pending",
	}
}

func UpdateShipmentStatus(s Shipment, status string) Shipment {
	s.Status = status
	return s
}