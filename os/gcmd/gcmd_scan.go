// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import "fmt"

// Scan prints <info> to stdout, reads and returns user input, which stops by '\n'.
func Scan(info ...interface{}) string {
	var s string
	fmt.Print(info...)
	fmt.Scanln(&s)
	return s
}

// Scanf prints <info> to stdout with <format>, reads and returns user input, which stops by '\n'.
func Scanf(format string, info ...interface{}) string {
	var s string
	fmt.Printf(format, info...)
	fmt.Scanln(&s)
	return s
}
