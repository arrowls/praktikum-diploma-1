package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
	UserID     uuid.UUID `json:"-"`
}
