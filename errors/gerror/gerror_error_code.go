// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import (
	"github.com/gogf/gf/v2/errors/gcode"
)

// Code returns the error code.
// It returns CodeNil if it has no error code.
func (err *Error) Code() gcode.Code {
	if err == nil {
		return gcode.CodeNil
	}
	if err.code == gcode.CodeNil {
		return Code(err.Next())
	}
	return err.code
}

// SetCode updates the internal code with given code.
func (err *Error) SetCode(code gcode.Code) {
	if err == nil {
		return
	}
	err.code = code
}
