package entity

// MonthlyCount represents a count for a specific month
type MonthlyCount struct {
	Month string `json:"month"` // Format: "YYYY-MM"
	Count int64  `json:"count"`
}
