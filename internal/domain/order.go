package domain

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type OrderRepository interface {
	Add(ctx context.Context, login, orderNumber string) (*models.Order, error)
	GetAllByLogin(ctx context.Context, login string) ([]*models.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error)
	GetNotAcceptedOrderNumbers(ctx context.Context, limit, updateTimeout uint) ([]*models.Order, error)
	UpdateStatus(ctx context.Context, order *models.Order) error
	IsExist(ctx context.Context, orderNumber string) (bool, error)
}
