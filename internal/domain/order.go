package domain

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type OrderRepository interface {
	Add(ctx context.Context, orderNumber, login string) (*models.Order, error)
	GetAllByLogin(ctx context.Context, login string) ([]*models.Order, error)
	GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, bool)
	GetOrderByLoginAndNumber(ctx context.Context, login, orderNumber string) (*models.Order, bool)
	GetNotAcceptedOrderNumbers(ctx context.Context, limit uint) ([]*models.Order, error)
	UpdateStatus(ctx context.Context, order *models.Order) error
}
