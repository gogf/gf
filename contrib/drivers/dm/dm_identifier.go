// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"strings"
)

// FoldIdentifier returns `name` in upper case, as DM folds unquoted identifiers to upper case.
func (d *Driver) FoldIdentifier(name string) string {
	return strings.ToUpper(name)
}

// foldSchemaName returns the configured schema name `name` as DM stores it: in upper case, as DM
// folds an unquoted name, or, if it is written inside double quotes, the name inside them as it is.
func foldSchemaName(name string) string {
	if len(name) > 1 && name[0] == '"' && name[len(name)-1] == '"' {
		return name[1 : len(name)-1]
	}
	return strings.ToUpper(name)
}

// splitTableName splits `table`, a table name optionally qualified with its schema, like `user`,
// `"user"`, `SYSDBA.user` or `"SYSDBA"."user"`, into the schema, which is empty if `table` is not
// qualified, and the table name as DM stores them: an unquoted identifier in upper case, and a
// quoted identifier as it is, without its quote characters. It returns false if `table` is not
// such a name.
func splitTableName(table string) (schema, name string, ok bool) {
	var end int
	if name, end, ok = scanIdentifier(table, 0); !ok {
		return "", "", false
	}
	if end == len(table) {
		return "", name, true
	}
	if table[end] != '.' {
		return "", "", false
	}
	schema = name
	if name, end, ok = scanIdentifier(table, end+1); !ok || end != len(table) {
		return "", "", false
	}
	return schema, name, true
}

// scanIdentifier returns the identifier that begins at index `i` of `s` as DM stores it, being an
// unquoted identifier in upper case and a quoted identifier as it is, without its quote
// characters, and the index following it. It returns false if no identifier begins there.
func scanIdentifier(s string, i int) (identifier string, end int, ok bool) {
	if i < len(s) && s[i] == '"' {
		n := strings.IndexByte(s[i+1:], '"')
		if n <= 0 {
			return "", 0, false
		}
		return s[i+1 : i+1+n], i + n + 2, true
	}
	end = i
	for end < len(s) && isIdentifierChar(s[end]) {
		end++
	}
	if end == i {
		return "", 0, false
	}
	return strings.ToUpper(s[i:end]), end, true
}

// isIdentifierChar reports whether `c` can be part of an unquoted identifier: a letter, a digit,
// `_`, `$`, `#`, or a byte of a multibyte character.
func isIdentifierChar(c byte) bool {
	return c == '_' || c == '$' || c == '#' || c >= 0x80 ||
		('0' <= c && c <= '9') || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
