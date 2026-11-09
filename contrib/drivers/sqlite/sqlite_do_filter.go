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
	sql = rewriteSubQueries(sql)
	return d.Core.DoFilter(ctx, link, sql, args)
}

// sqlToken is a lexical unit at the outermost parenthesis level of a sql statement: a word, a
// string literal, a quoted identifier, a comment, a punctuation character, or a whole
// parenthesized group.
type sqlToken struct {
	space string // The whitespace preceding the token.
	text  string
	group bool
}

// is reports whether the token is the keyword `word`, case-insensitively.
func (t sqlToken) is(word string) bool {
	return !t.group && strings.EqualFold(t.text, word)
}

// rewriteSubQueries rewrites the parenthesized queries the core builds into forms SQLite accepts,
// in `sql` itself and in every parenthesized group, innermost first:
//   - the operands of a compound query, `(SELECT ...) UNION (SELECT ...)`, become sub-queries in
//     FROM, `SELECT * FROM (SELECT ...) UNION SELECT * FROM (SELECT ...)`, as SQLite forbids
//     parentheses around them and ORDER BY or LIMIT on any operand but the last;
//   - a query in doubled parentheses, `((SELECT ...))`, keeps one pair, as SQLite rejects the
//     doubled pair after EXISTS.
//
// It returns `sql` unchanged if its parentheses or quotes are unbalanced.
func rewriteSubQueries(sql string) string {
	tokens, trailing, ok := scanSqlTokens(sql)
	if !ok {
		return sql
	}
	for i, token := range tokens {
		if !token.group {
			continue
		}
		inner := rewriteSubQueries(token.text[1 : len(token.text)-1])
		if query, ok := singleParenthesizedQuery(inner); ok {
			inner = query
		}
		tokens[i].text = "(" + inner + ")"
	}
	if isCompoundQuery(tokens) {
		for i, token := range tokens {
			if token.group && isCompoundOperand(tokens, i) {
				tokens[i].text = "SELECT * FROM " + token.text
			}
		}
	}
	return joinSqlTokens(tokens) + trailing
}

// singleParenthesizedQuery returns the content of `sql` if `sql` is nothing but one parenthesized
// query.
func singleParenthesizedQuery(sql string) (string, bool) {
	tokens, _, ok := scanSqlTokens(sql)
	if !ok || len(tokens) != 1 || !tokens[0].group {
		return "", false
	}
	content := tokens[0].text[1 : len(tokens[0].text)-1]
	inner, _, ok := scanSqlTokens(content)
	if !ok || len(inner) == 0 || !(inner[0].is("SELECT") || inner[0].is("WITH")) {
		return "", false
	}
	return content, true
}

// isCompoundQuery reports whether `tokens` combine queries with a set operator at their
// outermost level.
func isCompoundQuery(tokens []sqlToken) bool {
	for _, token := range tokens {
		if isSetOperator(token) {
			return true
		}
	}
	return false
}

// isCompoundOperand reports whether the group at index `i` of `tokens` is an operand of the
// compound query they form.
func isCompoundOperand(tokens []sqlToken, i int) bool {
	return i == 0 || isSetOperator(tokens[i-1]) ||
		(i > 1 && tokens[i-1].is("ALL") && tokens[i-2].is("UNION"))
}

// isSetOperator reports whether `token` is an operator combining queries.
func isSetOperator(token sqlToken) bool {
	return token.is("UNION") || token.is("EXCEPT") || token.is("INTERSECT")
}

// scanSqlTokens splits `sql` into tokens at its outermost parenthesis level, and returns them
// with the whitespace after the last one, so that joining them gives back `sql`. It returns
// false if a parenthesis, a quote or a comment in `sql` is unbalanced.
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

// closingParenthesis returns the index of the parenthesis closing the one at `open`, skipping
// string literals, quoted identifiers and comments, or -1 if there is none.
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

// quotedOrCommentEnd returns the index following the string literal, the quoted identifier or the
// comment that begins at index `i` of `sql`, or -1 if it is never closed. A doubled quote inside a
// literal or identifier is an escaped one.
func quotedOrCommentEnd(sql string, i int) int {
	switch c := sql[i]; {
	case c == '\'' || c == '"' || c == '`' || c == '[':
		closing := c
		if c == '[' {
			closing = ']'
		}
		for j := i + 1; j < len(sql); j++ {
			if sql[j] != closing {
				continue
			}
			if closing != ']' && j+1 < len(sql) && sql[j+1] == closing {
				j++
				continue
			}
			return j + 1
		}
		return -1
	case strings.HasPrefix(sql[i:], "--"):
		if n := strings.IndexByte(sql[i:], '\n'); n >= 0 {
			return i + n + 1
		}
		return len(sql)
	default:
		if n := strings.Index(sql[i+2:], "*/"); n >= 0 {
			return i + 2 + n + 2
		}
		return -1
	}
}

// isQuoteOrCommentStart reports whether a string literal, a quoted identifier or a comment begins
// at index `i` of `sql`.
func isQuoteOrCommentStart(sql string, i int) bool {
	switch sql[i] {
	case '\'', '"', '`', '[':
		return true
	}
	return strings.HasPrefix(sql[i:], "--") || strings.HasPrefix(sql[i:], "/*")
}

func isSqlSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isSqlWordChar(c byte) bool {
	return c == '_' || c == '$' || c == '.' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
