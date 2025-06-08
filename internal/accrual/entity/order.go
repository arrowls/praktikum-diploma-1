package entity

import (
	"github.com/google/uuid"
)

type Order struct {
	UserID uuid.UUID
	Number string
}
