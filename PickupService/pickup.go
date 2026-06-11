package main

func ProcessPickup(paymentStatus string, weight int) string {

	if paymentStatus == "paid" && weight > 0 {
		return "scheduled"
	}

	return "pending"
}
