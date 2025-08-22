package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
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

func (o *OrderService) GetNotAcceptedOrderNumbers(ctx context.Context, limit, updateTimeout uint) ([]*models.Order, error) {
	return o.orderRepo.GetNotAcceptedOrderNumbers(ctx, limit, updateTimeout)
}

func (o *OrderService) UpdateStatus(ctx context.Context, order *models.Order) error {
	return o.orderRepo.UpdateStatus(ctx, order)
}

func (o *OrderService) AddOrder(ctx context.Context, orderNumber, login string) (*models.Order, error) {
	if ok := common.CheckLuhnAlgorithm(orderNumber); !ok {
		return nil, common.ErrInvalidOrderNumber
	}

	ok, err := o.orderRepo.IsExist(ctx, orderNumber)
	if err != nil {
		return nil, err
	}

	if ok {
		existOrder, err := o.orderRepo.GetOrderByNumber(ctx, orderNumber)
		if err != nil {
			return nil, err
		}

		if existOrder.Login == login {
			return nil, common.ErrOrderAlreadyAdded
		}
		return nil, common.ErrOrderBelongAnotherUser
	}

	addedOrder, err := o.orderRepo.Add(ctx, login, orderNumber)
	if err != nil {
		return nil, err
	}
	return addedOrder, err
}
