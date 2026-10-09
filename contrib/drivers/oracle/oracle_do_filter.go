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
)

var (
	newSqlReplacementTmp = `
SELECT * FROM (
	SELECT GFORM.*, ROWNUM ROW_NUMBER__ FROM (%s) GFORM WHERE ROWNUM <= %d
) WHERE ROW_NUMBER__ > %d
`
)

const (
	rowNumLimitTmp  = `SELECT * FROM (%s) WHERE ROWNUM <= %d`
	derivedTableTmp = `SELECT * FROM (%s)`
)

func init() {
	var err error
	newSqlReplacementTmp, err = gdb.FormatMultiLineSqlToSingle(newSqlReplacementTmp)
	if err != nil {
		panic(err)
	}
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []any) (newSql string, newArgs []any, err error) {
	newSql = rewriteQuery(convertPlaceholders(sql))
	return d.Core.DoFilter(ctx, link, newSql, args)
}

// convertPlaceholders converts each placeholder '?' in `sql` to ":vN", N counting from 1,
// leaving string literals and comments as they are.
func convertPlaceholders(sql string) string {
	var (
		b     strings.Builder
		index int
	)
	b.Grow(len(sql))
	for i := 0; i < len(sql); i++ {
		if end := literalOrCommentEnd(sql, i); end >= 0 {
			b.WriteString(sql[i:end])
			i = end - 1
		} else if sql[i] == '?' {
			index++
			b.WriteString(":v" + strconv.Itoa(index))
		} else {
			b.WriteByte(sql[i])
		}
	}
	return b.String()
}

// literalOrCommentEnd returns the index following the string literal or the comment that begins
// at index `i` of `sql`, or -1 if none begins there. A quote that is never closed does not begin
// a literal, and a comment that is never closed runs to the end of `sql`.
func literalOrCommentEnd(sql string, i int) int {
	switch {
	case sql[i] == '\'':
		if end := closingQuote(sql, i); end >= 0 {
			return end + 1
		}
	case strings.HasPrefix(sql[i:], "--"):
		if n := strings.IndexByte(sql[i:], '\n'); n >= 0 {
			return i + n
		}
		return len(sql)
	case strings.HasPrefix(sql[i:], "/*"):
		if n := strings.Index(sql[i+2:], "*/"); n >= 0 {
			return i + 2 + n + 2
		}
		return len(sql)
	}
	return -1
}

// sqlToken is a lexical unit at the outermost parenthesis level of a sql statement: a word,
// a quoted literal, a punctuation character, or a whole parenthesized group.
type sqlToken struct {
	space string // The whitespace preceding the token.
	text  string
	group bool
}

// is reports whether the token is the keyword `word`, case-insensitively.
func (t sqlToken) is(word string) bool {
	return !t.group && strings.EqualFold(t.text, word)
}

// limitClause is a LIMIT clause in MySQL syntax, `LIMIT count`, `LIMIT offset,count` or
// `LIMIT count OFFSET offset`, spanning the tokens from index `begin` to index `end` exclusive.
type limitClause struct {
	begin  int
	end    int
	offset int
	count  int
}

// rewriteQuery rewrites the MySQL syntax that the core builds and Oracle rejects, in `sql`
// itself and in every parenthesized sub-query, innermost first: a LIMIT clause becomes a ROWNUM
// filter over the query it limits, keeping the clauses that follow it such as a lock clause,
// and a compound query is rewritten by rewriteCompoundQuery.
// It returns `sql` unchanged if its parentheses or quotes are unbalanced.
func rewriteQuery(sql string) string {
	tokens, trailing, ok := scanSqlTokens(sql)
	if !ok {
		return sql
	}
	for i, token := range tokens {
		if token.group {
			tokens[i].text = "(" + rewriteQuery(token.text[1:len(token.text)-1]) + ")"
		}
	}
	if len(tokens) > 0 && (tokens[0].group || tokens[0].is("SELECT") || tokens[0].is("WITH")) {
		tokens = rewriteLimitClause(tokens)
	}
	return joinSqlTokens(tokens) + trailing
}

// rewriteLimitClause rewrites the LIMIT clause at the outermost level of the query formed by
// `tokens` into a ROWNUM filter over the query it limits.
func rewriteLimitClause(tokens []sqlToken) []sqlToken {
	clause, ok := findLimitClause(tokens)
	if !ok {
		return rewriteCompoundQuery(tokens)
	}
	var (
		query   = strings.TrimSpace(joinSqlTokens(rewriteCompoundQuery(tokens[:clause.begin])))
		limited string
	)
	if clause.offset > 0 {
		limited = fmt.Sprintf(newSqlReplacementTmp, query, clause.offset+clause.count, clause.offset)
	} else {
		limited = fmt.Sprintf(rowNumLimitTmp, query, clause.count)
	}
	return append([]sqlToken{{space: tokens[0].space, text: limited}}, tokens[clause.end:]...)
}

// findLimitClause returns the first LIMIT clause among `tokens`.
func findLimitClause(tokens []sqlToken) (clause limitClause, ok bool) {
	for i, token := range tokens {
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

// rewriteCompoundQuery rewrites the compound query that the core builds as
// `(SELECT ...) UNION [ALL] (SELECT ...) ... [ORDER BY ...]`. Oracle rejects ORDER BY in a
// parenthesized operand, so such an operand becomes a sub-query in FROM. Oracle cannot resolve
// the column names of a `SELECT *` operand in the ORDER BY of the compound query either, so a
// compound query with ORDER BY becomes a derived table that the ORDER BY applies to.
func rewriteCompoundQuery(tokens []sqlToken) []sqlToken {
	if len(tokens) < 3 || !tokens[0].group || !tokens[1].is("UNION") {
		return tokens
	}
	var last int
	for i, token := range tokens {
		if !token.group || !isCompoundOperand(tokens, i) {
			continue
		}
		if operand := token.text[1 : len(token.text)-1]; hasOrderBy(operand) {
			tokens[i].text = "(" + fmt.Sprintf(derivedTableTmp, operand) + ")"
		}
		last = i
	}
	if !isOrderBy(tokens, last+1) {
		return tokens
	}
	compound := strings.TrimSpace(joinSqlTokens(tokens[:last+1]))
	return append(
		[]sqlToken{{space: tokens[0].space, text: fmt.Sprintf(derivedTableTmp, compound)}},
		tokens[last+1:]...,
	)
}

// isCompoundOperand reports whether the group at index `i` of `tokens` is an operand of the
// compound query formed by `tokens`.
func isCompoundOperand(tokens []sqlToken, i int) bool {
	return i == 0 || tokens[i-1].is("UNION") || (i > 1 && tokens[i-1].is("ALL") && tokens[i-2].is("UNION"))
}

// hasOrderBy reports whether the query `sql` has an ORDER BY clause at its outermost level.
func hasOrderBy(sql string) bool {
	tokens, _, _ := scanSqlTokens(sql)
	for i := range tokens {
		if isOrderBy(tokens, i) {
			return true
		}
	}
	return false
}

// isOrderBy reports whether an ORDER BY clause begins at index `i` of `tokens`.
func isOrderBy(tokens []sqlToken, i int) bool {
	return i+1 < len(tokens) && tokens[i].is("ORDER") && tokens[i+1].is("BY")
}

// sqlTokenNumber returns the integer that the token at index `i` of `tokens` is.
func sqlTokenNumber(tokens []sqlToken, i int) (int, bool) {
	if i >= len(tokens) || tokens[i].group {
		return 0, false
	}
	n, err := strconv.Atoi(tokens[i].text)
	return n, err == nil
}

// scanSqlTokens splits `sql` into tokens at its outermost parenthesis level, and returns them
// with the whitespace after the last one, so that joining them gives back `sql`. It returns
// false if a parenthesis or a quote in `sql` is unbalanced.
func scanSqlTokens(sql string) (tokens []sqlToken, trailing string, ok bool) {
	var end int
	for i := 0; i < len(sql); {
		c := sql[i]
		if isSqlSpace(c) {
			i++
			continue
		}
		var begin = i
		switch {
		case c == '(':
			if i = closingParenthesis(sql, i); i < 0 {
				return nil, "", false
			}
			i++
		case c == ')':
			return nil, "", false
		case c == '\'':
			if i = closingQuote(sql, i); i < 0 {
				return nil, "", false
			}
			i++
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

// closingParenthesis returns the index of the parenthesis closing the one at `open`, skipping
// quoted literals, or -1 if there is none.
func closingParenthesis(sql string, open int) int {
	var depth int
	for i := open; i < len(sql); i++ {
		switch sql[i] {
		case '\'':
			if i = closingQuote(sql, i); i < 0 {
				return -1
			}
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

// closingQuote returns the index of the quote closing the literal opened at `open`, taking a
// doubled quote as an escaped one, or -1 if there is none.
func closingQuote(sql string, open int) int {
	for i := open + 1; i < len(sql); i++ {
		if sql[i] != '\'' {
			continue
		}
		if i+1 < len(sql) && sql[i+1] == '\'' {
			i++
			continue
		}
		return i
	}
	return -1
}

func isSqlSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isSqlWordChar(c byte) bool {
	return c == '_' || c == '$' || c == '#' || c == '.' || c == ':' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
