package entity

import (
	"time"
)

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         float32   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawRequest struct {
	Order string  `json:"order" binding:"required"`
	Sum   float32 `json:"sum" binding:"required"`
}
