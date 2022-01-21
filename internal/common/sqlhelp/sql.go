package sqlhelp

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mobiletoly/moviex-backend/internal/common/service"
	"reflect"
	"strings"
)

func RowsFromQueryByIDs(ctx context.Context, db *sqlx.DB, tag string, query string, ids interface{}) (*sqlx.Rows, error) {
	service.LogEntry(ctx).Debugf("%s: %s ||| SQL params: %s", tag, strings.ReplaceAll(query, "\n", " "),
		trimSliceForLog(ids))
	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	rows, err := db.Queryx(query, args...)
	return rows, err
}

func PrepareNamedStmt(ctx context.Context, db *sqlx.DB, tag string, query string, arg interface{}) (*sqlx.NamedStmt, error) {
	service.LogEntry(ctx).Debugf("%s: %s ||| SQL params: %v", tag, strings.ReplaceAll(query, "\n", " "), arg)
	stmt, err := db.PrepareNamed(query)
	return stmt, err
}

func Get(ctx context.Context, db *sqlx.DB, tag string, dest interface{}, query string, args ...interface{}) error {
	service.LogEntry(ctx).Debugf("%s: %s ||| SQL params: %v", tag, strings.ReplaceAll(query, "\n", " "), args)
	return db.Get(dest, query, args...)
}

func Select(ctx context.Context, db *sqlx.DB, tag string, dest interface{}, query string, args ...interface{}) error {
	service.LogEntry(ctx).Debugf("%s: %s ||| SQL params: %v", tag, strings.ReplaceAll(query, "\n", " "), args)
	return db.Select(dest, query, args...)
}

func trimSliceForLog(idsRef interface{}) string {
	v := reflect.ValueOf(idsRef)
	if v.Kind() != reflect.Slice {
		panic("trimSliceForLog() given a non-slice type")
	}
	if v.IsNil() {
		return "[]"
	}
	const maxlen = 10
	if v.Len() > maxlen {
		return fmt.Sprintf("(%v...%d more elements)", v.Slice(0, maxlen), v.Len()-maxlen)
	}
	return fmt.Sprintf("%v", v)
}
