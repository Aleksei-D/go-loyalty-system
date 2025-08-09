package domain

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByLogin(ctx context.Context, login string) (*models.User, error)
	IsExist(ctx context.Context, login string) (bool, error)
}
