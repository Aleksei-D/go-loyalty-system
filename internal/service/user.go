package service

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/domain"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
	crypto2 "github.com/Aleksei-D/go-loyalty-system/internal/utils/crypto"
	"go.uber.org/zap"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) *UserService {
	return &UserService{userRepository}
}

func (u *UserService) Login(ctx context.Context, user *models.User) error {
	existUser, err := u.userRepo.GetByLogin(ctx, user.Login)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return err
	}
	if existUser == nil || !crypto2.CheckPasswordHash(user.Password, existUser.Password) {
		return common.ErrInvalidCredentials
	}

	return nil
}

func (u *UserService) CreateUser(ctx context.Context, newUser *models.User) (*models.User, error) {
	ok, err := u.userRepo.IsExist(ctx, newUser.Login)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return nil, err
	}

	if ok {
		return nil, common.ErrUserAlreadyExists
	}

	hashPassword, err := crypto2.HashPassword(newUser.Password)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return nil, err
	}

	newUser.Password = hashPassword
	user, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		logger.Log.Warn("User Create Error", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetToken(login, secretKey string) (string, error) {
	var tokenString string
	tokenString, err := crypto2.CreateToken(login, secretKey)
	if err != nil {
		logger.Log.Warn(err.Error(), zap.Error(err))
		return tokenString, err
	}
	return tokenString, nil
}
