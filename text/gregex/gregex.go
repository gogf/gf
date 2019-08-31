// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gregex provides high performance API for regular expression functionality.
package gregex

import (
	"regexp"
)

// Quote quotes <s> by replacing special chars in <s>
// to match the rules of regular expression pattern.
// And returns the copy.
//
// Eg: Quote(`[foo]`) returns `\[foo\]`.
func Quote(s string) string {
	return regexp.QuoteMeta(s)
}

// Validate checks whether given regular expression pattern <pattern> valid.
func Validate(pattern string) error {
	_, err := getRegexp(pattern)
	return err
}

// IsMatch checks whether given bytes <src> matches <pattern>.
func IsMatch(pattern string, src []byte) bool {
	if r, err := getRegexp(pattern); err == nil {
		return r.Match(src)
	}
	return false
}

// IsMatchString checks whether given string <src> matches <pattern>.
func IsMatchString(pattern string, src string) bool {
	return IsMatch(pattern, []byte(src))
}

// MatchString return bytes slice that matched <pattern>.
func Match(pattern string, src []byte) ([][]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindSubmatch(src), nil
	} else {
		return nil, err
	}
}

// MatchString return strings that matched <pattern>.
func MatchString(pattern string, src string) ([]string, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindStringSubmatch(src), nil
	} else {
		return nil, err
	}
}

// MatchAll return all bytes slices that matched <pattern>.
func MatchAll(pattern string, src []byte) ([][][]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindAllSubmatch(src, -1), nil
	} else {
		return nil, err
	}
}

// MatchAllString return all strings that matched <pattern>.
func MatchAllString(pattern string, src string) ([][]string, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.FindAllStringSubmatch(src, -1), nil
	} else {
		return nil, err
	}
}

// ReplaceString replace all matched <pattern> in bytes <src> with bytes <replace>.
func Replace(pattern string, replace, src []byte) ([]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.ReplaceAll(src, replace), nil
	} else {
		return nil, err
	}
}

// ReplaceString replace all matched <pattern> in string <src> with string <replace>.
func ReplaceString(pattern, replace, src string) (string, error) {
	r, e := Replace(pattern, []byte(replace), []byte(src))
	return string(r), e
}

// ReplaceFunc replace all matched <pattern> in bytes <src>
// with custom replacement function <replaceFunc>.
func ReplaceFunc(pattern string, src []byte, replaceFunc func(b []byte) []byte) ([]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.ReplaceAllFunc(src, replaceFunc), nil
	} else {
		return nil, err
	}
}

// ReplaceFuncMatch replace all matched <pattern> in bytes <src>
// with custom replacement function <replaceFunc>.
// The parameter <match> type for <replaceFunc> is [][]byte,
// which is the result contains all sub-patterns of <pattern> using Match function.
func ReplaceFuncMatch(pattern string, src []byte, replaceFunc func(match [][]byte) []byte) ([]byte, error) {
	if r, err := getRegexp(pattern); err == nil {
		return r.ReplaceAllFunc(src, func(bytes []byte) []byte {
			match, _ := Match(pattern, bytes)
			return replaceFunc(match)
		}), nil
	} else {
		return nil, err
	}
}

// ReplaceStringFunc replace all matched <pattern> in string <src>
// with custom replacement function <replaceFunc>.
func ReplaceStringFunc(pattern string, src string, replaceFunc func(s string) string) (string, error) {
	bytes, err := ReplaceFunc(pattern, []byte(src), func(bytes []byte) []byte {
		return []byte(replaceFunc(string(bytes)))
	})
	return string(bytes), err
}

// ReplaceStringFuncMatch replace all matched <pattern> in string <src>
// with custom replacement function <replaceFunc>.
// The parameter <match> type for <replaceFunc> is []string,
// which is the result contains all sub-patterns of <pattern> using MatchString function.
func ReplaceStringFuncMatch(pattern string, src string, replaceFunc func(match []string) string) (string, error) {
	if r, err := getRegexp(pattern); err == nil {
		return string(r.ReplaceAllFunc([]byte(src), func(bytes []byte) []byte {
			match, _ := MatchString(pattern, string(bytes))
			return []byte(replaceFunc(match))
		})), nil
	} else {
		return "", err
	}
}

// Split slices <src> into substrings separated by the expression and returns a slice of
// the substrings between those expression matches.
func Split(pattern string, src string) []string {
	if r, err := getRegexp(pattern); err == nil {
		return r.Split(src, -1)
	}
	return nil
}
