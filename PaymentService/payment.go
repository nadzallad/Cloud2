package main

func ValidatePayment(amount int, paid int) string {
	if paid == 0 {
		return "PENDING"
	}
	if paid < amount {
		return "FAILED"
	}
	return "PAID"
}