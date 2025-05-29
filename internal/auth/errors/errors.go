package autherrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/gin-gonic/gin"
)

var (
	ErrUsernameExists = errors.New("username exists")
)

type AuthErrorHandler struct {
	NextHandler apperrors.NextHandler
}

func (h *AuthErrorHandler) Handle(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var status int
	response := apperrors.DefaultErrorResponse{}

	switch {
	case errors.Is(err, ErrUsernameExists):
		status = http.StatusConflict
		response.Key = ErrUsernameExists.Error()
	default:
		h.NextHandler(c, err)
		return
	}
	log.Println(err.Error())

	c.JSON(status, response)
}
