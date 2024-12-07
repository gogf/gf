// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	tablesSqlTmp = `
SELECT
	c.relname
FROM
	pg_class c
INNER JOIN pg_namespace n ON
	c.relnamespace = n.oid
WHERE
	n.nspname = '%s'
	AND c.relkind IN ('r', 'p')
ORDER BY
	c.relname
`
)

func init() {
	var err error
	tablesSqlTmp, err = gdb.FormatMultiLineSqlToSingle(tablesSqlTmp)
	if err != nil {
		panic(err)
	}
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var (
		result     gdb.Result
		usedSchema = gutil.GetOrDefaultStr(d.GetConfig().Namespace, schema...)
	)
	if usedSchema == "" {
		usedSchema = d.GetConfig().User // GaussDB default schema use current user
	}
	// DO NOT use `usedSchema` as parameter for function `SlaveLink`.
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}

	var query = fmt.Sprintf(
		tablesSqlTmp,
		usedSchema,
	)

	query, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(query))
	result, err = d.DoSelect(ctx, link, query)
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
