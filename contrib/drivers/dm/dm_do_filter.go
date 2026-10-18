// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
)

var dmIndexKeywordReg = regexp.MustCompile(`(?i)\bINDEX\b`)

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []any,
) (newSql string, newArgs []any, err error) {
	// There should be no need to capitalize, because it has been done from field processing before
	// DM uses double quotes as identifier delimiters, so keep them intact.
	newSql = sql
	newSql = gstr.ReplaceI(gstr.ReplaceI(newSql, "GROUP_CONCAT", "LISTAGG"), "SEPARATOR", ",")

	// TODO The current approach is too rough. We should deal with the GROUP_CONCAT function and the
	// parsing of the index field from within the select from match.
	// ?GROUP_CONCAT DM  does not approve; index cannot be used as a query column name, and security characters need to be added, such as "index"?
	if strings.Contains(newSql, "INDEX") || strings.Contains(newSql, "index") {
		if !(strings.Contains(newSql, "_INDEX") || strings.Contains(newSql, "_index")) {
			newSql = quoteUnquotedIndexKeyword(newSql)
		}
	}

	// TODO i tried to do but it never work?
	// array, err := gregex.MatchAllString(`SELECT (.*INDEX.*) FROM .*`, newSql)
	// g.Dump("err:", err)
	// g.Dump("array:", array)
	// g.Dump("array:", array[0][1])

	// newSql, err = gregex.ReplaceString(`SELECT (.*INDEX.*) FROM .*`, l+"INDEX"+r, newSql)
	// g.Dump("err:", err)
	// g.Dump("newSql:", newSql)

	// re, err := regexp.Compile(`.*SELECT (.*INDEX.*) FROM .*`)
	// newSql = re.ReplaceAllStringFunc(newSql, func(data string) string {
	// 	fmt.Println("data:", data)
	// 	return data
	// })

	return d.Core.DoFilter(
		ctx,
		link,
		newSql,
		args,
	)
}

func quoteUnquotedIndexKeyword(sql string) string {
	if !dmIndexKeywordReg.MatchString(sql) {
		return sql
	}
	var (
		builder       strings.Builder
		start         int
		inSingleQuote bool
		inDoubleQuote bool
	)
	writeUnquotedPart := func(end int) {
		builder.WriteString(dmIndexKeywordReg.ReplaceAllString(sql[start:end], quoteChar+"INDEX"+quoteChar))
	}
	for i := 0; i < len(sql); i++ {
		switch sql[i] {
		case '\'':
			if inDoubleQuote {
				continue
			}
			if inSingleQuote && i+1 < len(sql) && sql[i+1] == '\'' {
				i++
				continue
			}
			if !inSingleQuote {
				writeUnquotedPart(i)
				start = i
				inSingleQuote = true
			} else {
				inSingleQuote = false
				builder.WriteString(sql[start : i+1])
				start = i + 1
			}

		case '"':
			if inSingleQuote {
				continue
			}
			if inDoubleQuote && i+1 < len(sql) && sql[i+1] == '"' {
				i++
				continue
			}
			if !inDoubleQuote {
				writeUnquotedPart(i)
				start = i
				inDoubleQuote = true
			} else {
				inDoubleQuote = false
				builder.WriteString(sql[start : i+1])
				start = i + 1
			}
		}
	}
	if start < len(sql) {
		if inSingleQuote || inDoubleQuote {
			builder.WriteString(sql[start:])
		} else {
			writeUnquotedPart(len(sql))
		}
	}
	return builder.String()
}
