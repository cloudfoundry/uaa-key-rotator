package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DbAwareQuerier struct {
	DB       *sqlx.DB
	DBScheme string
}

func (q DbAwareQuerier) Close() error {
	return q.DB.Close()
}

func (q DbAwareQuerier) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	reboundQuery, err := RebindForSQLDialect(query, q.DBScheme)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to query")
	}
	return q.DB.Queryx(reboundQuery, args...)
}
