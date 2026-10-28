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
	"sync"

	"github.com/gogf/gf/v2/database/gdb"
)

const (
	derivedTableTmp    = `SELECT * FROM (%s) AS TMP_`
	orderByNothing     = `ORDER BY (SELECT NULL)`
	orderByFirstColumn = `ORDER BY 1`
	// countTableAlias is the alias of the derived table whose rows the core counts for a grouped Count.
	countTableAlias = "count_alias"
	// countColumnAlias names the column of that derived table.
	countColumnAlias = "count_value"
)

// DoFilter deals with the sql string before commits it to underlying sql driver.
// A time argument is converted to the time zone of the driver, in which a column without time
// zone stores and returns its wall clock, as the type of the column it goes to is unknown here.
// `args` is copied before an argument is converted, and returned as it is if none is.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []any,
) (newSql string, newArgs []any, err error) {
	newSql = rewriteQuery(convertPlaceholders(rewriteSavePoint(sql)))
	if newSql, newArgs, err = d.Core.DoFilter(ctx, link, newSql, args); err != nil {
		return "", nil, err
	}
	var copied bool
	for i, arg := range newArgs {
		value, converted := d.timeInLocation(arg)
		if !converted {
			continue
		}
		if !copied {
			newArgs, copied = append([]any(nil), newArgs...), true
		}
		newArgs[i] = value
	}
	return newSql, newArgs, nil
}

// convertPlaceholders converts each placeholder '?' in `sql` to "@pN", N counting from 1,
// leaving string literals, quoted identifiers and comments as they are. It returns `sql` itself
// if it has no placeholder.
func convertPlaceholders(sql string) string {
	var (
		b     strings.Builder
		index int
		last  int
	)
	for i := 0; i < len(sql); i++ {
		if end := quotedOrCommentEnd(sql, i); end >= 0 {
			i = end - 1
			continue
		}
		if sql[i] != '?' {
			continue
		}
		if index == 0 {
			b.Grow(len(sql) + 4*strings.Count(sql[i:], "?"))
		}
		index++
		b.WriteString(sql[last:i])
		b.WriteString("@p")
		b.WriteString(strconv.Itoa(index))
		last = i + 1
	}
	if index == 0 {
		return sql
	}
	b.WriteString(sql[last:])
	return b.String()
}

// quotedOrCommentEnd returns the index following the string literal, the quoted identifier
// or the comment that begins at index `i` of `sql`, or -1 if none begins there. A quote that is
// never closed does not begin a literal, and a comment that is never closed runs to the end of
// `sql`. A line comment includes the line break that ends it.
func quotedOrCommentEnd(sql string, i int) int {
	switch {
	case sql[i] == '\'' || sql[i] == '"':
		if end := closingQuote(sql, i, sql[i]); end >= 0 {
			return end + 1
		}
	case sql[i] == '[':
		if end := closingQuote(sql, i, ']'); end >= 0 {
			return end + 1
		}
	case strings.HasPrefix(sql[i:], "--"):
		if n := strings.IndexByte(sql[i:], '\n'); n >= 0 {
			return i + n + 1
		}
		return len(sql)
	case strings.HasPrefix(sql[i:], "/*"):
		return blockCommentEnd(sql, i)
	}
	return -1
}

// blockCommentEnd returns the index following the block comment that begins at index `i`
// of `sql`, closing the block comments nested in it first as SQL Server does, or the length of
// `sql` if it is never closed.
func blockCommentEnd(sql string, i int) int {
	var depth int
	for i+1 < len(sql) {
		switch {
		case sql[i] == '/' && sql[i+1] == '*':
			depth++
			i += 2
		case sql[i] == '*' && sql[i+1] == '/':
			depth--
			i += 2
			if depth == 0 {
				return i
			}
		default:
			i++
		}
	}
	return len(sql)
}

// sqlToken is a lexical unit at the outermost parenthesis level of a sql statement: a word,
// a string literal, a quoted identifier, a comment, a punctuation character, or a whole
// parenthesized group.
type sqlToken struct {
	space string // The whitespace preceding the token.
	text  string
	group bool
}

// isComment reports whether the token is a comment.
func (t sqlToken) isComment() bool {
	return !t.group && (strings.HasPrefix(t.text, "--") || strings.HasPrefix(t.text, "/*"))
}

// is reports whether the token is the keyword `word`, case-insensitively.
func (t sqlToken) is(word string) bool {
	return !t.group && strings.EqualFold(t.text, word)
}

// sqlTokenBuffers pools the token slices that rewriteQuery scans statements into, as it keeps none
// of their tokens once it returns.
var sqlTokenBuffers = sync.Pool{
	New: func() any {
		return new([]sqlToken)
	},
}

// limitClause is a LIMIT clause in MySQL syntax, `LIMIT count`, `LIMIT offset,count` or
// `LIMIT count OFFSET offset`, or an OFFSET clause `OFFSET offset` without LIMIT, whose count
// is -1, spanning the tokens from index `begin` to index `end` exclusive.
type limitClause struct {
	begin  int
	end    int
	offset int
	count  int
}

// rewriteQuery rewrites the MySQL syntax that the core builds and SQL Server rejects, in
// `sql` itself and in every parenthesized sub-query, innermost first: a LIMIT or OFFSET clause
// becomes TOP or OFFSET FETCH, a lock clause becomes a table hint by moveLockClause, an operand
// of a compound query that has its own ORDER BY becomes a derived table, and the alias of the
// target table of an UPDATE or DELETE statement moves into a FROM clause by rewriteAliasedTarget,
// and the column of the derived table of a grouped Count is named by nameCountColumn.
// The comments that a query begins with are kept before it.
// It returns `sql` itself if it has nothing to rewrite, or if its parentheses or quotes are
// unbalanced.
func rewriteQuery(sql string) string {
	var (
		buffer               = sqlTokenBuffers.Get().(*[]sqlToken)
		tokens, trailing, ok = scanSqlTokens(sql, (*buffer)[:0])
		scanned              = tokens
	)
	defer func() {
		clear(scanned)
		*buffer = scanned[:0]
		sqlTokenBuffers.Put(buffer)
	}()
	if !ok {
		return sql
	}
	var rewritten bool
	for i, token := range tokens {
		if !token.group {
			continue
		}
		inner := token.text[1 : len(token.text)-1]
		if query := rewriteQuery(inner); query != inner {
			tokens[i].text = "(" + query + ")"
			rewritten = true
		}
	}
	var start int
	for start < len(tokens) && tokens[start].isComment() {
		start++
	}
	if start < len(tokens) {
		var (
			query   = tokens[start:]
			changed bool
		)
		switch first := tokens[start]; {
		case first.group || first.is("SELECT") || first.is("WITH"):
			var named, locked, operands, paged bool
			query, named = nameCountColumn(query)
			query, locked = moveLockClause(query)
			query, operands = rewriteCompoundOperands(query)
			query, paged = rewriteLimitClause(query)
			changed = named || locked || operands || paged
		case first.is("UPDATE") || first.is("DELETE"):
			query, changed = rewriteAliasedTarget(query)
		}
		if changed {
			tokens = append(tokens[:start:start], query...)
			rewritten = true
		}
	}
	if !rewritten {
		return sql
	}
	return joinSqlTokens(tokens) + trailing
}

// nameCountColumn names the column of the derived table of the statement formed by `tokens` in a
// column alias list if it is the one that the core builds for a grouped Count, `SELECT COUNT(1)
// FROM (SELECT COUNT(...) FROM ...) count_alias`, as SQL Server requires every column of a derived
// table to have a name, and reports whether it is.
func nameCountColumn(tokens []sqlToken) ([]sqlToken, bool) {
	if len(tokens) != 6 || !tokens[0].is("SELECT") || !tokens[1].is("COUNT") || tokens[2].text != "(1)" ||
		!tokens[3].is("FROM") || !tokens[4].group || !tokens[5].is(countTableAlias) {
		return tokens, false
	}
	tokens[5].text += "(" + countColumnAlias + ")"
	return tokens, true
}

// rewriteLimitClause rewrites the LIMIT or OFFSET clause at the outermost level of the
// query formed by `tokens`, and reports whether it has one. A simple query without offset selects
// the TOP rows, and any other query ends with OFFSET FETCH, ordered by nothing if it has no ORDER
// BY. As SQL Server only accepts the items of the select list in the ORDER BY of a compound or
// DISTINCT query, a compound query without ORDER BY becomes a derived table first, and a DISTINCT
// query without ORDER BY is ordered by its first column, whatever the names of its columns.
func rewriteLimitClause(tokens []sqlToken) ([]sqlToken, bool) {
	clause, ok := findLimitClause(tokens)
	if !ok {
		return tokens, false
	}
	var (
		query    = tokens[:clause.begin]
		compound = isCompoundQuery(query)
		distinct = query[0].is("SELECT") && len(query) > 1 && query[1].is("DISTINCT")
		result   = make([]sqlToken, 0, len(tokens)+1)
	)
	if clause.offset == 0 && clause.count >= 0 && query[0].is("SELECT") && !compound {
		var at = 1
		if at < len(query) && (query[at].is("DISTINCT") || query[at].is("ALL")) {
			at++
		}
		result = append(result, query[:at]...)
		result = append(result, sqlToken{space: " ", text: "TOP " + strconv.Itoa(clause.count)})
		result = append(result, query[at:]...)
		return append(result, tokens[clause.end:]...), true
	}
	var (
		space  = tokens[clause.begin].space
		paging = fmt.Sprintf("OFFSET %d ROWS", clause.offset)
	)
	if space == "" {
		space = " "
	}
	if clause.count >= 0 {
		paging += fmt.Sprintf(" FETCH NEXT %d ROWS ONLY", clause.count)
	}
	if !hasOrderBy(query) {
		switch {
		case compound:
			derived := fmt.Sprintf(derivedTableTmp, joinSqlTokens(query)[len(query[0].space):])
			query = []sqlToken{{space: query[0].space, text: derived}}
			paging = orderByNothing + " " + paging
		case distinct:
			paging = orderByFirstColumn + " " + paging
		default:
			paging = orderByNothing + " " + paging
		}
	}
	result = append(result, query...)
	result = append(result, sqlToken{space: space, text: paging})
	return append(result, tokens[clause.end:]...), true
}

// findLimitClause returns the first LIMIT clause among `tokens`, or the first OFFSET
// clause in MySQL syntax if it comes first.
func findLimitClause(tokens []sqlToken) (clause limitClause, ok bool) {
	for i, token := range tokens {
		if token.is("OFFSET") {
			offset, ok := sqlTokenNumber(tokens, i+1)
			if ok && !isRowsKeyword(tokens, i+2) {
				return limitClause{begin: i, end: i + 2, offset: offset, count: -1}, true
			}
			continue
		}
		if !token.is("LIMIT") {
			continue
		}
		count, ok := sqlTokenNumber(tokens, i+1)
		if !ok {
			continue
		}
		clause = limitClause{begin: i, end: i + 2, count: count}
		if i+2 < len(tokens) && tokens[i+2].text == "," {
			if n, ok := sqlTokenNumber(tokens, i+3); ok {
				clause.offset, clause.count, clause.end = count, n, i+4
			}
		} else if i+2 < len(tokens) && tokens[i+2].is("OFFSET") {
			if n, ok := sqlTokenNumber(tokens, i+3); ok {
				clause.offset, clause.end = n, i+4
			}
		}
		return clause, true
	}
	return clause, false
}

// rewriteCompoundOperands turns each parenthesized operand of the compound query formed by
// `tokens` that has its own ORDER BY into a derived table, as SQL Server rejects ORDER BY in an
// operand, and reports whether it turns any. The operand is given OFFSET 0 ROWS unless it has TOP
// or OFFSET, as SQL Server rejects ORDER BY in a derived table without either.
func rewriteCompoundOperands(tokens []sqlToken) ([]sqlToken, bool) {
	if !isCompoundQuery(tokens) {
		return tokens, false
	}
	var rewritten bool
	for i, token := range tokens {
		if !token.group || !isCompoundOperand(tokens, i) {
			continue
		}
		var (
			operand             = token.text[1 : len(token.text)-1]
			operandTokens, _, _ = scanSqlTokens(operand, nil)
		)
		if !hasOrderBy(operandTokens) {
			continue
		}
		if !hasPaging(operandTokens) {
			operand += " OFFSET 0 ROWS"
		}
		tokens[i].text = "(" + fmt.Sprintf(derivedTableTmp, operand) + ")"
		rewritten = true
	}
	return tokens, rewritten
}

// isCompoundQuery reports whether the query formed by `tokens` combines queries with a set
// operator at its outermost level.
func isCompoundQuery(tokens []sqlToken) bool {
	for _, token := range tokens {
		if isSetOperator(token) {
			return true
		}
	}
	return false
}

// isCompoundOperand reports whether the group at index `i` of `tokens` is an operand of the
// compound query formed by `tokens`.
func isCompoundOperand(tokens []sqlToken, i int) bool {
	return i == 0 || isSetOperator(tokens[i-1]) ||
		(i > 1 && tokens[i-1].is("ALL") && tokens[i-2].is("UNION"))
}

// isSetOperator reports whether `token` is an operator combining queries.
func isSetOperator(token sqlToken) bool {
	return token.is("UNION") || token.is("EXCEPT") || token.is("INTERSECT")
}

// hasOrderBy reports whether an ORDER BY clause is among `tokens`.
func hasOrderBy(tokens []sqlToken) bool {
	for i := 0; i+1 < len(tokens); i++ {
		if tokens[i].is("ORDER") && tokens[i+1].is("BY") {
			return true
		}
	}
	return false
}

// hasPaging reports whether a TOP clause or an OFFSET clause in T-SQL syntax is among
// `tokens`.
func hasPaging(tokens []sqlToken) bool {
	for i, token := range tokens {
		if token.is("TOP") || (token.is("OFFSET") && isRowsKeyword(tokens, i+2)) {
			return true
		}
	}
	return false
}

// isRowsKeyword reports whether the token at index `i` of `tokens` is ROW or ROWS.
func isRowsKeyword(tokens []sqlToken, i int) bool {
	return i < len(tokens) && (tokens[i].is("ROWS") || tokens[i].is("ROW"))
}

// sqlTokenNumber returns the integer that the token at index `i` of `tokens` is.
func sqlTokenNumber(tokens []sqlToken, i int) (int, bool) {
	if i >= len(tokens) || tokens[i].group {
		return 0, false
	}
	n, err := strconv.Atoi(tokens[i].text)
	return n, err == nil
}

// scanSqlTokens splits `sql` into tokens at its outermost parenthesis level, appends them to
// `tokens` and returns them with the whitespace after the last one, so that joining them gives back
// `sql`. It returns false if a parenthesis, a quote or a bracket in `sql` is unbalanced.
func scanSqlTokens(sql string, tokens []sqlToken) ([]sqlToken, string, bool) {
	var end int
	for i := 0; i < len(sql); {
		c := sql[i]
		if isSqlSpace(c) {
			i++
			continue
		}
		var begin = i
		switch {
		case isQuoteOrCommentStart(sql, i):
			if i = quotedOrCommentEnd(sql, i); i < 0 {
				return nil, "", false
			}
		case c == '(':
			if i = closingParenthesis(sql, i); i < 0 {
				return nil, "", false
			}
			i++
		case c == ')':
			return nil, "", false
		case isSqlWordChar(c):
			for i < len(sql) && isSqlWordChar(sql[i]) {
				i++
			}
		default:
			i++
		}
		tokens = append(tokens, sqlToken{space: sql[end:begin], text: sql[begin:i], group: c == '('})
		end = i
	}
	return tokens, sql[end:], true
}

// joinSqlTokens joins `tokens` back into sql.
func joinSqlTokens(tokens []sqlToken) string {
	var b strings.Builder
	for _, token := range tokens {
		b.WriteString(token.space)
		b.WriteString(token.text)
	}
	return b.String()
}

// closingParenthesis returns the index of the parenthesis closing the one at `open`,
// skipping string literals, quoted identifiers and comments, or -1 if there is none.
func closingParenthesis(sql string, open int) int {
	var depth int
	for i := open; i < len(sql); i++ {
		if isQuoteOrCommentStart(sql, i) {
			end := quotedOrCommentEnd(sql, i)
			if end < 0 {
				return -1
			}
			i = end - 1
			continue
		}
		switch sql[i] {
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

// closingQuote returns the index of the `quote` closing the string literal or the quoted
// identifier opened at `open`, taking a doubled `quote` as an escaped one, or -1 if there is none.
func closingQuote(sql string, open int, quote byte) int {
	for i := open + 1; i < len(sql); i++ {
		if sql[i] != quote {
			continue
		}
		if i+1 < len(sql) && sql[i+1] == quote {
			i++
			continue
		}
		return i
	}
	return -1
}

// isQuoteOrCommentStart reports whether a string literal, a quoted identifier or a comment
// begins at index `i` of `sql`.
func isQuoteOrCommentStart(sql string, i int) bool {
	return sql[i] == '\'' || sql[i] == '"' || sql[i] == '[' ||
		strings.HasPrefix(sql[i:], "--") || strings.HasPrefix(sql[i:], "/*")
}

func isSqlSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isSqlWordChar(c byte) bool {
	return c == '_' || c == '@' || c == '#' || c == '$' || c == '.' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
