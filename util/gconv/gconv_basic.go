// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

// Byte converts `any` to byte.
func Byte(any any) byte {
	v, _ := defaultConverter.Uint8(any)
	return v
}

// Bytes converts `any` to []byte.
func Bytes(any any) []byte {
	v, _ := defaultConverter.Bytes(any)
	return v
}

// Rune converts `any` to rune.
func Rune(any any) rune {
	v, _ := defaultConverter.Rune(any)
	return v
}

// Runes converts `any` to []rune.
func Runes(any any) []rune {
	v, _ := defaultConverter.Runes(any)
	return v
}

// String converts `any` to string.
// It's most commonly used converting function.
func String(any any) string {
	v, _ := defaultConverter.String(any)
	return v
}

// Bool converts `any` to bool.
// It returns false if `any` is: false, "", 0, "false", "off", "no", empty slice/map.
func Bool(any any) bool {
	v, _ := defaultConverter.Bool(any)
	return v
}
