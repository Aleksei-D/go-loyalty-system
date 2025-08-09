package domain

import (
	"context"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetByLogin(ctx context.Context, username string) (*models.User, bool)
}
