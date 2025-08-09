package datasource

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewDatabase(databaseUri string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseUri)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
