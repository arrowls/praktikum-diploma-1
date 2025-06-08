package handlers

import (
	"net/http"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/balance/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type Service interface {
	GetWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Withdrawal, error)
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*entity.Balance, error)
	Withdraw(ctx context.Context, userID uuid.UUID, req *entity.WithdrawRequest) error
}

type BalanceHandlers struct {
	service      Service
	errorHandler apperrors.NextHandler
}

func (h *BalanceHandlers) GetBalance(c *gin.Context) {
	userID, err := entity.GinToUserUUID(c)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	balance, err := h.service.GetBalanceByUserID(c.Request.Context(), userID)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (h *BalanceHandlers) GetWithdrawals(c *gin.Context) {
	userID, err := entity.GinToUserUUID(c)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	withdrawals, err := h.service.GetWithdrawalsByUserID(c.Request.Context(), userID)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	if len(withdrawals) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, withdrawals)
}

func (h *BalanceHandlers) Withdraw(c *gin.Context) {
	var req entity.WithdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorHandler(c, apperrors.ErrBadRequest)
		return
	}

	userID, err := entity.GinToUserUUID(c)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	err = h.service.Withdraw(c.Request.Context(), userID, &req)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.Status(http.StatusOK)
}
