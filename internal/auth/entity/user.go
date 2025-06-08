package entity

import (
	"github.com/google/uuid"
)

type User struct {
	Username string
	ID       uuid.UUID
}

type LoginRequest struct {
	Username string `json:"login" binding:"required,min=3"`
	Password string `json:"password" binding:"required"`
}
