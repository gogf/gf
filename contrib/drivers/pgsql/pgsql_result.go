package pgsql

import "database/sql"

type Result struct {
	sql.Result
	affected          int64
	lastInsertId      int64
	lastInsertIdError error
}

func (pgr Result) RowsAffected() (int64, error) {
	return pgr.affected, nil
}

func (pgr Result) LastInsertId() (int64, error) {
	return pgr.lastInsertId, pgr.lastInsertIdError
}
