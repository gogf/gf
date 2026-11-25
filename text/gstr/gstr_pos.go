// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// toLowerOffsets converts `s` to lower case like strings.ToLower, but it keeps an invalid UTF-8 byte as it is.
// The lower case of a rune can have another byte length, for example "İ", so `offsets` gives for each byte
// of `lower` the offset in `s` of the rune it comes from, followed by len(s).
// `offsets` is nil if `s` is ASCII, as then each byte of `lower` is at the same offset in `s`.
func toLowerOffsets(s string) (lower string, offsets []int) {
	ascii := true
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			ascii = false
			break
		}
	}
	if ascii {
		return strings.ToLower(s), nil
	}
	buf := make([]byte, 0, len(s))
	offsets = make([]int, 0, len(s)+1)
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		n := len(buf)
		if r == utf8.RuneError && size == 1 {
			buf = append(buf, s[i])
		} else {
			buf = utf8.AppendRune(buf, unicode.ToLower(r))
		}
		for ; n < len(buf); n++ {
			offsets = append(offsets, i)
		}
		i += size
	}
	return string(buf), append(offsets, len(s))
}

// indexFold returns the byte range in `s` of the first occurrence of `search`, or of the last one
// if `last` is true, case-insensitively. It returns -1, -1 if `search` is not found.
// The range can be shorter or longer than `search`, as "İ" takes two bytes and its lower case "i" one.
func indexFold(s, search string, last bool) (start, end int) {
	lower, offsets := toLowerOffsets(s)
	searchLower, _ := toLowerOffsets(search)
	if searchLower == "" {
		if last {
			return len(s), len(s)
		}
		return 0, 0
	}
	for from, to := 0, len(lower); ; {
		var pos int
		if last {
			pos = strings.LastIndex(lower[:to], searchLower)
		} else if pos = strings.Index(lower[from:], searchLower); pos != -1 {
			pos += from
		}
		if pos == -1 {
			return -1, -1
		}
		end = pos + len(searchLower)
		if offsets == nil {
			return pos, end
		}
		// A match that starts or ends inside a rune of `s` is not a match.
		if (pos == 0 || offsets[pos] != offsets[pos-1]) && offsets[end] != offsets[end-1] {
			return offsets[pos], offsets[end]
		}
		from, to = pos+1, end-1
	}
}

// Pos returns the position of the first occurrence of `needle`
// in `haystack` from `startOffset`, case-sensitively.
// It returns -1, if not found.
func Pos(haystack, needle string, startOffset ...int) int {
	length := len(haystack)
	offset := 0
	if len(startOffset) > 0 {
		offset = startOffset[0]
	}
	if length == 0 || offset > length || -offset > length {
		return -1
	}
	if offset < 0 {
		offset += length
	}
	pos := strings.Index(haystack[offset:], needle)
	if pos == NotFoundIndex {
		return NotFoundIndex
	}
	return pos + offset
}

// PosRune acts like function Pos but considers `haystack` and `needle` as unicode string.
func PosRune(haystack, needle string, startOffset ...int) int {
	pos := Pos(haystack, needle, startOffset...)
	if pos < 2 {
		return pos
	}
	return len([]rune(haystack[:pos]))
}

// PosI returns the position of the first occurrence of `needle`
// in `haystack` from `startOffset`, case-insensitively.
// It returns -1, if not found.
func PosI(haystack, needle string, startOffset ...int) int {
	length := len(haystack)
	offset := 0
	if len(startOffset) > 0 {
		offset = startOffset[0]
	}
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		offset += length
	}
	pos, _ := indexFold(haystack[offset:], needle, false)
	if pos == -1 {
		return -1
	}
	return pos + offset
}

// PosIRune acts like function PosI but considers `haystack` and `needle` as unicode string.
func PosIRune(haystack, needle string, startOffset ...int) int {
	pos := PosI(haystack, needle, startOffset...)
	if pos < 2 {
		return pos
	}
	return len([]rune(haystack[:pos]))
}

// PosR returns the position of the last occurrence of `needle`
// in `haystack` from `startOffset`, case-sensitively.
// It returns -1, if not found.
func PosR(haystack, needle string, startOffset ...int) int {
	offset := 0
	if len(startOffset) > 0 {
		offset = startOffset[0]
	}
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos = strings.LastIndex(haystack, needle)
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// PosRRune acts like function PosR but considers `haystack` and `needle` as unicode string.
func PosRRune(haystack, needle string, startOffset ...int) int {
	pos := PosR(haystack, needle, startOffset...)
	if pos < 2 {
		return pos
	}
	return len([]rune(haystack[:pos]))
}

// PosRI returns the position of the last occurrence of `needle`
// in `haystack` from `startOffset`, case-insensitively.
// It returns -1, if not found.
func PosRI(haystack, needle string, startOffset ...int) int {
	offset := 0
	if len(startOffset) > 0 {
		offset = startOffset[0]
	}
	pos, length := 0, len(haystack)
	if length == 0 || offset > length || -offset > length {
		return -1
	}

	if offset < 0 {
		haystack = haystack[:offset+length+1]
	} else {
		haystack = haystack[offset:]
	}
	pos, _ = indexFold(haystack, needle, true)
	if offset > 0 && pos != -1 {
		pos += offset
	}
	return pos
}

// PosRIRune acts like function PosRI but considers `haystack` and `needle` as unicode string.
func PosRIRune(haystack, needle string, startOffset ...int) int {
	pos := PosRI(haystack, needle, startOffset...)
	if pos < 2 {
		return pos
	}
	return len([]rune(haystack[:pos]))
}
