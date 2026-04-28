package main

type Delivery struct {
	ID      string `json:"id"`
	Courier string `json:"courier"`
	Status  string `json:"status"`
}

func CreateDelivery(id string, courier string) Delivery {
	return Delivery{
		ID:      id,
		Courier: courier,
		Status:  "On Process",
	}
}

func UpdateDeliveryStatus(d Delivery, status string) Delivery {
	d.Status = status
	return d
}