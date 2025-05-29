package service

import (
	"context"
	"fmt"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/auth/entity"
	"github.com/arrowls/praktikum-diploma-1/internal/token"
)

type AuthRepo interface {
	AddUser(ctx context.Context, username, password string) (*entity.User, error)
	GetUser(ctx context.Context, username, password string) (*entity.User, error)
}

type AuthService struct {
	repo AuthRepo
}

func (s *AuthService) generateTokens(user *entity.User) (*entity.TokenResponse, error) {
	accessToken, err := token.New(user)
	if err != nil {
		return nil, fmt.Errorf("error generating access token: %v %w", err, apperrors.ErrUnknown)
	}

	return &entity.TokenResponse{AccessToken: accessToken}, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*entity.TokenResponse, error) {
	user, err := s.repo.GetUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("error validating user: %w", err)
	}

	return s.generateTokens(user)
}

func (s *AuthService) Register(ctx context.Context, username, password string) (*entity.TokenResponse, error) {
	user, err := s.repo.AddUser(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return s.generateTokens(user)
}
