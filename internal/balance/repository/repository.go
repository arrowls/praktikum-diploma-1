package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arrowls/praktikum-diploma-1/internal/balance/entity"
	balanceerrors "github.com/arrowls/praktikum-diploma-1/internal/balance/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type BalanceRepository struct {
	logger *logrus.Logger
	db     *pgx.Conn
}

func (r *BalanceRepository) initUserBalance(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
			insert into user_balance (user_id, current, withdrawn) values ($1, 0, 0)
		`, userID)
	if err != nil {
		return fmt.Errorf("error creating emty balance for user %s %w", userID.String(), err)
	}

	return nil
}

func (r *BalanceRepository) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*entity.Balance, error) {
	var balance entity.Balance
	err := r.db.QueryRow(ctx, `
		select current, withdrawn from user_balance where user_id = $1
	`, userID).Scan(&balance.Current, &balance.Withdrawn)

	if errors.Is(err, pgx.ErrNoRows) {
		err := r.initUserBalance(ctx, userID)
		if err != nil {
			return nil, err
		}
		return r.GetBalanceByUserID(ctx, userID)
	}

	if err != nil {
		return nil, fmt.Errorf("error retrieving balance for user %s: %w", userID.String(), err)
	}

	return &balance, nil
}

func (r *BalanceRepository) GetWithdrawalsByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Withdrawal, error) {
	rows, err := r.db.Query(ctx, `
		select order_number, sum, processed_at from withdrawals where user_id = $1;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving withdrawals for user %s: %w", userID.String(), err)
	}

	list := make([]entity.Withdrawal, 0)
	for rows.Next() {
		var withdrawal entity.Withdrawal
		err := rows.Scan(&withdrawal.Order, &withdrawal.Sum, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning for withdrawals for user %s: %w", userID.String(), err)
		}

		list = append(list, withdrawal)
	}

	return list, nil
}

func (r *BalanceRepository) Withdraw(ctx context.Context, userID uuid.UUID, req *entity.WithdrawRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error begin tx: %w", err)
	}

	defer func() {
		errRollback := tx.Rollback(ctx)
		if err != nil {
			r.logger.Errorf("error during tx rollback: %v", errRollback)
		}
	}()

	ok, err := r.validateWithdrawalIsNew(ctx, tx, req.Order)
	if err != nil {
		return fmt.Errorf("error checking withdrawal: %v", err)
	}
	if !ok {
		return balanceerrors.ErrWithdrawalExists
	}

	var currentBalance float32
	err = tx.QueryRow(ctx, `
		select current from user_balance where user_id = $1 for update 
	`, userID).Scan(&currentBalance)

	if errors.Is(err, pgx.ErrNoRows) {
		err := r.initUserBalance(ctx, userID)
		if err != nil {
			return err
		}

		return r.Withdraw(ctx, userID, req)
	}

	if err != nil {
		return fmt.Errorf("error reading user balance: %v", err)
	}

	if currentBalance < req.Sum {
		return balanceerrors.ErrNotEnoughCurrency
	}

	res, err := tx.Exec(ctx, `
		update user_balance set current = current - $1 where user_id = $2 and current >= $1;
	`, req.Sum, userID)
	if err != nil {
		return fmt.Errorf("error updating user (%s) balance: %v", userID.String(), err)
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("no balance records found for user %s", userID.String())
	}

	_, err = tx.Exec(ctx, `
		update user_balance set withdrawn = user_balance.withdrawn + $1 where user_id = $2
	`, req.Sum, userID)
	if err != nil {
		return fmt.Errorf("error updating user (%s) withdrawals: %v", userID.String(), err)
	}

	_, err = tx.Exec(ctx, `
		insert into withdrawals (user_id, order_number, sum) values ($1, $2, $3)
	`, userID, req.Order, req.Sum)

	if err != nil && strings.Contains(err.Error(), "violates foreign key constraint") {
		return balanceerrors.ErrOrderNumberNotFound
	}

	if err != nil {
		return fmt.Errorf("error creating a withdrawal record: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting tx: %v", err)
	}
	return nil
}

func (r *BalanceRepository) validateWithdrawalIsNew(ctx context.Context, tx pgx.Tx, order string) (bool, error) {
	var exists bool
	err := tx.QueryRow(ctx, `
		select exists (select 1 from withdrawals where order_number = $1)
	`, order).Scan(&exists)
	if err != nil {
		return false, err
	}

	return !exists, nil
}
