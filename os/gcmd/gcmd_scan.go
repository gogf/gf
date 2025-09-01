// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/text/gstr"
)

// Scan prints `info` to stdout, reads and returns user input, which stops by '\n'.
func Scan(info ...any) string {
	fmt.Print(info...)
	return readline()
}

// Scanf prints `info` to stdout with `format`, reads and returns user input, which stops by '\n'.
func Scanf(format string, info ...any) string {
	fmt.Printf(format, info...)
	return readline()
}

func readline() string {
	var s string
	reader := bufio.NewReader(os.Stdin)
	s, _ = reader.ReadString('\n')
	s = gstr.Trim(s)
	return s
}
