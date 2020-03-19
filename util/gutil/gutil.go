// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gutil provides utility functions.
package gutil

import (
	"github.com/gogf/gf/internal/empty"
)

// Throw throws out an exception, which can be caught be TryCatch or recover.
func Throw(exception interface{}) {
	panic(exception)
}

// TryCatch implements try...catch... logistics.
func TryCatch(try func(), catch ...func(exception interface{})) {
	if len(catch) > 0 {
		// If <catch> is given, it's used to handle the exception.
		defer func() {
			if e := recover(); e != nil {
				catch[0](e)
			}
		}()
	} else {
		// If no <catch> function passed, it filters the exception.
		defer func() {
			recover()
		}()
	}
	try()
}

// IsEmpty checks given <value> empty or not.
// It returns false if <value> is: integer(0), bool(false), slice/map(len=0), nil;
// or else returns true.
func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}
