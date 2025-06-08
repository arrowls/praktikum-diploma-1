package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/auth/entity"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*entity.TokenResponse, error)
	Register(ctx context.Context, username, password string) (*entity.TokenResponse, error)
}

type AuthHandlers struct {
	service      Service
	errorHandler apperrors.NextHandler
}

func (h *AuthHandlers) Register(c *gin.Context) {
	var req entity.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.errorHandler(c, fmt.Errorf("error validating register request: %v %w", err, apperrors.ErrBadRequest))
		return
	}

	token, err := h.service.Register(c.Request.Context(), req.Username, req.Password)

	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	c.JSON(http.StatusOK, token)
}

func (h *AuthHandlers) Login(c *gin.Context) {
	var req entity.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		h.errorHandler(c, fmt.Errorf("error validating login request: %v %w", err, apperrors.ErrBadRequest))
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	c.JSON(http.StatusOK, token)
}
