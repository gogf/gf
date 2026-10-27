// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements SQLite user-table discovery.

package sqlitecgo

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

const (
	// tablesSqlTmp excludes SQLite's reserved internal table prefix from discovery.
	tablesSqlTmp = `SELECT NAME FROM SQLITE_MASTER WHERE TYPE='table' AND LOWER(SUBSTR(NAME,1,7)) <> 'sqlite_' ORDER BY NAME`
)

// Tables retrieves and returns the user tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}

	result, err = d.DoSelect(ctx, link, tablesSqlTmp)
	if err != nil {
		return
	}
	for _, m := range result {
		for _, v := range m {
			tables = append(tables, v.String())
		}
	}
	return
}
