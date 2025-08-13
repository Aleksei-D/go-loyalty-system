package domain

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type WithdrawalRepository interface {
	GetAllByLogin(ctx context.Context, login string) ([]*models.Withdrawal, error)
	Withdraw(ctx context.Context, withdraw *models.Withdrawal) error
	IsExist(ctx context.Context, withdraw *models.Withdrawal) (bool, error)
}
