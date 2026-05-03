package main

func CheckWarehouse(stock int) string {
	if stock > 0 {
		return "available"
	}
	return "out_of_stock"
}