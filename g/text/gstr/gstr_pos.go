// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "strings"

// Find the position of the first occurrence of a substring in a string.
// It returns -1, if none found.
//
// 返回 needle 在 haystack 中首次出现的数字位置，找不到返回-1。
func Pos(haystack, needle string, startOffset...int) int {
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
    pos := strings.Index(haystack[offset : ], needle)
    if pos == -1 {
        return -1
    }
    return pos + offset
}

// Find the position of the first occurrence of a case-insensitive substring in a string.
// It returns -1, if none found.
//
// 返回在字符串 haystack 中 needle 首次出现的数字位置（不区分大小写），找不到返回-1。
func PosI(haystack, needle string, startOffset...int) int {
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
    pos := strings.Index(strings.ToLower(haystack[offset : ]), strings.ToLower(needle))
    if pos == -1 {
        return -1
    }
    return pos + offset
}

// Find the position of the last occurrence of a substring in a string.
// It returns -1, if none found.
//
// 查找指定字符串在目标字符串中最后一次出现的位置，找不到返回-1。
func PosR(haystack, needle string, startOffset...int) int {
    offset := 0
    if len(startOffset) > 0 {
        offset = startOffset[0]
    }
    pos, length := 0, len(haystack)
    if length == 0 || offset > length || -offset > length {
        return -1
    }

    if offset < 0 {
        haystack = haystack[ : offset + length + 1]
    } else {
        haystack = haystack[offset : ]
    }
    pos = strings.LastIndex(haystack, needle)
    if offset > 0 && pos != -1 {
        pos += offset
    }
    return pos
}

// Find the position of the last occurrence of a case-insensitive substring in a string.
// It returns -1, if none found.
//
// 以不区分大小写的方式查找指定字符串在目标字符串中最后一次出现的位置，找不到返回-1。
func PosRI(haystack, needle string, startOffset...int) int {
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
    pos = strings.LastIndex(strings.ToLower(haystack), strings.ToLower(needle))
    if offset > 0 && pos != -1 {
        pos += offset
    }
    return pos
}