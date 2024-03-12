// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	var (
		keys        []string // Field names.
		valueHolder = make([]string, 0)
	)
	// Handle the field names and placeholders.
	for k := range list[0] {
		keys = append(keys, k)
		valueHolder = append(valueHolder, "?")
	}
	// Prepare the batch result pointer.
	var (
		charL, charR = d.Core.GetChars()
		keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
		holderStr    = strings.Join(valueHolder, ",")
		tx           gdb.TX
		stmt         *gdb.Stmt
	)
	tx, err = d.Core.Begin(ctx)
	if err != nil {
		return
	}
	// It here uses defer to guarantee transaction be committed or roll-backed.
	defer func() {
		if err == nil {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()
	stmt, err = tx.Prepare(fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES (%s)",
		d.QuotePrefixTableName(table), keysStr,
		holderStr,
	))
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		// Values that will be committed to underlying database driver.
		params := make([]interface{}, 0)
		for _, k := range keys {
			params = append(params, list[i][k])
		}
		// Prepare is allowed to execute only once in a transaction opened by clickhouse
		result, err = stmt.ExecContext(ctx, params...)
		if err != nil {
			return
		}
	}
	return
}
