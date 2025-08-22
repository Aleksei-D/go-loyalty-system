package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
)

type WithdrawalService struct {
	withdrawalRepo domain.WithdrawalRepository
}

func NewWithdrawalService(repository domain.WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{repository}
}

func (w *WithdrawalService) GetAllByLogin(ctx context.Context, login string) ([]*models.Withdrawal, error) {
	withdrawals, err := w.withdrawalRepo.GetAllByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if len(withdrawals) == 0 {
		return nil, common.ErrNoContent
	}

	return w.withdrawalRepo.GetAllByLogin(ctx, login)
}

func (w *WithdrawalService) Withdraw(ctx context.Context, withdrawal *models.Withdrawal) error {
	if ok := common.CheckLuhnAlgorithm(withdrawal.OrderNumber); !ok {
		return common.ErrInvalidOrderNumber
	}
	ok, err := w.withdrawalRepo.IsExist(ctx, withdrawal)
	if err != nil {
		return err
	}
	if ok {
		return common.ErrOrderAlreadyAdded
	}

	return w.withdrawalRepo.Withdraw(ctx, withdrawal)
}
