// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql

import (
	"database/sql"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Result struct {
	sql.Result
	affected          int64
	lastInsertId      int64
	lastInsertIdError error
	returningRecords  []gdb.Record
	returningFields   []string
}

func (pgr Result) RowsAffected() (int64, error) {
	return pgr.affected, nil
}

func (pgr Result) LastInsertId() (int64, error) {
	return pgr.lastInsertId, pgr.lastInsertIdError
}

// ReturningRecords retrieves all records returned by RETURNING clause
func (pgr Result) ReturningRecords() ([]gdb.Record, error) {
	return pgr.returningRecords, nil
}

// ReturningValues retrieves all values of the specified field
func (pgr Result) ReturningValues(field string) ([]interface{}, error) {
	var values []interface{}
	for _, record := range pgr.returningRecords {
		if value, ok := record[field]; ok {
			values = append(values, value)
		}
	}
	return values, nil
}

// ReturningFirst retrieves the first returned record
func (pgr Result) ReturningFirst() (gdb.Record, error) {
	if len(pgr.returningRecords) > 0 {
		return pgr.returningRecords[0], nil
	}
	return nil, gerror.New("no returning records")
}

// ReturningCount retrieves the count of returned records
func (pgr Result) ReturningCount() int {
	return len(pgr.returningRecords)
}
