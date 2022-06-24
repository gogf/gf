// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
)

// IsGNUVersion checks and returns whether given `version` is valid GNU version string.
func IsGNUVersion(version string) bool {
	if version != "" && (version[0] == 'v' || version[0] == 'V') {
		version = version[1:]
	}
	if version == "" {
		return false
	}
	var array = strings.Split(version, ".")
	if len(array) > 3 {
		return false
	}
	for _, v := range array {
		if v == "" {
			return false
		}
		if !IsNumeric(v) {
			return false
		}
		if v[0] == '-' || v[0] == '+' {
			return false
		}
	}
	return true
}

// CompareVersion compares `a` and `b` as standard GNU version.
// It returns  1 if `a` > `b`.
// It returns -1 if `a` < `b`.
// It returns  0 if `a` = `b`.
// GNU standard version is like:
// v1.0
// 1
// 1.0.0
// v1.0.1
// v2.10.8
// 10.2.0
// etc.
func CompareVersion(a, b string) int {
	if a != "" && a[0] == 'v' {
		a = a[1:]
	}
	if b != "" && b[0] == 'v' {
		b = b[1:]
	}
	var (
		array1 = strings.Split(a, ".")
		array2 = strings.Split(b, ".")
		diff   = 0
	)
	diff = len(array2) - len(array1)
	for i := 0; i < diff; i++ {
		array1 = append(array1, "0")
	}
	diff = len(array1) - len(array2)
	for i := 0; i < diff; i++ {
		array2 = append(array2, "0")
	}
	v1 := 0
	v2 := 0
	for i := 0; i < len(array1); i++ {
		v1 = gconv.Int(array1[i])
		v2 = gconv.Int(array2[i])
		if v1 > v2 {
			return 1
		}
		if v1 < v2 {
			return -1
		}
	}
	return 0
}

// CompareVersionGo compares `a` and `b` as standard Golang version.
// It returns  1 if `a` > `b`.
// It returns -1 if `a` < `b`.
// It returns  0 if `a` = `b`.
// Golang standard version is like:
// 1.0.0
// v1.0.1
// v2.10.8
// 10.2.0
// v0.0.0-20190626092158-b2ccc519800e
// v4.20.0+incompatible
// etc.
func CompareVersionGo(a, b string) int {
	a = Trim(a)
	b = Trim(b)
	if a != "" && a[0] == 'v' {
		a = a[1:]
	}
	if b != "" && b[0] == 'v' {
		b = b[1:]
	}
	var (
		rawA = a
		rawB = b
	)
	if Count(a, "-") > 1 {
		if i := PosR(a, "-"); i > 0 {
			a = a[:i]
		}
	}
	if Count(b, "-") > 1 {
		if i := PosR(b, "-"); i > 0 {
			b = b[:i]
		}
	}
	if i := Pos(a, "+"); i > 0 {
		a = a[:i]
	}
	if i := Pos(b, "+"); i > 0 {
		b = b[:i]
	}
	a = Replace(a, "-", ".")
	b = Replace(b, "-", ".")
	var (
		array1 = strings.Split(a, ".")
		array2 = strings.Split(b, ".")
		diff   = 0
	)
	// Specially in Golang:
	// "v1.12.2-0.20200413154443-b17e3a6804fa" < "v1.12.2"
	if len(array1) > 3 && len(array2) <= 3 {
		return -1
	}
	if len(array1) <= 3 && len(array2) > 3 {
		return 1
	}

	diff = len(array2) - len(array1)
	for i := 0; i < diff; i++ {
		array1 = append(array1, "0")
	}
	diff = len(array1) - len(array2)
	for i := 0; i < diff; i++ {
		array2 = append(array2, "0")
	}
	v1 := 0
	v2 := 0
	for i := 0; i < len(array1); i++ {
		v1 = gconv.Int(array1[i])
		v2 = gconv.Int(array2[i])
		if v1 > v2 {
			return 1
		}
		if v1 < v2 {
			return -1
		}
	}
	// Specially in Golang:
	// "v4.20.1+incompatible" < "v4.20.1"
	if Contains(rawA, "incompatible") {
		return -1
	}
	if Contains(rawB, "incompatible") {
		return 1
	}
	return 0
}
