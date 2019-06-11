// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import "strings"

// Trim strips whitespace (or other characters) from the beginning and end of a string.
func Trim(str string, characterMask ...string) string {
    if len(characterMask) > 0 {
        return strings.Trim(str, characterMask[0])
    } else {
        return strings.TrimSpace(str)
    }
}

// TrimLeft strips whitespace (or other characters) from the beginning of a string.
func TrimLeft(str string, characterMask ...string) string {
    mask := ""
    if len(characterMask) == 0 {
        mask = string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0})
    } else {
        mask = characterMask[0]
    }
    return strings.TrimLeft(str, mask)
}

// TrimLeftStr strips all of the given <cut> string from the beginning of a string.
func TrimLeftStr(str string, cut string) string {
    for str[0 : len(cut)] == cut {
        str = str[len(cut) : ]
    }
    return str
}

// TrimRight strips whitespace (or other characters) from the end of a string.
func TrimRight(str string, characterMask ...string) string {
    mask := ""
    if len(characterMask) == 0 {
        mask = string([]byte{'\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0})
    } else {
        mask = characterMask[0]
    }
    return strings.TrimRight(str, mask)
}

// TrimRightStr strips all of the given <cut> string from the end of a string.
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
