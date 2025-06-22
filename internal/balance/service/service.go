package service

import (
	"context"

	"github.com/arrowls/praktikum-diploma-1/internal/balance/entity"
	"github.com/google/uuid"
)

type BalanceRepo interface {
	GetWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Withdrawal, error)
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*entity.Balance, error)
	Withdraw(ctx context.Context, userID uuid.UUID, req *entity.WithdrawRequest) error
}

type BalanceService struct {
	repo BalanceRepo
}

func (s *BalanceService) GetWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Withdrawal, error) {
	return s.repo.GetWithdrawalsByUserID(ctx, userID)
}

func (s *BalanceService) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*entity.Balance, error) {
	return s.repo.GetBalanceByUserID(ctx, userID)
}

func (s *BalanceService) Withdraw(ctx context.Context, userID uuid.UUID, req *entity.WithdrawRequest) error {
	return s.repo.Withdraw(ctx, userID, req)
}
