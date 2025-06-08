package entity

import (
	"fmt"
	"strings"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GinToOrderNumber(c *gin.Context) (string, error) {
	body, err := c.GetRawData()
	if err != nil {
		return "", fmt.Errorf("error reading request body %w", apperrors.ErrBadRequest)
	}

	if string(body) == "" {
		return "", fmt.Errorf("empty body %w", apperrors.ErrBadRequest)
	}

	if strings.Trim(string(body), "0123456789") == "" {
		return string(body), nil
	}

	return "", fmt.Errorf("request contains invalid digits %w", apperrors.ErrBadRequest)
}

func GinToUserUUID(c *gin.Context) (uuid.UUID, error) {
	user, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, apperrors.ErrUnauthorized
	}

	userID, ok := user.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperrors.ErrUnauthorized
	}

	return userID, nil
}
