package domain

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type BalanceRepository interface {
	Get(ctx context.Context, login string) (*models.Balance, error)
}
