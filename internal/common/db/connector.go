package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func (dbcfg *DatabaseConfig) Connect() *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		dbcfg.Host, dbcfg.Port, dbcfg.Name, dbcfg.User, dbcfg.Password, dbcfg.SslMode)
	db, err := sqlx.ConnectContext(context.Background(), "postgres", connStr)
	if err != nil {
		logrus.Fatalln("error connecting to mpnex database:", err)
	}
	return db
}
