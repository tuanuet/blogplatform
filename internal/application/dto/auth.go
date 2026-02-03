package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100" binding:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email" binding:"required,email"`
	Password string `json:"password" validate:"required,min=8" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" binding:"required,email"`
	Password string `json:"password" validate:"required" binding:"required"`
}

type LoginWithSocialRequest struct {
	Provider string `json:"provider" validate:"required" binding:"required"`
	Code     string `json:"code" validate:"required" binding:"required"`
}

type AuthResponse struct {
	SessionID string    `json:"sessionId"`
	UserID    uuid.UUID `json:"userId"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
}
