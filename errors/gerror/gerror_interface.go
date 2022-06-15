// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import (
	"github.com/gogf/gf/v2/errors/gcode"
)

// iIs is the interface for Is feature.
type iIs interface {
	Is(target error) bool
}

// iEqual is the interface for Equal feature.
type iEqual interface {
	Equal(target error) bool
}

// iCode is the interface for Code feature.
type iCode interface {
	Error() string
	Code() gcode.Code
}

// iStack is the interface for Stack feature.
type iStack interface {
	Error() string
	Stack() string
}

// iCause is the interface for Cause feature.
type iCause interface {
	Error() string
	Cause() error
}

// iCurrent is the interface for Current feature.
type iCurrent interface {
	Error() string
	Current() error
}

// iNext is the interface for Next feature.
type iNext interface {
	Error() string
	Next() error
}

// iUnwrap is the interface for Unwrap feature.
type iUnwrap interface {
	Error() string
	Unwrap() error
}
