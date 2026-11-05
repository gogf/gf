// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"strings"
)

// targetAliasEndKeywords are the keywords that may follow the target table of an UPDATE or DELETE
// statement, so that none of them is taken as its alias.
var targetAliasEndKeywords = map[string]bool{
	"SET": true, "WHERE": true, "OPTION": true, "OUTPUT": true, "WITH": true, "FROM": true,
}

// tableWithoutAlias returns the table of the table expression `table` without the alias that
// follows it, with or without AS, keeping its schema, its database and its quotes, like `dbo.t` of
// `dbo.t AS a`, or `table` itself if it is not a table name followed by an alias.
func tableWithoutAlias(table string) string {
	tokens, _, ok := scanSqlTokens(table, nil)
	if !ok {
		return table
	}
	var words []string
	for _, token := range tokens {
		if token.group || token.isComment() {
			return table
		}
		if len(words) == 0 || token.space != "" {
			words = append(words, token.text)
		} else {
			words[len(words)-1] += token.text
		}
	}
	var alias string
	switch {
	case len(words) == 2 && !strings.EqualFold(words[1], "AS"):
		alias = words[1]
	case len(words) == 3 && strings.EqualFold(words[1], "AS"):
		alias = words[2]
	default:
		return table
	}
	if _, ok = splitTableName(words[0]); !ok {
		return table
	}
	if parts, ok := splitTableName(alias); !ok || len(parts) != 1 {
		return table
	}
	return words[0]
}

// rewriteAliasedTarget rewrites the UPDATE or DELETE statement formed by `tokens` whose target table
// has an alias, which SQL Server rejects there, into its form naming the alias as the target and
// the table with the alias in a FROM clause: `UPDATE t AS a SET ... WHERE ...` becomes
// `UPDATE a SET ... FROM t AS a WHERE ...`, and `DELETE FROM t a WHERE ...` becomes
// `DELETE a FROM t AS a WHERE ...`. It reports whether it rewrites the statement, and returns
// `tokens` unchanged for any other statement.
func rewriteAliasedTarget(tokens []sqlToken) ([]sqlToken, bool) {
	var (
		isDelete = tokens[0].is("DELETE")
		begin    = 1
	)
	if isDelete {
		if len(tokens) < 2 || !tokens[1].is("FROM") {
			return tokens, false
		}
		begin = 2
	}
	if begin >= len(tokens) {
		return tokens, false
	}
	end := begin + 1
	for end < len(tokens) && tokens[end].space == "" {
		end++
	}
	target := joinSqlTokens(tokens[begin:end])[len(tokens[begin].space):]
	if _, ok := splitTableName(target); !ok {
		return tokens, false
	}
	at := end
	if at < len(tokens) && tokens[at].is("AS") {
		at++
	}
	if at >= len(tokens) || !isTargetAlias(tokens[at]) {
		return tokens, false
	}
	next := at + 1
	if next < len(tokens) && tokens[next].space == "" {
		return tokens, false
	}
	alias := tokens[at].text
	if isDelete {
		if next < len(tokens) && !tokens[next].is("WHERE") && !tokens[next].is("OPTION") {
			return tokens, false
		}
		result := []sqlToken{tokens[0], {space: " ", text: alias + " FROM " + target + " AS " + alias}}
		return append(result, tokens[next:]...), true
	}
	if next >= len(tokens) || !tokens[next].is("SET") {
		return tokens, false
	}
	from := len(tokens)
	for i := next + 1; i < len(tokens); i++ {
		if tokens[i].is("WHERE") || tokens[i].is("OPTION") {
			from = i
			break
		}
	}
	result := make([]sqlToken, 0, len(tokens)+2)
	result = append(result, tokens[0], sqlToken{space: " ", text: alias})
	result = append(result, tokens[next:from]...)
	result = append(result, sqlToken{space: " ", text: "FROM " + target + " AS " + alias})
	return append(result, tokens[from:]...), true
}

// isTargetAlias reports whether `token` is an identifier that can be the alias of the target table
// of an UPDATE or DELETE statement.
func isTargetAlias(token sqlToken) bool {
	if token.group || targetAliasEndKeywords[strings.ToUpper(token.text)] {
		return false
	}
	parts, ok := splitTableName(token.text)
	return ok && len(parts) == 1
}
