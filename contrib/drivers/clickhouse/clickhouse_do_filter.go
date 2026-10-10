// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
)

// DoFilter handles the sql before posts it to database.
//
// A placeholder '?' becomes "$N" outside string literals, quoted identifiers and comments, a
// compound query that the core builds becomes one that ClickHouse accepts, and an UPDATE or a
// DELETE that the framework generates becomes the ALTER TABLE mutation of ClickHouse.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, originSql string, args []any,
) (newSql string, newArgs []any, err error) {
	newSql = rewriteCompoundQueries(originSql)
	if len(args) > 0 {
		newSql = convertPlaceholders(newSql)
	}

	// Only SQL generated through the framework is processed.
	if !d.getNeedParsedSqlFromCtx(ctx) {
		return newSql, args, nil
	}

	// replace STD SQL to Clickhouse SQL grammar
	modeRes, err := gregex.MatchString(filterTypePattern, strings.TrimSpace(newSql))
	if err != nil {
		return "", nil, err
	}
	if len(modeRes) == 0 {
		return newSql, args, nil
	}

	// Only delete/ UPDATE statements require filter
	switch strings.ToUpper(modeRes[0]) {
	case "UPDATE":
		// MySQL eg: UPDATE table_name SET field1=new-value1, field2=new-value2 [WHERE Clause]
		// Clickhouse eg: ALTER TABLE [db.]table UPDATE column1 = expr1 [, ...] WHERE filter_expr
		newSql, err = gregex.ReplaceStringFuncMatch(
			updateFilterPattern, newSql,
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
			deleteFilterPattern, newSql,
			func(s []string) string {
				return fmt.Sprintf("ALTER TABLE %s DELETE", s[1])
			},
		)
		if err != nil {
			return "", nil, err
		}
		return newSql, args, nil

	default:
		return newSql, args, nil
	}
}

func (d *Driver) getNeedParsedSqlFromCtx(ctx context.Context) bool {
	return ctx.Value(needParsedSqlInCtx) != nil
}

// convertPlaceholders converts each placeholder '?' in `sql` to "$N", N counting from 1,
// leaving string literals, quoted identifiers and comments as they are.
func convertPlaceholders(sql string) string {
	var (
		b     strings.Builder
		index int
	)
	b.Grow(len(sql))
	for i := 0; i < len(sql); i++ {
		if end := quotedOrCommentEnd(sql, i); end > i {
			b.WriteString(sql[i:end])
			i = end - 1
		} else if sql[i] == '?' {
			index++
			b.WriteString("$" + strconv.Itoa(index))
		} else {
			b.WriteByte(sql[i])
		}
	}
	return b.String()
}

// rewriteCompoundQueries rewrites the compound queries that the core builds as
// `(SELECT ...) UNION [ALL] (SELECT ...) ... [ORDER BY ...] [LIMIT ...]`, in `sql` itself and in
// every parenthesized sub-query: a UNION without ALL or DISTINCT becomes UNION DISTINCT, as
// ClickHouse requires one of them by default, and a compound query followed by clauses becomes a
// derived table that the clauses apply to, as ClickHouse applies them to its last operand.
func rewriteCompoundQueries(sql string) string {
	var (
		b         strings.Builder
		start     = len(sql) - len(strings.TrimLeft(sql, " \t\r\n"))
		compound  = start < len(sql) && sql[start] == '('
		hasUnion  bool
		lastWord  string
		lastEnd   = -1
		isOperand = compound
	)
	b.Grow(len(sql) + 16)
	b.WriteString(sql[:start])
	for i := start; i < len(sql); i++ {
		switch c := sql[i]; {
		case quotedOrCommentEnd(sql, i) > i:
			end := quotedOrCommentEnd(sql, i)
			b.WriteString(sql[i:end])
			i = end - 1

		case c == '(':
			end := closingParenthesis(sql, i)
			if end < 0 {
				return sql
			}
			b.WriteString("(" + rewriteCompoundQueries(sql[i+1:end]) + ")")
			if compound && isOperand {
				lastEnd = b.Len()
			}
			isOperand = false
			i = end

		case isWordChar(c):
			end := i
			for end < len(sql) && isWordChar(sql[end]) {
				end++
			}
			word := strings.ToUpper(sql[i:end])
			b.WriteString(sql[i:end])
			if word == "UNION" {
				hasUnion = true
				if strings.HasPrefix(strings.TrimLeft(sql[end:], " \t\r\n"), "(") {
					b.WriteString(" DISTINCT")
				}
			}
			isOperand = word == "UNION" || (lastWord == "UNION" && (word == "ALL" || word == "DISTINCT"))
			lastWord = word
			i = end - 1

		default:
			b.WriteByte(c)
		}
	}
	rewritten := b.String()
	if !compound || !hasUnion || lastEnd < 0 || strings.TrimSpace(rewritten[lastEnd:]) == "" {
		return rewritten
	}
	return rewritten[:start] + "SELECT * FROM (" + rewritten[start:lastEnd] + ")" + rewritten[lastEnd:]
}

// quotedOrCommentEnd returns the index following the string literal, the quoted identifier or
// the comment that begins at index `i` of `sql`, or `i` if none begins there.
func quotedOrCommentEnd(sql string, i int) int {
	switch c := sql[i]; {
	case c == '\'' || c == '"' || c == '`':
		for j := i + 1; j < len(sql); j++ {
			switch sql[j] {
			case '\\':
				j++
			case c:
				if j+1 < len(sql) && sql[j+1] == c {
					j++
					continue
				}
				return j + 1
			}
		}
		return len(sql)
	case strings.HasPrefix(sql[i:], "--"):
		if n := strings.IndexByte(sql[i:], '\n'); n >= 0 {
			return i + n + 1
		}
		return len(sql)
	case strings.HasPrefix(sql[i:], "/*"):
		if n := strings.Index(sql[i+2:], "*/"); n >= 0 {
			return i + 2 + n + 2
		}
		return len(sql)
	}
	return i
}

// closingParenthesis returns the index of the parenthesis closing the one at index `open` of
// `sql`, skipping string literals, quoted identifiers and comments, or -1 if there is none.
func closingParenthesis(sql string, open int) int {
	var depth int
	for i := open; i < len(sql); i++ {
		if end := quotedOrCommentEnd(sql, i); end > i {
			i = end - 1
			continue
		}
		switch sql[i] {
		case '(':
			depth++
		case ')':
			if depth--; depth == 0 {
				return i
			}
		}
	}
	return -1
}

func isWordChar(c byte) bool {
	return c == '_' || c >= 0x80 || ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
