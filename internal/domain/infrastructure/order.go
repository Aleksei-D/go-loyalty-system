package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Aleksei-D/go-loyalty-system/internal/logger"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (p *PostgresOrderRepository) Add(ctx context.Context, login, orderNumber string) (*models.Order, error) {
	var order models.Order
	var status string
	var accrual sql.NullFloat64
	var loginFromDB string
	var orderNumberFromDB string
	var uploadedAt time.Time

	row := p.db.QueryRowContext(
		ctx,
		"INSERT INTO orders (number, login) VALUES ($1, $2) RETURNING login, number, status, accrual, uploaded_at",
		orderNumber,
		login,
	)
	if err := row.Err(); err != nil {
		logger.Log.Info(err.Error(), zap.Error(err))
		return &order, err
	}

	err := row.Scan(&loginFromDB, &orderNumberFromDB, &status, &accrual, &uploadedAt)
	if err != nil {
		logger.Log.Info(err.Error(), zap.Error(err))
		return &order, err
	}

	if accrual.Valid {
		order.Accrual = &accrual.Float64
	}
	order.Login = loginFromDB
	order.Status = status
	order.Number = orderNumberFromDB
	order.UploadedAt = models.CustomTime{Time: uploadedAt}
	return &order, nil
}

func (p *PostgresOrderRepository) GetAllByLogin(ctx context.Context, login string) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)
	rows, err := p.db.QueryContext(
		ctx,
		"SELECT status, number, accrual, uploaded_at FROM orders WHERE login = $1 ORDER BY uploaded_at DESC",
		login,
	)
	if rows.Err() != nil {
		return orders, err
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(err.Error(), zap.Error(err))
			return orders, nil
		}
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var accrual sql.NullFloat64
		var uploadedAt time.Time
		var number string
		var order models.Order
		err := rows.Scan(&status, &number, &accrual, &uploadedAt)
		if err != nil {
			return orders, err
		}

		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}
		order.Number = number
		order.UploadedAt = models.CustomTime{Time: uploadedAt}
		order.Status = status
		order.Login = login
		orders = append(orders, &order)
	}
	return orders, nil
}

func (p *PostgresOrderRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, error) {
	var order models.Order
	var status string
	var accrual sql.NullFloat64
	var login string
	var orderFromDB string
	var uploadedAt time.Time
	err := p.db.QueryRowContext(
		ctx,
		"SELECT login, number, status, accrual, uploaded_at FROM orders WHERE number = $1",
		orderNumber,
	).Scan(&login, &orderFromDB, &status, &accrual, &uploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("order - %s not found", orderNumber), zap.Error(err))
			return nil, err
		}
		return &order, err
	}

	if accrual.Valid {
		order.Accrual = &accrual.Float64
	}
	order.Number = orderFromDB
	order.UploadedAt = models.CustomTime{Time: uploadedAt}
	order.Status = status
	order.Login = login
	return &order, nil
}

func (p *PostgresOrderRepository) GetNotAcceptedOrderNumbers(ctx context.Context, limit uint) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return orders, err
	}
	rows, err := tx.QueryContext(
		ctx,
		"SELECT number FROM orders WHERE status = ANY($1) and is_update_status = false limit $2",
		pq.Array([]string{models.OrderStatusNew, models.OrderStatusProcessing}),
		limit,
	)
	if rows.Err() != nil {
		return orders, err
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(err.Error(), zap.Error(err))
			return orders, nil
		}
		return orders, err
	}
	defer rows.Close()

	orderNumbers := make([]string, len(orders))
	for rows.Next() {
		var number string
		err := rows.Scan(&number)
		if err != nil {
			return orders, err
		}

		orders = append(orders, &models.Order{Number: number})
		orderNumbers = append(orderNumbers, number)
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE orders SET is_update_status = true WHERE number = ANY($1)",
		pq.Array(orderNumbers),
	)
	if err != nil {
		err := tx.Rollback()
		return orders, err
	}
	if err := tx.Commit(); err != nil {
		return orders, err
	}

	return orders, nil
}

func (p *PostgresOrderRepository) UpdateStatus(ctx context.Context, order *models.Order) error {
	var loginFromDB string

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	var isUpdateStatus bool
	switch order.Status {
	case models.OrderStatusProcessed, models.OrderStatusInvalid:
		isUpdateStatus = true
	default:
		isUpdateStatus = false

	}

	err = tx.QueryRowContext(
		ctx,
		"UPDATE orders SET status = $1, accrual = $2, is_update_status = $3 WHERE number = $4 RETURNING login",
		order.Status,
		order.Accrual,
		isUpdateStatus,
		order.Number,
	).Scan(&loginFromDB)

	if err != nil {
		err = tx.Rollback()
		return err
	}

	if order.Accrual != nil {
		_, err = tx.ExecContext(
			ctx,
			"UPDATE balance SET current = current + $1 WHERE login = $2",
			order.Accrual,
			loginFromDB,
		)

		if err != nil {
			err = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (p *PostgresOrderRepository) IsExist(ctx context.Context, orderNumber string) (bool, error) {
	var status string
	err := p.db.QueryRowContext(ctx, "SELECT number FROM orders WHERE number = $1", orderNumber).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Info(fmt.Sprintf("order - %s not found", orderNumber), zap.Error(err))
			return false, nil
		}
		return false, err
	}
	return true, nil
}
