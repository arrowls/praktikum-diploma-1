package entity

import (
	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
