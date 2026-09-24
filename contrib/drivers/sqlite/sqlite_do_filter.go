// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gstr"
)

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []any,
) (newSql string, newArgs []any, err error) {
	// Special insert/ignore operation for sqlite.
	switch {
	case gstr.HasPrefix(sql, gdb.InsertOperationIgnore):
		sql = "INSERT OR IGNORE" + sql[len(gdb.InsertOperationIgnore):]

	case gstr.HasPrefix(sql, gdb.InsertOperationReplace):
		sql = "INSERT OR REPLACE" + sql[len(gdb.InsertOperationReplace):]
	}
	// The core parenthesizes compound and EXISTS sub-queries in ways SQLite rejects.
	sql = unwrapDoubledExistsParentheses(sql)
	sql = wrapUnionOperandsAsSubQueries(sql)
	return d.Core.DoFilter(ctx, link, sql, args)
}

// wrapUnionOperandsAsSubQueries rewrites the compound query the core builds as
// `(SELECT ...) UNION (SELECT ...)` into `SELECT * FROM (SELECT ...) UNION SELECT * FROM (SELECT ...)`.
// SQLite forbids parentheses around the operands of a compound SELECT, and also forbids
// ORDER BY and LIMIT on any operand but the last; wrapping each operand as a sub-query in
// FROM keeps both legal and preserves the meaning.
func wrapUnionOperandsAsSubQueries(sql string) string {
	if !strings.HasPrefix(sql, "(") {
		return sql
	}
	var (
		b        strings.Builder
		operands = 0
	)
	b.Grow(len(sql) + 64)
	for i := 0; i < len(sql); {
		if sql[i] != '(' {
			b.WriteByte(sql[i])
			i++
			continue
		}
		if !precededByUnionOrStart(sql, i) {
			b.WriteByte(sql[i])
			i++
			continue
		}
		end := matchingParenthesis(sql, i)
		if end < 0 {
			return sql
		}
		operands++
		b.WriteString("SELECT * FROM ")
		b.WriteString(sql[i : end+1])
		i = end + 1
	}
	if operands < 2 {
		return sql
	}
	return b.String()
}

// precededByUnionOrStart reports whether the parenthesis at `pos` opens a compound
// operand: it is the first character, or follows `UNION` / `UNION ALL`.
func precededByUnionOrStart(sql string, pos int) bool {
	if pos == 0 {
		return true
	}
	head := strings.ToUpper(strings.TrimRight(sql[:pos], " "))
	return strings.HasSuffix(head, " UNION") || strings.HasSuffix(head, " UNION ALL")
}

// unwrapDoubledExistsParentheses rewrites `EXISTS ((SELECT ...))`, which the core builds
// by parenthesizing an already parenthesized sub-query, into `EXISTS (SELECT ...)`.
// SQLite rejects the doubled parentheses as a syntax error.
func unwrapDoubledExistsParentheses(sql string) string {
	const marker = "EXISTS (("
	for {
		pos := gstr.PosI(sql, marker)
		if pos < 0 {
			return sql
		}
		outer := pos + len(marker) - 2
		end := matchingParenthesis(sql, outer)
		if end < 0 {
			return sql
		}
		sql = sql[:outer] + sql[outer+1:end] + sql[end+1:]
	}
}

// matchingParenthesis returns the index of the parenthesis closing the one at `open`,
// ignoring parentheses inside quoted strings and identifiers, or -1 if unbalanced.
func matchingParenthesis(sql string, open int) int {
	var (
		depth = 0
		quote byte
	)
	for i := open; i < len(sql); i++ {
		c := sql[i]
		if quote != 0 {
			if c == quote {
				if quote == '\'' && i+1 < len(sql) && sql[i+1] == '\'' {
					i++
					continue
				}
				quote = 0
			}
			continue
		}
		switch c {
		case '\'', '"', '`':
			quote = c
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
