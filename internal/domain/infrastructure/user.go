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

	userCreateQuery := `INSERT INTO user (login, password) VALUES ($1, $2)`
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

	_, err = stmtBalance.ExecContext(ctx, user.Login, user.Password)
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

func (p *PostgresUserRepository) GetByLogin(ctx context.Context, username string) (*models.User, bool) {
	var user *models.User
	var login string
	var password string
	row := p.db.QueryRowContext(ctx, "SELECT login, password FROM user WHERE login = $1", username)

	if err := row.Err(); err != nil {
		if errors.Is(sql.ErrNoRows, row.Err()) {
			logger.Log.Info(fmt.Sprintf("user - %s not found", username), zap.Error(err))
			return user, false
		}
		return user, false
	}

	err := row.Scan(&login, &password)
	if err != nil {
		return user, false
	}
	user.Login = &login
	user.Password = &password
	return user, true
}
