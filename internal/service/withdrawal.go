package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type WithdrawalService struct {
	withdrawalRepo domain.WithdrawalRepository
}

func NewWithdrawalService(repository domain.WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{repository}
}

func (w *WithdrawalService) GetAllByLogin(ctx context.Context, login string) ([]*models.Withdrawal, error) {
	return w.withdrawalRepo.GetAllByLogin(ctx, login)
}

func (w *WithdrawalService) Withdraw(ctx context.Context, withdraw *models.Withdrawal) error {
	return w.withdrawalRepo.Withdraw(ctx, withdraw)
}

func (w *WithdrawalService) IsExist(ctx context.Context, withdraw *models.Withdrawal) (bool, error) {
	return w.withdrawalRepo.IsExist(ctx, withdraw)
}
