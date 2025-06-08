package handlers

import (
	"context"
	"net/http"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/orders/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service interface {
	AddOrder(ctx context.Context, orderNumber string, userID uuid.UUID) (*entity.Order, bool, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
}

type OrderHandlers struct {
	service      Service
	errorHandler apperrors.NextHandler
}

func (h *OrderHandlers) AddOrder(c *gin.Context) {
	orderNumber, err := entity.GinToOrderNumber(c)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	userID, err := entity.GinToUserUUID(c)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	_, exists, err := h.service.AddOrder(c.Request.Context(), orderNumber, userID)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	if !exists {
		c.Status(http.StatusAccepted)
		return
	}

	c.Status(http.StatusOK)
}

func (h *OrderHandlers) GetList(c *gin.Context) {
	userID, err := entity.GinToUserUUID(c)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	orders, err := h.service.GetAllByUserID(c.Request.Context(), userID)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, orders)
}
