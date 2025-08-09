package service

import "github.com/Aleksei-D/go-loyalty-system/internal/domain"

type Service struct {
	BalanceService    *BalanceService
	UserService       *UserService
	OrderService      *OrderService
	WithdrawalService *WithdrawalService
}

func NewService(balanceRepo domain.BalanceRepository, orderRepo domain.OrderRepository, userRepo domain.UserRepository, withdrawalRepo domain.WithdrawalRepository) *Service {
	return &Service{
		BalanceService:    NewBalanceService(balanceRepo),
		UserService:       NewUserService(userRepo),
		OrderService:      NewOrderService(orderRepo),
		WithdrawalService: NewWithdrawalService(withdrawalRepo),
	}
}
