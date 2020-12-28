// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror_test

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
)

func ExampleNewCode() {
	err := gerror.NewCode(10000, "My Error")
	fmt.Println(err.Error())
	fmt.Println(gerror.Code(err))

	// Output:
	// My Error
	// 10000
}

func ExampleNewCodef() {
	err := gerror.NewCodef(10000, "It's %s", "My Error")
	fmt.Println(err.Error())
	fmt.Println(gerror.Code(err))

	// Output:
	// It's My Error
	// 10000
}

func ExampleWrapCode() {
	err1 := errors.New("permission denied")
	err2 := gerror.WrapCode(10000, err1, "Custom Error")
	fmt.Println(err2.Error())
	fmt.Println(gerror.Code(err2))

	// Output:
	// Custom Error: permission denied
	// 10000
}

func ExampleWrapCodef() {
	err1 := errors.New("permission denied")
	err2 := gerror.WrapCodef(10000, err1, "It's %s", "Custom Error")
	fmt.Println(err2.Error())
	fmt.Println(gerror.Code(err2))

	// Output:
	// It's Custom Error: permission denied
	// 10000
}
