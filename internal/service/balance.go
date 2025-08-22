package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type BalanceService struct {
	balanceRepo domain.BalanceRepository
}

func NewBalanceService(repo domain.BalanceRepository) *BalanceService {
	return &BalanceService{balanceRepo: repo}
}

func (b *BalanceService) Get(ctx context.Context, login string) (*models.Balance, error) {
	return b.balanceRepo.Get(ctx, login)
}
