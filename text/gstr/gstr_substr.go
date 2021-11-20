// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

// SubStr returns a portion of string `str` specified by the `start` and `length` parameters.
// The parameter `length` is optional, it uses the length of `str` in default.
func SubStr(str string, start int, length ...int) (substr string) {
	strLength := len(str)
	// Simple border checks.
	if start < 0 {
		start = 0
	}
	if start >= strLength {
		start = strLength
	}
	end := strLength
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = strLength
		}
	}
	if end > strLength {
		end = strLength
	}
	return str[start:end]
}

// SubStrRune returns a portion of string `str` specified by the `start` and `length` parameters.
// SubStrRune considers parameter `str` as unicode string.
// The parameter `length` is optional, it uses the length of `str` in default.
func SubStrRune(str string, start int, length ...int) (substr string) {
	// Converting to []rune to support unicode.
	var (
		runes       = []rune(str)
		runesLength = len(runes)
	)

	// Simple border checks.
	if start < 0 {
		start = 0
	}
	if start >= runesLength {
		start = runesLength
	}
	end := runesLength
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = runesLength
		}
	}
	if end > runesLength {
		end = runesLength
	}
	return string(runes[start:end])
}

// StrLimit returns a portion of string `str` specified by `length` parameters, if the length
// of `str` is greater than `length`, then the `suffix` will be appended to the result string.
func StrLimit(str string, length int, suffix ...string) string {
	if len(str) < length {
		return str
	}
	suffixStr := defaultSuffixForStrLimit
	if len(suffix) > 0 {
		suffixStr = suffix[0]
	}
	return str[0:length] + suffixStr
}

// StrLimitRune returns a portion of string `str` specified by `length` parameters, if the length
// of `str` is greater than `length`, then the `suffix` will be appended to the result string.
// StrLimitRune considers parameter `str` as unicode string.
func StrLimitRune(str string, length int, suffix ...string) string {
	runes := []rune(str)
	if len(runes) < length {
		return str
	}
	suffixStr := defaultSuffixForStrLimit
	if len(suffix) > 0 {
		suffixStr = suffix[0]
	}
	return string(runes[0:length]) + suffixStr
}

// SubStrFrom returns a portion of string `str` starting from first occurrence of and including `need`
// to the end of `str`.
func SubStrFrom(str string, need string) (substr string) {
	pos := Pos(str, need)
	if pos < 0 {
		return ""
	}
	return str[pos:]
}

// SubStrFromEx returns a portion of string `str` starting from first occurrence of and excluding `need`
// to the end of `str`.
func SubStrFromEx(str string, need string) (substr string) {
	pos := Pos(str, need)
	if pos < 0 {
		return ""
	}
	return str[pos+len(need):]
}

// SubStrFromR returns a portion of string `str` starting from last occurrence of and including `need`
// to the end of `str`.
func SubStrFromR(str string, need string) (substr string) {
	pos := PosR(str, need)
	if pos < 0 {
		return ""
	}
	return str[pos:]
}

// SubStrFromREx returns a portion of string `str` starting from last occurrence of and excluding `need`
// to the end of `str`.
func SubStrFromREx(str string, need string) (substr string) {
	pos := PosR(str, need)
	if pos < 0 {
		return ""
	}
	return str[pos+len(need):]
}
