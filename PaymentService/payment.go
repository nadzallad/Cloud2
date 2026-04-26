package main

func ValidatePayment(amount int, paid int) string {
	if paid >= amount {
		return "paid"
	}
	return "pending"
}