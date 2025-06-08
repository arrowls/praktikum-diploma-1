package entity

import (
	"time"
)

type Withdrawal struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type WithdrawRequest struct {
	Order string `json:"order" binding:"required"`
	Sum   int    `json:"sum" binding:"required"`
}
