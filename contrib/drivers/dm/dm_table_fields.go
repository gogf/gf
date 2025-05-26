// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	tableFieldsSqlTmp = `SELECT c.COLUMN_NAME, c.DATA_TYPE, c.DATA_DEFAULT, cc.COMMENTS FROM ALL_TAB_COLUMNS c LEFT JOIN ALL_COL_COMMENTS cc ON c.COLUMN_NAME=cc.COLUMN_NAME AND c.TABLE_NAME=cc.TABLE_NAME AND c.OWNER=cc.OWNER WHERE c.TABLE_NAME= '%s' AND c.OWNER= '%s'`
)

// TableFields retrieves and returns the fields' information of specified table of current schema.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	var (
		result gdb.Result
		link   gdb.Link
		// When no schema is specified, the configuration item is returned by default
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	// When usedSchema is empty, return the default link
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	// The link has been distinguished and no longer needs to judge the owner
	result, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(
			tableFieldsSqlTmp,
			strings.ToUpper(table),
			strings.ToUpper(d.GetSchema()),
		),
	)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		// m[NULLABLE] returns "N" "Y"
		// "N" means not null
		// "Y" means could be null
		var nullable bool
		if m["NULLABLE"].String() != "N" {
			nullable = true
		}
		fields[m["COLUMN_NAME"].String()] = &gdb.TableField{
			Index:   i,
			Name:    m["COLUMN_NAME"].String(),
			Type:    m["DATA_TYPE"].String(),
			Null:    nullable,
			Default: m["DATA_DEFAULT"].Val(),
			// Key:     m["Key"].String(),
			// Extra:   m["Extra"].String(),
			Comment: m["COMMENTS"].String(),
		}
	}
	return fields, nil
}
