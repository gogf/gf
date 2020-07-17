// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"github.com/gogf/gf/util/gconv"
	"strings"
)

// CompareVersion compares <a> and <b> as standard GNU version.
// It returns  1 if <a> > <b>.
// It returns -1 if <a> < <b>.
// It returns  0 if <a> = <b>.
// GNU standard version is like:
// v1.0
// 1
// 1.0.0
// v1.0.1
// v2.10.8
// 10.2.0
// etc.
func CompareVersion(a, b string) int {
	if a[0] == 'v' {
		a = a[1:]
	}
	if b[0] == 'v' {
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

// CompareVersionGo compares <a> and <b> as standard Golang version.
// It returns  1 if <a> > <b>.
// It returns -1 if <a> < <b>.
// It returns  0 if <a> = <b>.
// Golang standard version is like:
// 1.0.0
// v1.0.1
// v2.10.8
// 10.2.0
// v0.0.0-20190626092158-b2ccc519800e
// v4.20.0+incompatible
// etc.
func CompareVersionGo(a, b string) int {
	if a[0] == 'v' {
		a = a[1:]
	}
	if b[0] == 'v' {
		b = b[1:]
	}
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
	return 0
}
