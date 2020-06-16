// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gutil"
)

// NewVar returns a gvar.Var.
func NewVar(i interface{}, safe ...bool) Var {
	return gvar.New(i, safe...)
}

// Wait blocks until all the web servers shutdown.
func Wait() {
	ghttp.Wait()
}

// Dump dumps a variable to stdout with more manually readable.
func Dump(i ...interface{}) {
	gutil.Dump(i...)
}

// Export exports a variable to string with more manually readable.
func Export(i ...interface{}) string {
	return gutil.Export(i...)
}

// Throw throws a exception, which can be caught by TryCatch function.
// It always be used in TryCatch function.
func Throw(exception interface{}) {
	gutil.Throw(exception)
}

// TryCatch does the try...catch... mechanism.
func TryCatch(try func(), catch ...func(exception interface{})) {
	gutil.TryCatch(try, catch...)
}

// IsNil checks whether given <value> is nil.
// Note that it might use reflect feature which affects performance a little bit.
func IsNil(value interface{}) bool {
	return empty.IsNil(value)
}

// IsEmpty checks whether given <value> empty.
// It returns true if <value> is in: 0, nil, false, "", len(slice/map/chan) == 0.
// Or else it returns true.
func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}
