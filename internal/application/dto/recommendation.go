package dto

// UpdateUserInterestsRequest represents the request to update user interests
type UpdateUserInterestsRequest struct {
	TagIDs []string `json:"tagIds" binding:"required"`
}
