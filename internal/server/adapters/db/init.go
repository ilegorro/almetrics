package db

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Adapter struct {
	Conn *sql.DB
}

func New(dsn string) (*Adapter, error) {
	a := &Adapter{}
	db, err := sql.Open("pgx", dsn)
	if err == nil {
		a.Conn = db
	}

	return a, err
}

func (a *Adapter) Close() {
	a.Conn.Close()
}
