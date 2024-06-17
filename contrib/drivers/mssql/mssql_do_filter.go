// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	selectSqlTmp          = `SELECT * FROM (SELECT TOP %d * FROM (SELECT TOP %d %s) as TMP1_ ) as TMP2_ `
	selectWithOrderSqlTmp = `
SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROWNUMBER_, %s ) as TMP_ 
WHERE TMP_.ROWNUMBER_ > %d AND TMP_.ROWNUMBER_ <= %d
`
)

func init() {
	var err error
	selectWithOrderSqlTmp, err = gdb.FormatMultiLineSqlToSingle(selectWithOrderSqlTmp)
	if err != nil {
		panic(err)
	}
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []interface{},
) (newSql string, newArgs []interface{}, err error) {
	var index int
	// Convert placeholder char '?' to string "@px".
	newSql, err = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})
	if err != nil {
		return "", nil, err
	}
	newSql, err = gregex.ReplaceString("\"", "", newSql)
	if err != nil {
		return "", nil, err
	}
	newSql, err = d.parseSql(newSql)
	if err != nil {
		return "", nil, err
	}
	newArgs = args
	return d.Core.DoFilter(ctx, link, newSql, newArgs)
}

// parseSql does some replacement of the sql before commits it to underlying driver,
// for support of microsoft sql server.
func (d *Driver) parseSql(toBeCommittedSql string) (string, error) {
	var (
		err       error
		operation = gstr.StrTillEx(toBeCommittedSql, " ")
		keyword   = strings.ToUpper(gstr.Trim(operation))
	)
	switch keyword {
	case "SELECT":
		toBeCommittedSql, err = d.handleSelectSqlReplacement(toBeCommittedSql)
		if err != nil {
			return "", err
		}
	}
	return toBeCommittedSql, nil
}

func (d *Driver) handleSelectSqlReplacement(toBeCommittedSql string) (newSql string, err error) {
	// SELECT * FROM USER WHERE ID=1 LIMIT 1
	match, err := gregex.MatchString(`^SELECT(.+)LIMIT 1$`, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(match) > 1 {
		return fmt.Sprintf(`SELECT TOP 1 %s`, match[1]), nil
	}

	// SELECT * FROM USER WHERE AGE>18 ORDER BY ID DESC LIMIT 100, 200
	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if gregex.IsMatchString(patten, toBeCommittedSql) == false {
		return toBeCommittedSql, nil
	}
	allMatch, err := gregex.MatchAllString(patten, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	var index = 1
	// LIMIT statement checks.
	if len(allMatch) < 2 ||
		(strings.HasPrefix(allMatch[index][0], "LIMIT") == false &&
			strings.HasPrefix(allMatch[index][0], "limit") == false) {
		return toBeCommittedSql, nil
	}
	if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", toBeCommittedSql) == false {
		return toBeCommittedSql, nil
	}
	// ORDER BY statement checks.
	var (
		selectStr = ""
		orderStr  = ""
		haveOrder = gregex.IsMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", toBeCommittedSql)
	)
	if haveOrder {
		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)ORDER BY)", toBeCommittedSql)
		if len(queryExpr) != 4 ||
			strings.EqualFold(queryExpr[1], "SELECT") == false ||
			strings.EqualFold(queryExpr[3], "ORDER BY") == false {
			return toBeCommittedSql, nil
		}
		selectStr = queryExpr[2]
		orderExpr, _ := gregex.MatchString("((?i)ORDER BY)(.+)((?i)LIMIT)", toBeCommittedSql)
		if len(orderExpr) != 4 ||
			strings.EqualFold(orderExpr[1], "ORDER BY") == false ||
			strings.EqualFold(orderExpr[3], "LIMIT") == false {
			return toBeCommittedSql, nil
		}
		orderStr = orderExpr[2]
	} else {
		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", toBeCommittedSql)
		if len(queryExpr) != 4 ||
			strings.EqualFold(queryExpr[1], "SELECT") == false ||
			strings.EqualFold(queryExpr[3], "LIMIT") == false {
			return toBeCommittedSql, nil
		}
		selectStr = queryExpr[2]
	}
	first, limit := 0, 0
	for i := 1; i < len(allMatch[index]); i++ {
		if len(strings.TrimSpace(allMatch[index][i])) == 0 {
			continue
		}
		if strings.HasPrefix(allMatch[index][i], "LIMIT") ||
			strings.HasPrefix(allMatch[index][i], "limit") {
			first, _ = strconv.Atoi(allMatch[index][i+1])
			limit, _ = strconv.Atoi(allMatch[index][i+2])
			break
		}
	}
	if haveOrder {
		toBeCommittedSql = fmt.Sprintf(
			selectWithOrderSqlTmp,
			orderStr, selectStr, first, first+limit,
		)
		return toBeCommittedSql, nil
	}

	if first == 0 {
		first = limit
	}
	toBeCommittedSql = fmt.Sprintf(
		selectSqlTmp,
		limit, first+limit, selectStr,
	)
	return toBeCommittedSql, nil
}
