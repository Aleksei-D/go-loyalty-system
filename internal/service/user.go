package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository}
}

func (u *UserService) Create(ctx context.Context, user *models.User) (*models.User, error) {
	return u.userRepo.Create(ctx, user)
}

func (u *UserService) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	return u.userRepo.GetByLogin(ctx, login)
}

func (u *UserService) IsExist(ctx context.Context, login string) (bool, error) {
	return u.userRepo.IsExist(ctx, login)
}
