package common

import "time"

// User model
type User struct {
	UserID               int
	LastCardUsePlace     string
	LastCardUseTimestamp time.Time
}

// Airport model
type Airport struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
