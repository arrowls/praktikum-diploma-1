package ordererrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/gin-gonic/gin"
)

var (
	ErrLuhnInvalid           = errors.New("wrong number format")
	ErrOrderAddedByOtherUser = errors.New("order was already added by other user")
)

type OrderErrorHandler struct {
	NextHandler apperrors.NextHandler
}

func (h *OrderErrorHandler) Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var status int
	response := apperrors.DefaultErrorResponse{}

	switch {
	case errors.Is(err, ErrLuhnInvalid):
		status = http.StatusUnprocessableEntity
		response.Key = ErrLuhnInvalid.Error()
	case errors.Is(err, ErrOrderAddedByOtherUser):
		status = http.StatusConflict
		response.Key = ErrOrderAddedByOtherUser.Error()
	default:
		h.NextHandler(c, err)
		return
	}
	log.Println(err.Error())

	c.JSON(status, response)
}
