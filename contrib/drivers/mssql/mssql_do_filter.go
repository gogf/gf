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
	selectWithOrderSqlTmp = `
SELECT * FROM (
    SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROW_NUMBER__, %s 
    FROM (%s) as InnerQuery
) as TMP_ 
WHERE TMP_.ROW_NUMBER__ > %d AND TMP_.ROW_NUMBER__ <= %d`
	selectWithoutOrderSqlTmp = `
SELECT * FROM (
    SELECT ROW_NUMBER() OVER (ORDER BY (SELECT NULL)) as ROW_NUMBER__, %s 
    FROM (%s) as InnerQuery
) as TMP_ 
WHERE TMP_.ROW_NUMBER__ > %d AND TMP_.ROW_NUMBER__ <= %d`
)

func init() {
	var err error
	selectWithOrderSqlTmp, err = gdb.FormatMultiLineSqlToSingle(selectWithOrderSqlTmp)
	if err != nil {
		panic(err)
	}
	selectWithoutOrderSqlTmp, err = gdb.FormatMultiLineSqlToSingle(selectWithoutOrderSqlTmp)
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
	match, err := gregex.MatchString(`^SELECT(.+?)LIMIT\s+1$`, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(match) > 1 {
		return fmt.Sprintf(`SELECT TOP 1 %s`, strings.TrimSpace(match[1])), nil
	}

	// SELECT * FROM USER WHERE AGE>18 ORDER BY ID DESC LIMIT 100, 200
	pattern := `(?i)SELECT(.+?)(ORDER BY.+?)?\s*LIMIT\s*(\d+)(?:\s*,\s*(\d+))?`
	if !gregex.IsMatchString(pattern, toBeCommittedSql) {
		return toBeCommittedSql, nil
	}

	allMatch, err := gregex.MatchString(pattern, toBeCommittedSql)
	if err != nil {
		return "", err
	}

	// Extract SELECT part
	selectStr := strings.TrimSpace(allMatch[1])

	// Extract ORDER BY part
	orderStr := ""
	if len(allMatch[2]) > 0 {
		orderStr = strings.TrimSpace(allMatch[2])
		// Remove "ORDER BY" prefix as it will be used in OVER clause
		orderStr = strings.TrimPrefix(orderStr, "ORDER BY")
		orderStr = strings.TrimSpace(orderStr)
	}

	// Calculate LIMIT and OFFSET values
	first, _ := strconv.Atoi(allMatch[3]) // LIMIT first parameter
	limit := 0
	if len(allMatch) > 4 && allMatch[4] != "" {
		limit, _ = strconv.Atoi(allMatch[4]) // LIMIT second parameter
	} else {
		limit = first
		first = 0
	}

	// Build the final query
	if orderStr != "" {
		// Have ORDER BY clause
		newSql = fmt.Sprintf(
			selectWithOrderSqlTmp,
			orderStr,                            // ORDER BY clause for ROW_NUMBER
			"*",                                 // Select all columns
			fmt.Sprintf("SELECT %s", selectStr), // Original SELECT
			first,                               // OFFSET
			first+limit,                         // OFFSET + LIMIT
		)
	} else {
		// Without ORDER BY clause
		newSql = fmt.Sprintf(
			selectWithoutOrderSqlTmp,
			"*",                                 // Select all columns
			fmt.Sprintf("SELECT %s", selectStr), // Original SELECT
			first,                               // OFFSET
			first+limit,                         // OFFSET + LIMIT
		)
	}
	return newSql, nil
}
