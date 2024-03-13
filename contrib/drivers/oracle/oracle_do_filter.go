// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

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
	newSqlReplacementTmp = `
SELECT * FROM (
	SELECT GFORM.*, ROWNUM ROWNUM_ FROM (%s %s) GFORM WHERE ROWNUM <= %d
) 
	WHERE ROWNUM_ > %d
`
)

func init() {
	var err error
	newSqlReplacementTmp, err = gdb.FormatMultiLineSqlToSingle(newSqlReplacementTmp)
	if err != nil {
		panic(err)
	}
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	var index int
	newArgs = args
	// Convert placeholder char '?' to string ":vx".
	newSql, err = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf(":v%d", index)
	})
	if err != nil {
		return
	}
	newSql, err = gregex.ReplaceString("\"", "", newSql)
	if err != nil {
		return
	}
	newSql, err = d.parseSql(newSql)
	if err != nil {
		return
	}
	return d.Core.DoFilter(ctx, link, newSql, newArgs)
}

// parseSql does some replacement of the sql before commits it to underlying driver,
// for support of oracle server.
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
	var (
		match  [][]string
		patten = `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,{0,1}\s*(\d*))`
	)
	match, err = gregex.MatchAllString(patten, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(match) == 0 {
		return toBeCommittedSql, nil
	}
	var index = 1
	if len(match) < 2 || strings.HasPrefix(match[index][0], "LIMIT") == false {
		return toBeCommittedSql, nil
	}
	// only handle `SELECT ... LIMIT ...` statement.
	queryExpr, err := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(queryExpr) == 0 {
		return toBeCommittedSql, nil
	}
	if len(queryExpr) != 4 ||
		strings.EqualFold(queryExpr[1], "SELECT") == false ||
		strings.EqualFold(queryExpr[3], "LIMIT") == false {
		return toBeCommittedSql, nil
	}
	page, limit := 0, 0
	for i := 1; i < len(match[index]); i++ {
		if len(strings.TrimSpace(match[index][i])) == 0 {
			continue
		}
		if strings.HasPrefix(match[index][i], "LIMIT") {
			if match[index][i+2] != "" {
				page, err = strconv.Atoi(match[index][i+1])
				if err != nil {
					return "", err
				}
				limit, err = strconv.Atoi(match[index][i+2])
				if err != nil {
					return "", err
				}
				if page <= 0 {
					page = 1
				}
				limit = (page/limit + 1) * limit
				page, err = strconv.Atoi(match[index][i+1])
				if err != nil {
					return "", err
				}
			} else {
				limit, err = strconv.Atoi(match[index][i+1])
				if err != nil {
					return "", err
				}
			}
			break
		}
	}
	var newReplacedSql = fmt.Sprintf(
		newSqlReplacementTmp,
		queryExpr[1], queryExpr[2], limit, page,
	)
	return newReplacedSql, nil
}
