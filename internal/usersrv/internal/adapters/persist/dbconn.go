package persist

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // driver required by sqlx
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/app"
	"github.com/mobiletoly/moviex-backend/internal/usersrv/internal/core/outport"
)

type dbConnector struct {
	db *sqlx.DB
}

func NewDBConnector(cfg *app.Config) outport.DBConnector {
	db := cfg.DB.Connect()
	db.SetMaxOpenConns(4)
	return &dbConnector{db: db}
}

func (impl dbConnector) DB() *sqlx.DB {
	return impl.db
}
