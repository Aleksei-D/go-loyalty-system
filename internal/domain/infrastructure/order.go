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

func (p *PostgresOrderRepository) Add(ctx context.Context, orderNumber, login string) (*models.Order, error) {
	var order *models.Order
	var status string
	var accrual sql.NullFloat64
	var loginFromDB string
	var orderNumberFromDB string
	var uploadedAt time.Time

	row := p.db.QueryRowContext(
		ctx,
		"INSERT INTO order (orderNumber, login) VALUES ($1, $2) RETURNING login, orderNumber, status, accrual, uploaded_at",
		orderNumber,
		login,
	)
	if err := row.Err(); err != nil {
		logger.Log.Info(err.Error(), zap.Error(err))
		return order, err
	}

	err := row.Scan(&loginFromDB, &orderNumberFromDB, &status, &accrual, &uploadedAt)
	if err != nil {
		logger.Log.Info(err.Error(), zap.Error(err))
		return order, err
	}

	if accrual.Valid {
		order.Accrual = &accrual.Float64
	}
	order.Login = &loginFromDB
	order.Status = &status
	order.Number = &orderNumberFromDB
	order.UploadedAt = &models.CustomTime{Time: uploadedAt}
	return order, nil
}

func (p *PostgresOrderRepository) GetAllByLogin(ctx context.Context, login string) ([]*models.Order, error) {
	var orders []*models.Order
	rows, err := p.db.QueryContext(
		ctx,
		"SELECT status, number, accrual, uploaded_at FROM order WHERE login = $1 ORDER BY uploaded_at DESC",
		login,
	)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
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
		var order *models.Order
		err := rows.Scan(&status, &number, &accrual, &uploadedAt)
		if err != nil {
			return orders, err
		}

		if accrual.Valid {
			order.Accrual = &accrual.Float64
		}
		order.Number = &number
		order.UploadedAt = &models.CustomTime{Time: uploadedAt}
		order.Status = &status
		order.Login = &login
		orders = append(orders, order)
	}
	return orders, nil
}

func (p *PostgresOrderRepository) GetOrderByNumber(ctx context.Context, orderNumber string) (*models.Order, bool) {
	var order *models.Order
	var status string
	var accrual sql.NullFloat64
	var login string
	var orderFromDB string
	var uploadedAt time.Time
	row := p.db.QueryRowContext(
		ctx,
		"SELECT login, number, status, accrual, uploaded_at FROM order WHERE number = $1",
		orderNumber,
	)

	if err := row.Err(); err != nil {
		if errors.Is(sql.ErrNoRows, row.Err()) {
			logger.Log.Info(fmt.Sprintf("order - %s not found", orderNumber), zap.Error(err))
			return order, false
		}
		return order, false
	}

	err := row.Scan(&login, &orderFromDB, &status, &accrual, &uploadedAt)
	if err != nil {
		return order, false
	}

	if accrual.Valid {
		order.Accrual = &accrual.Float64
	}
	order.Number = &orderFromDB
	order.UploadedAt = &models.CustomTime{Time: uploadedAt}
	order.Status = &status
	order.Login = &login
	return order, true
}

func (p *PostgresOrderRepository) GetOrderByLoginAndNumber(ctx context.Context, login, orderNumber string) (*models.Order, bool) {
	var order *models.Order
	var status string
	var accrual sql.NullFloat64
	var loginFromDB string
	var uploadedAt time.Time
	var orderNumberFromDB string
	row := p.db.QueryRowContext(
		ctx,
		"SELECT login, number, status, accrual, uploaded_at FROM order WHERE number = $1 AND login = $2",
		orderNumber,
		login,
	)

	if err := row.Err(); err != nil {
		if errors.Is(sql.ErrNoRows, row.Err()) {
			logger.Log.Info(fmt.Sprintf("order - %s not found", orderNumber), zap.Error(err))
			return order, false
		}
		return order, false
	}

	err := row.Scan(&loginFromDB, &orderNumberFromDB, &status, &accrual, &uploadedAt)
	if err != nil {
		return order, false
	}

	if accrual.Valid {
		order.Accrual = &accrual.Float64
	}
	order.Number = &orderNumberFromDB
	order.UploadedAt = &models.CustomTime{Time: uploadedAt}
	order.Status = &status
	order.Login = &loginFromDB
	return order, true
}

func (p *PostgresOrderRepository) GetNotAcceptedOrderNumbers(ctx context.Context, limit uint) ([]*models.Order, error) {
	var orders []*models.Order
	rows, err := p.db.QueryContext(
		ctx,
		"SELECT number FROM order WHERE status = ANY($1) limit $2",
		pq.Array([]string{models.OrderStatusNew, models.OrderStatusProcessing}),
		limit,
	)
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var number string
		err := rows.Scan(&number)
		if err != nil {
			return orders, err
		}

		orders = append(orders, &models.Order{Number: &number})
	}
	return orders, nil
}

func (p *PostgresOrderRepository) UpdateStatus(ctx context.Context, order *models.Order) error {
	_, err := p.db.ExecContext(
		ctx,
		"UPDATE order SET status = $1, accrual = $2 WHERE number = $3",
		order.Status,
		order.Accrual,
		order.Number,
	)
	if err != nil {
		return err
	}
	return nil
}
