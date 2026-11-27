// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// dm_tables.go lists the tables of one DM schema.

package dm

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	// tablesSqlTmp lists the tables of one owner. ALL_TABLES returns every table the
	// login can access, so the session schema alone does not limit the result.
	tablesSqlTmp = `SELECT TABLE_NAME FROM ALL_TABLES WHERE OWNER = ?`
)

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var (
		result gdb.Result
		// When no schema is specified, the configuration item is returned by default.
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		// owner is that schema as DM stores it, folded the same way as TableFields.
		owner = foldSchemaName(usedSchema)
	)
	// An empty owner matches no rows of ALL_TABLES, so a missing schema is an error.
	// GetSchema already falls back to the configured name, which Open requires.
	if owner == "" {
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			`dm.Tables failed: schema is empty`,
		)
	}
	// When schema is empty, return the default link.
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, tablesSqlTmp, owner)
	if err != nil {
		return
	}
	for _, m := range result {
		if v, ok := lookupRecordValue(m, "TABLE_NAME"); ok {
			tables = append(tables, v.String())
		}
	}
	return
}
