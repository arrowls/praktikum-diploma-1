package balanceerrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/gin-gonic/gin"
)

var (
	ErrNotEnoughCurrency   = errors.New("not enough currency")
	ErrOrderNumberNotFound = errors.New("order number not found")
	ErrWithdrawalExists    = errors.New("withdrawal exists")
)

type BalanceErrorHandler struct {
	NextHandler apperrors.NextHandler
}

func (h *BalanceErrorHandler) Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var status int
	response := apperrors.DefaultErrorResponse{}

	switch {
	case errors.Is(err, ErrNotEnoughCurrency):
		status = http.StatusPaymentRequired
		response.Key = ErrNotEnoughCurrency.Error()
	case errors.Is(err, ErrOrderNumberNotFound):
		status = http.StatusUnprocessableEntity
		response.Key = ErrOrderNumberNotFound.Error()
	case errors.Is(err, ErrWithdrawalExists):
		status = http.StatusUnprocessableEntity
		response.Key = ErrWithdrawalExists.Error()
	default:
		h.NextHandler(c, err)
		return
	}
	log.Println(err.Error())

	c.JSON(status, response)
}
