package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type OrderService struct {
	orderRepo domain.OrderRepository
}

func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{orderRepo: repo}
}

func (o *OrderService) GetAllByLogin(ctx context.Context, login string) ([]*models.Order, error) {
	return o.orderRepo.GetAllByLogin(ctx, login)
}

func (o *OrderService) Add(ctx context.Context, login, orderNumber string) (*models.Order, error) {
	return o.orderRepo.Add(ctx, login, orderNumber)
}

func (o *OrderService) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	return o.orderRepo.GetOrderByNumber(ctx, orderNumber)
}

func (o *OrderService) GetNotAcceptedOrderNumbers(ctx context.Context, limit, updateTimeout uint) ([]*models.Order, error) {
	return o.orderRepo.GetNotAcceptedOrderNumbers(ctx, limit, updateTimeout)
}

func (o *OrderService) UpdateStatus(ctx context.Context, order *models.Order) error {
	return o.orderRepo.UpdateStatus(ctx, order)
}

func (o *OrderService) IsExist(ctx context.Context, orderNumber string) (bool, error) {
	return o.orderRepo.IsExist(ctx, orderNumber)
}
