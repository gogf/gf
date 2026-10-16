// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// sqlTokenKind is the kind of a lexical unit of a sql statement.
type sqlTokenKind int

const (
	sqlTokenSpace   sqlTokenKind = iota // Whitespace.
	sqlTokenWord                        // An unquoted identifier, keyword or number.
	sqlTokenQuoted                      // A string literal or a quoted identifier.
	sqlTokenComment                     // A line or block comment.
	sqlTokenSymbol                      // Any other character.
)

// sqlToken is a lexical unit of a sql statement.
type sqlToken struct {
	kind sqlTokenKind
	text string
}

// is reports whether the token is the word `word`, case-insensitively.
func (t sqlToken) is(word string) bool {
	return t.kind == sqlTokenWord && strings.EqualFold(t.text, word)
}

// dmlKeywords are the keywords beginning the statements in which INDEX is an identifier.
var dmlKeywords = map[string]bool{
	"SELECT": true, "WITH": true, "INSERT": true, "UPDATE": true, "DELETE": true, "MERGE": true,
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
// It rewrites the calls of GROUP_CONCAT, a MySQL function that DM does not have, to LISTAGG, and
// quotes INDEX, a reserved word of DM, where it is an unquoted identifier: anywhere in a SELECT,
// INSERT, UPDATE, DELETE or MERGE statement. String literals, quoted identifiers and comments are
// left as they are.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []any,
) (newSql string, newArgs []any, err error) {
	tokens := scanSqlTokens(sql)
	return d.Core.DoFilter(ctx, link, rewriteSqlTokens(tokens, isDmlStatement(tokens)), args)
}

// isDmlStatement reports whether `tokens` form a SELECT, INSERT, UPDATE, DELETE or MERGE
// statement, which may begin with a parenthesis.
func isDmlStatement(tokens []sqlToken) bool {
	for _, token := range tokens {
		switch {
		case token.kind == sqlTokenSpace || token.kind == sqlTokenComment || token.text == "(":
			continue
		case token.kind == sqlTokenWord:
			return dmlKeywords[strings.ToUpper(token.text)]
		default:
			return false
		}
	}
	return false
}

// rewriteSqlTokens joins `tokens` into sql, rewriting the calls of GROUP_CONCAT with
// rewriteGroupConcat, and quoting the words INDEX if `quoteIndex` is true.
func rewriteSqlTokens(tokens []sqlToken, quoteIndex bool) string {
	var b strings.Builder
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]
		if token.is("GROUP_CONCAT") {
			if open := nextSqlToken(tokens, i+1); open < len(tokens) && tokens[open].text == "(" {
				if end := closingParenthesis(tokens, open); end > 0 {
					if call, ok := rewriteGroupConcat(tokens[open+1:end], quoteIndex); ok {
						b.WriteString(call)
						i = end
						continue
					}
				}
			}
		}
		if quoteIndex && token.is("INDEX") {
			b.WriteString(`"INDEX"`)
			continue
		}
		b.WriteString(token.text)
	}
	return b.String()
}

// rewriteGroupConcat returns the call of LISTAGG that computes the call of GROUP_CONCAT with the
// arguments `args`, being `[DISTINCT] expr [ORDER BY ...] [SEPARATOR str]`, like
// `LISTAGG(DISTINCT expr, str) WITHIN GROUP (ORDER BY ...)`. The separator is ',' if it is not
// given, as it is on MySQL. It returns false for a call that concatenates several expressions,
// which is left to DM to reject.
func rewriteGroupConcat(args []sqlToken, quoteIndex bool) (string, bool) {
	var (
		distinct  string
		separator = "','"
		begin     = nextSqlToken(args, 0)
		order     = -1 // The index of ORDER of the ORDER BY clause.
		by        = -1 // The index of BY of the ORDER BY clause.
		sep       = -1 // The index of SEPARATOR.
		depth     int
	)
	if begin < len(args) && args[begin].is("DISTINCT") {
		distinct, begin = "DISTINCT ", nextSqlToken(args, begin+1)
	}
	for i := begin; i < len(args); i++ {
		switch {
		case args[i].text == "(":
			depth++
		case args[i].text == ")":
			depth--
		case depth > 0:
		case args[i].text == "," && order < 0 && sep < 0:
			return "", false
		case args[i].is("ORDER") && order < 0 && sep < 0:
			if next := nextSqlToken(args, i+1); next < len(args) && args[next].is("BY") {
				order, by = i, next
			}
		case args[i].is("SEPARATOR") && sep < 0:
			sep = i
		}
	}
	expressionEnd, orderEnd := len(args), len(args)
	if sep >= 0 {
		literal := nextSqlToken(args, sep+1)
		if literal == len(args) || args[literal].kind != sqlTokenQuoted || args[literal].text[0] != '\'' ||
			nextSqlToken(args, literal+1) != len(args) {
			return "", false
		}
		separator = args[literal].text
		expressionEnd, orderEnd = sep, sep
	}
	if order >= 0 {
		expressionEnd = order
	}
	expression := strings.TrimSpace(rewriteSqlTokens(args[begin:expressionEnd], quoteIndex))
	if expression == "" {
		return "", false
	}
	call := "LISTAGG(" + distinct + expression + ", " + separator + ")"
	if order >= 0 {
		orderBy := strings.TrimSpace(rewriteSqlTokens(args[by+1:orderEnd], quoteIndex))
		call += " WITHIN GROUP (ORDER BY " + orderBy + ")"
	}
	return call, true
}

// nextSqlToken returns the index of the first token from index `i` of `tokens` that is neither
// whitespace nor a comment, or the length of `tokens` if there is none.
func nextSqlToken(tokens []sqlToken, i int) int {
	for i < len(tokens) && (tokens[i].kind == sqlTokenSpace || tokens[i].kind == sqlTokenComment) {
		i++
	}
	return i
}

// closingParenthesis returns the index of the token closing the parenthesis at index `open` of
// `tokens`, or -1 if there is none.
func closingParenthesis(tokens []sqlToken, open int) int {
	var depth int
	for i := open; i < len(tokens); i++ {
		switch tokens[i].text {
		case "(":
			depth++
		case ")":
			if depth--; depth == 0 {
				return i
			}
		}
	}
	return -1
}

// scanSqlTokens splits `sql` into its tokens, so that joining their texts gives back `sql`.
// A string literal, a quoted identifier or a comment that is not closed runs to the end of `sql`.
func scanSqlTokens(sql string) (tokens []sqlToken) {
	for i := 0; i < len(sql); {
		var (
			c     = sql[i]
			begin = i
			kind  sqlTokenKind
		)
		switch {
		case isSqlSpace(c):
			kind = sqlTokenSpace
			for i < len(sql) && isSqlSpace(sql[i]) {
				i++
			}
		case c == '\'' || c == '"':
			kind = sqlTokenQuoted
			i = quotedEnd(sql, i)
		case strings.HasPrefix(sql[i:], "--"):
			kind = sqlTokenComment
			if n := strings.IndexByte(sql[i:], '\n'); n >= 0 {
				i += n + 1
			} else {
				i = len(sql)
			}
		case strings.HasPrefix(sql[i:], "/*"):
			kind = sqlTokenComment
			if n := strings.Index(sql[i+2:], "*/"); n >= 0 {
				i += n + 4
			} else {
				i = len(sql)
			}
		case isSqlWordChar(c):
			kind = sqlTokenWord
			for i < len(sql) && isSqlWordChar(sql[i]) {
				i++
			}
		default:
			kind = sqlTokenSymbol
			i++
		}
		tokens = append(tokens, sqlToken{kind: kind, text: sql[begin:i]})
	}
	return tokens
}

// quotedEnd returns the index following the string literal or quoted identifier that begins at
// index `open` of `sql`, taking a doubled quote character as an escaped one, or the length of
// `sql` if it is not closed.
func quotedEnd(sql string, open int) int {
	quote := sql[open]
	for i := open + 1; i < len(sql); i++ {
		if sql[i] != quote {
			continue
		}
		if i+1 < len(sql) && sql[i+1] == quote {
			i++
			continue
		}
		return i + 1
	}
	return len(sql)
}

// isSqlSpace reports whether `c` is whitespace.
func isSqlSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// isSqlWordChar reports whether `c` can be part of an unquoted identifier, keyword or number.
func isSqlWordChar(c byte) bool {
	return c == '_' || c == '$' || c == '#' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
