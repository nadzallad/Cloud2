package main

// ValidateNotification menentukan status notifikasi berdasarkan pembayaran
func ValidateNotification(amount int, paid int) string {
	if paid >= amount {
		return "paid"
	}
	return "pending"
}