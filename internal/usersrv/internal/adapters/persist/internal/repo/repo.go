package repo

import "github.com/jmoiron/sqlx"

type DBRepo struct {
	db *sqlx.DB
}

func NewDBRepo(db *sqlx.DB) DBRepo {
	return DBRepo{db: db}
}
