package gstr

import "strings"

// Find the position of the first occurrence of a substring in a string.
func Pos(haystack, needle string, offset int) int {
    length := len(haystack)
    if length == 0 || offset > length || -offset > length {
        return -1
    }

    if offset < 0 {
        offset += length
    }
    pos := strings.Index(haystack[offset:], needle)
    if pos == -1 {
        return -1
    }
    return pos + offset
}

// Find the position of the first occurrence of a case-insensitive substring in a string.
func PosI(haystack, needle string, offset int) int {
    length := len(haystack)
    if length == 0 || offset > length || -offset > length {
        return -1
    }

    haystack = haystack[offset:]
    if offset < 0 {
        offset += length
    }
    pos := strings.Index(strings.ToLower(haystack), strings.ToLower(needle))
    if pos == -1 {
        return -1
    }
    return pos + offset
}

// Find the position of the last occurrence of a substring in a string.
func PosR(haystack, needle string, offset int) int {
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

// Find the position of the last occurrence of a case-insensitive substring in a string.
func PosRI(haystack, needle string, offset int) int {
    pos, length := 0, len(haystack)
    if length == 0 || offset > length || -offset > length {
        return -1
    }

    if offset < 0 {
        haystack = haystack[:offset+length+1]
    } else {
        haystack = haystack[offset:]
    }
    pos = strings.LastIndex(strings.ToLower(haystack), strings.ToLower(needle))
    if offset > 0 && pos != -1 {
        pos += offset
    }
    return pos
}