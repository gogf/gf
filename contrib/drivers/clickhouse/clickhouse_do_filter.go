// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
)

// DoFilter handles the sql before posts it to database.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, originSql string, args []interface{},
) (newSql string, newArgs []interface{}, err error) {
	if len(args) == 0 {
		return originSql, args, nil
	}
	// Convert placeholder char '?' to string "$x".
	var index int
	originSql, _ = gregex.ReplaceStringFunc(`\?`, originSql, func(s string) string {
		index++
		return fmt.Sprintf(`$%d`, index)
	})

	// Only SQL generated through the framework is processed.
	if !d.getNeedParsedSqlFromCtx(ctx) {
		return originSql, args, nil
	}

	// replace STD SQL to Clickhouse SQL grammar
	modeRes, err := gregex.MatchString(filterTypePattern, strings.TrimSpace(originSql))
	if err != nil {
		return "", nil, err
	}
	if len(modeRes) == 0 {
		return originSql, args, nil
	}

	// Only delete/ UPDATE statements require filter
	switch strings.ToUpper(modeRes[0]) {
	case "UPDATE":
		// MySQL eg: UPDATE table_name SET field1=new-value1, field2=new-value2 [WHERE Clause]
		// Clickhouse eg: ALTER TABLE [db.]table UPDATE column1 = expr1 [, ...] WHERE filter_expr
		newSql, err = gregex.ReplaceStringFuncMatch(
			updateFilterPattern, originSql,
			func(s []string) string {
				return fmt.Sprintf("ALTER TABLE %s UPDATE", s[1])
			},
		)
		if err != nil {
			return "", nil, err
		}
		return newSql, args, nil

	case "DELETE":
		// MySQL eg: DELETE FROM table_name [WHERE Clause]
		// Clickhouse eg: ALTER TABLE [db.]table [ON CLUSTER cluster] DELETE WHERE filter_expr
		newSql, err = gregex.ReplaceStringFuncMatch(
			deleteFilterPattern, originSql,
			func(s []string) string {
				return fmt.Sprintf("ALTER TABLE %s DELETE", s[1])
			},
		)
		if err != nil {
			return "", nil, err
		}
		return newSql, args, nil

	}
	return originSql, args, nil
}

func (d *Driver) getNeedParsedSqlFromCtx(ctx context.Context) bool {
	if ctx.Value(needParsedSqlInCtx) != nil {
		return true
	}
	return false
}
