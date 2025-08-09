package infrastructure

import (
	"context"
	"database/sql"
	"github.com/Aleksei-D/go-loyalty-system/internal/models"
)

type PostgresBalanceRepository struct {
	db *sql.DB
}

func NewPostgresBalanceRepository(db *sql.DB) *PostgresBalanceRepository {
	return &PostgresBalanceRepository{db: db}
}

func (p *PostgresBalanceRepository) Get(ctx context.Context, login string) (*models.Balance, error) {
	var balance *models.Balance
	var current float64
	var withdrawn float64
	var loginFromDB string
	row := p.db.QueryRowContext(ctx, "SELECT b.login, b.current, SUM(w.sum) as withdrawns FROM balance as b JOIN withdrawal as w ON b.login = w.login WHERE b.login = $1 GROUP BY b.login", login)
	if err := row.Err(); err != nil {
		return balance, err
	}

	err := row.Scan(&loginFromDB, &current, &withdrawn)
	if err != nil {
		return balance, err
	}
	balance.Login = &loginFromDB
	balance.Current = &current
	balance.Withdrawn = &withdrawn
	return balance, nil
}
