package apperrors

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
	ErrUnknown      = errors.New("internal server error")
)

type DefaultErrorResponse struct {
	Key string `json:"key"`
}

func DefaultErrorHandler(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var status int
	response := DefaultErrorResponse{}

	log.Println(err.Error())

	switch {
	case errors.Is(err, ErrUnauthorized):
		status = http.StatusUnauthorized
		response.Key = ErrUnauthorized.Error()
	case errors.Is(err, ErrBadRequest):
		status = http.StatusBadRequest
		response.Key = ErrBadRequest.Error()
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
		response.Key = ErrNotFound.Error()
	default:
		status = http.StatusInternalServerError
		response.Key = ErrUnknown.Error()
	}

	c.JSON(status, response)
}

type NextHandler func(c *gin.Context, err error)
