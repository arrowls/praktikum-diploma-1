package repository

import (
	"context"
	"fmt"

	"github.com/arrowls/praktikum-diploma-1/internal/accrual/entity"
	orders "github.com/arrowls/praktikum-diploma-1/internal/orders/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

type AccrualRepository struct {
	db     *pgx.Conn
	logger *logrus.Logger
}

func (r *AccrualRepository) UpdateBalanceAndOrderStatus(ctx context.Context, userID uuid.UUID, amount float32, orderNumber, status string) error {
	err := r.initUserBalance(ctx, userID)
	if err != nil {
		return fmt.Errorf("error initializing user (%s) balance: %w", userID.String(), err)
	}

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

	_, err = tx.Exec(ctx, `
		update user_balance set current = current + $1 where user_id = $2
	`, amount, userID)
	if err != nil {
		return fmt.Errorf("error updating user (%s) balance: %w", userID.String(), err)
	}

	_, err = tx.Exec(ctx, `
		update orders set status = $1, accrual = $3 where number = $2
	`, status, orderNumber, amount)
	if err != nil {
		return fmt.Errorf("error updating order (%s) status: %w", orderNumber, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("error commiting tx: %v", err)
	}

	return nil
}

func (r *AccrualRepository) UpdateOrderStatus(ctx context.Context, orderNumber, status string) error {
	_, err := r.db.Exec(ctx, `
		update orders set status = $1 where number = $2
	`, status, orderNumber)

	return err
}

func (r *AccrualRepository) GetUnprocessedOrders(ctx context.Context) ([]entity.Order, error) {
	rows, err := r.db.Query(ctx, `
		select number, user_id  from orders where status in ($1, $2)
	`, orders.StatusNew, orders.StatusProcessing)
	if err != nil {
		return nil, fmt.Errorf("could not get rows for orders: %v", err)
	}

	list := make([]entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.Number, &order.UserID); err != nil {
			r.logger.Errorf("unexpected error while retrieving order number: %v", err)
			continue
		}

		list = append(list, order)
	}

	return list, nil
}

func (r *AccrualRepository) initUserBalance(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
			insert into user_balance (user_id, current, withdrawn) values ($1, 0, 0) on conflict do nothing
		`, userID)
	if err != nil {
		return fmt.Errorf("error creating emty balance for user %s %w", userID.String(), err)
	}

	return nil
}
