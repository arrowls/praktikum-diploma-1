package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/arrowls/praktikum-diploma-1/internal/accrual/entity"
	orders "github.com/arrowls/praktikum-diploma-1/internal/orders/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"resty.dev/v3"
)

type Repo interface {
	GetUnprocessedOrders(ctx context.Context) ([]entity.Order, error)
	UpdateBalanceAndOrderStatus(ctx context.Context, userID uuid.UUID, amount float32, orderNumber, status string) error
	UpdateOrderStatus(ctx context.Context, orderNumber, status string) error
}

type AccrualService struct {
	repo            Repo
	ticker          *time.Ticker
	baseInterval    int
	logger          *logrus.Logger
	accrualEndpoint string
}

func (s *AccrualService) Run() {
	for range s.ticker.C {
		s.resetTicker()
		ctx := context.Background()
		orderList := s.extractUnprocessedOrders(ctx)
		s.processOrders(orderList)
	}
}

func (s *AccrualService) processOrders(orders []entity.Order) {
	for _, order := range orders {
		accrualResponse := s.makeRequest(order.Number)
		if accrualResponse == nil {
			continue
		}

		if accrualResponse.Status == entity.StatusProcessed {
			s.processSuccessCase(order, accrualResponse)
		}

		if accrualResponse.Status == entity.StatusInvalid {
			s.processErrorCase(order)
		}
	}
}

func (s *AccrualService) extractUnprocessedOrders(ctx context.Context) []entity.Order {
	orderList, err := s.repo.GetUnprocessedOrders(ctx)
	if err != nil {
		s.logger.Errorf("error getting unprocessed orders")
		return []entity.Order{}
	}

	return orderList
}

func (s *AccrualService) makeRequest(orderNumber string) *entity.AccrualResponse {
	resp, err := resty.New().R().SetResult(&entity.AccrualResponse{}).Get(s.accrualEndpoint + "/api/orders/" + orderNumber)
	if err != nil {
		s.logger.Errorf("request failed: %v", err)
		return nil
	}

	if resp.StatusCode() == http.StatusTooManyRequests {
		retryAfter, err := strconv.Atoi(resp.Header().Get("Retry-After"))
		if err != nil {
			return nil
		}

		s.setNewTicker(time.Duration(retryAfter))
	}

	if status := resp.StatusCode(); status != http.StatusOK {
		s.logger.Infof("got response with status %d %s", status, orderNumber)
		return nil
	}

	result, ok := resp.Result().(*entity.AccrualResponse)
	if !ok {
		s.logger.Errorf("got bad response: %v", resp.Result())
		return nil
	}

	s.logger.Infof("got result %v", result)

	return result
}

func (s *AccrualService) processSuccessCase(order entity.Order, response *entity.AccrualResponse) {
	err := s.repo.UpdateBalanceAndOrderStatus(context.Background(), order.UserID, response.Accrual, order.Number, orders.StatusProcessed)
	if err != nil {
		s.logger.Errorf("error handling succes case: %v", err)
	}
}

func (s *AccrualService) processErrorCase(order entity.Order) {
	err := s.repo.UpdateOrderStatus(context.Background(), order.Number, orders.StatusInvalid)
	if err != nil {
		s.logger.Errorf("error handling error case: %v", err)
	}
}

func (s *AccrualService) setNewTicker(duration time.Duration) {
	s.ticker = time.NewTicker(duration * time.Second)
}

func (s *AccrualService) resetTicker() {
	s.ticker = time.NewTicker(time.Duration(s.baseInterval) * time.Second)
}
