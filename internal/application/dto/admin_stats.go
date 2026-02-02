package dto

type MonthlyStat struct {
	Month       string `json:"month"` // Format: "YYYY-MM"
	NewUsers    int64  `json:"new_users"`
	NewBlogs    int64  `json:"new_blogs"`
	NewComments int64  `json:"new_comments"`
}

type DashboardStatsResponse struct {
	Stats []MonthlyStat `json:"stats"`
}
