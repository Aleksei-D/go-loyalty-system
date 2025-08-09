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

type PostgresBalanceRepository struct {
	db *sql.DB
}

func NewPostgresBalanceRepository(db *sql.DB) *PostgresBalanceRepository {
	return &PostgresBalanceRepository{db: db}
}

func (p *PostgresBalanceRepository) Get(ctx context.Context, login string) (*models.Balance, error) {
	var balance models.Balance
	var current float64
	var withdrawn float64
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx, "SELECT current FROM balance WHERE login = $1", login).Scan(&current)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("balalance for user - %s not found", login), zap.Error(err))
			return &balance, nil
		}
		return &balance, err
	}

	balance.Login = login
	balance.Current = current

	err = tx.QueryRowContext(ctx, "SELECT login, SUM(sum) as withdrawns FROM withdrawals WHERE login = $1 group by login", login).Scan(&withdrawn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("balalance for user - %s not found", login), zap.Error(err))
			return &balance, nil
		}
		return &balance, err
	}

	balance.Withdrawn = withdrawn
	return &balance, nil
}
