package main

// Function untuk hitung ongkir
func CalculateShippingCost(weightKg float64, distanceKm float64) float64 {
	const pricePerKg = 5000
	const pricePerKm = 2000

	return (weightKg * pricePerKg) + (distanceKm * pricePerKm)
}

// Function untuk hitung total harga
func CalculateTotalPrice(basePrice float64, shippingCost float64) float64 {
	return basePrice + shippingCost
}