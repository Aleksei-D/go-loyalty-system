package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"go.uber.org/zap"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (p *PostgresUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	userCreateQuery := `INSERT INTO users (login, password) VALUES ($1, $2)`
	stmtUser, err := tx.Prepare(userCreateQuery)
	if err != nil {
		return nil, err
	}
	defer stmtUser.Close()

	_, err = stmtUser.ExecContext(ctx, user.Login, user.Password)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	balanceCreateQuery := `INSERT INTO balance (login) VALUES ($1)`
	stmtBalance, err := tx.Prepare(balanceCreateQuery)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	_, err = stmtBalance.ExecContext(ctx, user.Login)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {

		return nil, err
	}
	return user, nil
}

func (p *PostgresUserRepository) GetByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	var loginFromDb string
	var password string
	err := p.db.QueryRowContext(ctx, "SELECT login, password FROM users WHERE login = $1", login).Scan(&loginFromDb, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("user - %s not found", login), zap.Error(err))
			return &user, nil
		}
		return &user, err
	}

	user.Login = loginFromDb
	user.Password = password
	return &user, err
}

func (p *PostgresUserRepository) IsExist(ctx context.Context, login string) (bool, error) {
	var loginFromDb string
	var password string
	err := p.db.QueryRowContext(ctx, "SELECT login, password FROM users WHERE login = $1", login).Scan(&loginFromDb, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("user - %s not found", login), zap.Error(err))
			return false, nil
		}
		return false, err
	}

	return true, nil
}
