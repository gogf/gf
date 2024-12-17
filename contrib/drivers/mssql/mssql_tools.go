package mssql

import "strings"

// RemoveSymbols removes all symbols from string and lefts only numbers and letters.
func removeSymbols(s string) string {
	var b = make([]rune, 0, len(s))
	for _, c := range s {
		if c > 127 {
			b = append(b, c)
		} else if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			b = append(b, c)
		}
	}
	return string(b)
}

// equalFoldWithoutChars checks string `s1` and `s2` equal case-insensitively,
// with/without chars '-'/'_'/'.'/' '.
func equalFoldWithoutChars(s1, s2 string) bool {
	return strings.EqualFold(removeSymbols(s1), removeSymbols(s2))
}
