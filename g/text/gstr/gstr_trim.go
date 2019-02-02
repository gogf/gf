// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "strings"

// Strip whitespace (or other characters) from the beginning and end of a string.
//
// 去除字符串首尾处的空白字符（或者其他字符）。
func Trim(str string, characterMask ...string) string {
    if len(characterMask) > 0 {
        return strings.Trim(str, characterMask[0])
    } else {
        return strings.TrimSpace(str)
    }
}

// Strip whitespace (or other characters) from the beginning of a string.
//
// 去除字符串首的空白字符（或者其他字符）。
func TrimLeft(str string, characterMask ...string) string {
    mask := ""
    if len(characterMask) == 0 {
        mask = string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0})
    } else {
        mask = characterMask[0]
    }
    return strings.TrimLeft(str, mask)
}

// Strip all of the given <cut> string from the beginning of a string.
//
// 去除字符串首的给定字符串。
func TrimLeftStr(str string, cut string) string {
    for str[0 : len(cut)] == cut {
        str = str[len(cut) : ]
    }
    return str
}

// Strip whitespace (or other characters) from the end of a string.
//
// 去除字符串尾的空白字符（或者其他字符）。
func TrimRight(str string, characterMask ...string) string {
    mask := ""
    if len(characterMask) == 0 {
        mask = string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0})
    } else {
        mask = characterMask[0]
    }
    return strings.TrimRight(str, mask)
}

// Strip all of the given <cut> string from the end of a string.
//
// 去除字符串尾的给定字符串。
func TrimRightStr(str string, cut string) string {
    for {
        length := len(str)
        if str[length - len(cut) : length] == cut {
            str = str[ : length - len(cut)]
        } else {
            break
        }
    }
    return str
}
