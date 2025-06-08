package service

import (
	"context"
	"errors"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/orders/entity"
	ordererrors "github.com/arrowls/praktikum-diploma-1/internal/orders/errors"
	"github.com/google/uuid"
)

type OrderRepo interface {
	AddOrder(ctx context.Context, number string, userID uuid.UUID) (*entity.Order, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.Order, error)
}

type OrderService struct {
	repo OrderRepo
}

func (s *OrderService) AddOrder(ctx context.Context, orderNumber string, userID uuid.UUID) (*entity.Order, error) {
	isValidLuhn := s.validateLuhn(orderNumber)
	if !isValidLuhn {
		return nil, ordererrors.ErrLuhnInvalid
	}

	existingOrder, err := s.repo.GetOrderByNumber(ctx, orderNumber)
	if errors.Is(err, apperrors.ErrNotFound) {
		order, err := s.repo.AddOrder(ctx, orderNumber, userID)
		if err != nil {
			return nil, err
		}

		return order, nil
	}

	if err != nil {
		return nil, err
	}

	if existingOrder.UserID.String() != userID.String() {
		return nil, ordererrors.ErrOrderAddedByOtherUser
	}

	return existingOrder, nil

}

func (s *OrderService) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	return s.repo.GetAllByUserID(ctx, userID)
}

func (s *OrderService) validateLuhn(orderNumber string) bool {
	sum := 0
	nDigits := len(orderNumber)
	parity := nDigits % 2

	for i := 0; i < nDigits; i++ {
		digit := int(orderNumber[i] - '0')
		if digit < 0 || digit > 9 {
			return false
		}

		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
