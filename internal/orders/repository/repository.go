package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/arrowls/praktikum-diploma-1/internal/apperrors"
	"github.com/arrowls/praktikum-diploma-1/internal/orders/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

const orderDefaultStatus = entity.StatusNew

type OrdersRepository struct {
	logger *logrus.Logger
	db     *pgx.Conn
}

func (r *OrdersRepository) AddOrder(ctx context.Context, number string, userID uuid.UUID) (*entity.Order, error) {
	_, err := r.db.Exec(ctx, `
		insert into orders (number, user_id, status) values ($1, $2, $3)
	`, number, userID, orderDefaultStatus)

	if err != nil {
		return nil, fmt.Errorf("%v %w", err, apperrors.ErrBadRequest)
	}
	return r.GetOrderByNumber(ctx, number)
}

func (r *OrdersRepository) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Order, error) {
	rows, err := r.db.Query(ctx, `
		select number, status, accrual::decimal, uploaded_at, user_id from orders where user_id = $1 order by uploaded_at desc 
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("could not get rows for user %s: %v", userID.String(), err)
	}

	list := make([]entity.Order, 0)
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID); err != nil {
			r.logger.Errorf("unexpected error while retrieving orders list for id %s: %v", userID.String(), err)
			continue
		}

		list = append(list, order)
	}

	return list, nil
}

func (r *OrdersRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*entity.Order, error) {
	var order entity.Order

	err := r.db.QueryRow(ctx, `
		select number, status, accrual::decimal, uploaded_at, user_id from orders where number = $1
	`, orderNumber).Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt, &order.UserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%v %w", err, apperrors.ErrNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("error retrieving order with number %s %v", orderNumber, err)
	}

	return &order, nil
}
