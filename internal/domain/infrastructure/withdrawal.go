package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/Aleksei-D/go-loyalty-system/internal/utils/common"
	"go.uber.org/zap"
	"time"
)

type PostgresWithdrawalRepository struct {
	db *sql.DB
}

func NewPostgresWithdrawalRepository(db *sql.DB) *PostgresWithdrawalRepository {
	return &PostgresWithdrawalRepository{db: db}
}

func (p *PostgresWithdrawalRepository) GetAllByLogin(ctx context.Context, login string) ([]*models.Withdrawal, error) {
	withdrawals := make([]*models.Withdrawal, 0)
	query := `SELECT order_number, sum, processed_at FROM withdrawals WHERE login = $1 ORDER BY processed_at DESC`
	stmt, err := p.db.Prepare(query)
	if err != nil {
		return withdrawals, err
	}

	rows, err := stmt.QueryContext(ctx, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return withdrawals, nil
		}
		logger.Log.Info(err.Error(), zap.Error(err))
		return withdrawals, err
	}
	if rows.Err() != nil {
		return withdrawals, rows.Err()
	}
	defer rows.Close()

	for rows.Next() {
		var processedAt time.Time
		var order string
		var sum float64
		var withdrawal models.Withdrawal
		err := rows.Scan(&order, &sum, &processedAt)
		if err != nil {
			return withdrawals, err
		}

		withdrawal.Login = login
		withdrawal.OrderNumber = order
		withdrawal.ProcessedAt = models.CustomTime{Time: processedAt}
		withdrawal.Sum = sum
		withdrawals = append(withdrawals, &withdrawal)
	}
	return withdrawals, nil
}

func (p *PostgresWithdrawalRepository) Withdraw(ctx context.Context, withdraw *models.Withdrawal) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var currentBalance float64
	err = tx.QueryRowContext(ctx, "SELECT current FROM balance WHERE login = $1", withdraw.Login).Scan(&currentBalance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("user balance - %s not found", withdraw.Login), zap.Error(err))
			return err
		}
		return err
	}

	if currentBalance < withdraw.Sum {
		return common.ErrPaymentInsufficient
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO withdrawals (login, order_number, sum) VALUES ($1, $2, $3)",
		withdraw.Login,
		withdraw.OrderNumber,
		withdraw.Sum,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE balance SET current = current - $1 WHERE login = $2",
		withdraw.Sum,
		withdraw.Login,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresWithdrawalRepository) IsExist(ctx context.Context, withdraw *models.Withdrawal) (bool, error) {
	var orderNumber string
	err := p.db.QueryRowContext(ctx, "SELECT order_number FROM withdrawals WHERE order_number = $1", withdraw.OrderNumber).Scan(&orderNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("order - %s already used", orderNumber), zap.Error(err))
			return false, nil
		}
		return false, err
	}

	return true, nil
}
