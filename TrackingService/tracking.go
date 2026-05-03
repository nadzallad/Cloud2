package main

func ValidateTracking(status string) string {
	switch status {
	case "PICKED_UP", "IN_TRANSIT", "DELIVERED":
		return status
	default:
		return "UNKNOWN"
	}
}