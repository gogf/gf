// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstr

import (
	"regexp"
	"strconv"
)

var (
	// octReg is the regular expression object for checks octal string.
	octReg = regexp.MustCompile(`\\[0-7]{3}`)
)

// OctStr converts string container octal string to its original string,
// for example, to Chinese string.
// Eg: `\346\200\241` -> 怡
func OctStr(str string) string {
	return octReg.ReplaceAllStringFunc(
		str,
		func(s string) string {
			i, _ := strconv.ParseInt(s[1:], 8, 0)
			return string([]byte{byte(i)})
		},
	)
}
