// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Byte converts `any` to byte.
func Byte(anyInput any) byte {
	v, _ := defaultConverter.Uint8(anyInput)
	return v
}

// Bytes converts `any` to []byte.
func Bytes(anyInput any) []byte {
	v, _ := defaultConverter.Bytes(anyInput)
	return v
}

// Rune converts `any` to rune.
func Rune(anyInput any) rune {
	v, _ := defaultConverter.Rune(anyInput)
	return v
}

// Runes converts `any` to []rune.
func Runes(anyInput any) []rune {
	v, _ := defaultConverter.Runes(anyInput)
	return v
}

// String converts `any` to string.
// It's most commonly used converting function.
func String(anyInput any) string {
	v, _ := defaultConverter.String(anyInput)
	return v
}

// Bool converts `any` to bool.
// It returns false if `any` is: false, "", 0, "false", "off", "no", empty slice/map.
func Bool(anyInput any) bool {
	v, _ := defaultConverter.Bool(anyInput)
	return v
}
