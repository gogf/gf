// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	chdriver "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
)

const (
	identifierPattern   = "(?:`(?:[^`\\\\]|\\\\.|``)*`|\"(?:[^\"\\\\]|\\\\.|\"\")*\"|[\\w\\x{80}-\\x{10FFFF}]+)"
	tableNamePattern    = identifierPattern + `(?:\.` + identifierPattern + `)?`
	updateFilterPattern = `(?is)^\s*UPDATE\s+(` + tableNamePattern + `)\s+SET\b`
	deleteFilterPattern = `(?is)^\s*DELETE\s+FROM\s+(` + tableNamePattern + `)(\s|$)`
)

// DoFilter handles the sql before posts it to database.
//
// A compound query that the core builds becomes one that ClickHouse accepts, the string
// literals, quoted identifiers and comments of a statement with arguments are escaped so that
// the underlying driver binds only the placeholders, and an UPDATE or a DELETE that the
// framework generates becomes the ALTER TABLE mutation of ClickHouse. A statement whose
// placeholders do not match its arguments fails with an error of code
// gcode.CodeInvalidParameter, as database/sql does for MySQL, since the underlying driver does
// not check them.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, originSql string, args []any,
) (newSql string, newArgs []any, err error) {
	newSql = rewriteCompoundQueries(originSql)
	if len(args) > 0 {
		var placeholders int
		if newSql, placeholders, err = escapePlaceholders(newSql, isNamed(args)); err != nil {
			return "", nil, err
		}
		if placeholders != len(args) {
			return "", nil, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`the statement has %d placeholders but %d arguments: %s`, placeholders, len(args), originSql,
			)
		}
	}

	// Only SQL generated through the framework is processed.
	if !d.getNeedParsedSqlFromCtx(ctx) {
		return newSql, args, nil
	}

	// MySQL eg: UPDATE table_name SET field1=new-value1, field2=new-value2 [WHERE Clause]
	// Clickhouse eg: ALTER TABLE [db.]table UPDATE column1 = expr1 [, ...] WHERE filter_expr
	match, err := gregex.MatchString(updateFilterPattern, newSql)
	if err != nil {
		return "", nil, err
	}
	if len(match) > 0 {
		return "ALTER TABLE " + match[1] + " UPDATE" + newSql[len(match[0]):], args, nil
	}

	// MySQL eg: DELETE FROM table_name [WHERE Clause]
	// Clickhouse eg: ALTER TABLE [db.]table [ON CLUSTER cluster] DELETE WHERE filter_expr
	match, err = gregex.MatchString(deleteFilterPattern, newSql)
	if err != nil {
		return "", nil, err
	}
	if len(match) > 0 {
		return "ALTER TABLE " + match[1] + " DELETE" + match[2] + newSql[len(match[0]):], args, nil
	}
	return newSql, args, nil
}

func (d *Driver) getNeedParsedSqlFromCtx(ctx context.Context) bool {
	return ctx.Value(needParsedSqlInCtx) != nil
}

// escapePlaceholders returns `sql` with the characters of its string literals, quoted identifiers
// and comments that the underlying driver would take for placeholders escaped, and the number of
// its placeholders: of '?' and of distinct "$N", or of distinct "@name" if `named`.
//
// The underlying driver binds named arguments to every "@name", and other arguments to every "$N"
// if the statement has one or else to every '?' not preceded by a backslash, wherever they are,
// and refuses a statement that has both a "$N" and a '?'. So in a literal or an identifier a '?',
// the '$' of a "$N" and, if `named`, the '@' of an "@name" is written as its escape sequence, like
// `\x3F`, and a comment holding one of them becomes a space. A statement with both '?' and "$N"
// placeholders fails with an error of code gcode.CodeInvalidParameter.
func escapePlaceholders(sql string, named bool) (string, int, error) {
	var (
		b          strings.Builder
		positional int
		numeric    = make(map[string]struct{})
		names      = make(map[string]struct{})
	)
	b.Grow(len(sql) + 8)
	for i := 0; i < len(sql); i++ {
		end := quotedOrCommentEnd(sql, i)
		switch {
		case end > i && !isCommentStart(sql, i):
			escapeQuoted(&b, sql[i:end], named)
			i = end - 1
		case end > i && hasPlaceholder(sql[i:end], named):
			b.WriteByte(' ')
			i = end - 1
		case end > i:
			b.WriteString(sql[i:end])
			i = end - 1
		default:
			end = placeholderEnd(sql, i, named)
			switch {
			case end == i:
				b.WriteByte(sql[i])
				continue
			case sql[i] == '?':
				positional++
			case sql[i] == '$':
				numeric[sql[i:end]] = struct{}{}
			default:
				names[sql[i:end]] = struct{}{}
			}
			b.WriteString(sql[i:end])
			i = end - 1
		}
	}
	switch {
	case named:
		return b.String(), len(names), nil
	case positional > 0 && len(numeric) > 0:
		return "", 0, gerror.NewCodef(
			gcode.CodeInvalidParameter, `the statement mixes placeholders "?" and "$N": %s`, sql,
		)
	default:
		return b.String(), positional + len(numeric), nil
	}
}

// escapeQuoted writes the string literal or quoted identifier `quoted` into `b` with each
// character that placeholderEnd takes for the beginning of a placeholder written as its escape
// sequence, like `\x3F` for '?', keeping a backslash before it as a backslash.
func escapeQuoted(b *strings.Builder, quoted string, named bool) {
	for i := 0; i < len(quoted); i++ {
		switch {
		case placeholderEnd(quoted, i, named) > i:
			fmt.Fprintf(b, `\x%02X`, quoted[i])
		case quoted[i] == '\\' && i+1 < len(quoted) && placeholderEnd(quoted, i+1, named) > i+1:
			b.WriteString(`\\`)
		case quoted[i] == '\\' && i+1 < len(quoted):
			b.WriteString(quoted[i : i+2])
			i++
		default:
			b.WriteByte(quoted[i])
		}
	}
}

// hasPlaceholder reports whether a placeholder that placeholderEnd recognizes begins in `s`.
func hasPlaceholder(s string, named bool) bool {
	for i := range len(s) {
		if placeholderEnd(s, i, named) > i {
			return true
		}
	}
	return false
}

// placeholderEnd returns the index following the placeholder that the underlying driver binds
// beginning at index `i` of `s`, a '?', a "$N" with N being digits or, if `named`, an "@name"
// with name being letters, digits and underscores, or `i` if none begins there.
func placeholderEnd(s string, i int, named bool) int {
	switch {
	case s[i] == '?':
		return i + 1
	case s[i] == '$' || (named && s[i] == '@'):
		end := i + 1
		for end < len(s) && (('0' <= s[end] && s[end] <= '9') ||
			(s[i] == '@' && (s[end] == '_' || ('a' <= s[end] && s[end] <= 'z') || ('A' <= s[end] && s[end] <= 'Z')))) {
			end++
		}
		if end > i+1 {
			return end
		}
	}
	return i
}

// isNamed reports whether all of `args` are named arguments, which the underlying driver binds to
// "@name" placeholders.
func isNamed(args []any) bool {
	for _, arg := range args {
		switch arg.(type) {
		case sql.NamedArg, chdriver.NamedValue, chdriver.NamedDateValue:
		default:
			return false
		}
	}
	return true
}

// rewriteCompoundQueries rewrites the compound queries that the core builds as
// `(SELECT ...) UNION [ALL] (SELECT ...) ... [ORDER BY ...] [LIMIT ...]`, in `sql` itself and in
// every parenthesized sub-query: a UNION without ALL or DISTINCT becomes UNION DISTINCT, as
// ClickHouse requires one of them by default, and a compound query followed by clauses becomes a
// derived table that the clauses apply to, as ClickHouse applies them to its last operand.
func rewriteCompoundQueries(sql string) string {
	if !containsFold(sql, "UNION") {
		return sql
	}
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
				if next := skipSpacesAndComments(sql, end); next < len(sql) && sql[next] == '(' {
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
	case c == '#' || strings.HasPrefix(sql[i:], "--"):
		if n := strings.IndexByte(sql[i:], '\n'); n >= 0 {
			return i + n + 1
		}
		return len(sql)
	case strings.HasPrefix(sql[i:], "/*"):
		var depth int
		for j := i; j+1 < len(sql); j++ {
			switch sql[j : j+2] {
			case "/*":
				depth++
				j++
			case "*/":
				if depth--; depth == 0 {
					return j + 2
				}
				j++
			}
		}
		return len(sql)
	}
	return i
}

// isCommentStart reports whether a comment begins at index `i` of `sql`: "--" or '#' up to the
// end of the line, or "/*" up to its "*/", which ClickHouse allows to nest.
func isCommentStart(sql string, i int) bool {
	return sql[i] == '#' || strings.HasPrefix(sql[i:], "--") || strings.HasPrefix(sql[i:], "/*")
}

// skipSpacesAndComments returns the index of the first character from index `i` of `sql` that is
// neither a white space nor in a comment, or the length of `sql` if there is none.
func skipSpacesAndComments(sql string, i int) int {
	for i < len(sql) {
		switch {
		case sql[i] == ' ' || sql[i] == '\t' || sql[i] == '\r' || sql[i] == '\n':
			i++
		case isCommentStart(sql, i):
			i = quotedOrCommentEnd(sql, i)
		default:
			return i
		}
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

// containsFold reports whether `s` contains `substr`, ignoring the case of letters.
func containsFold(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if strings.EqualFold(s[i:i+len(substr)], substr) {
			return true
		}
	}
	return false
}

func isWordChar(c byte) bool {
	return c == '_' || c >= 0x80 || ('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
