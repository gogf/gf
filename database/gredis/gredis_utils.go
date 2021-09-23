// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

func defaultOption() *Option {
	return &Option{}
}

func parseArgsToArgsAndOption(arguments []interface{}) (args []interface{}, option *Option) {
	if length := len(arguments); length > 0 {
		lastElement := arguments[length-1]
		switch result := lastElement.(type) {
		case Option:
			return arguments[:length-1], &result

		case *Option:
			return arguments[:length-1], result
		}
	}
	return arguments, defaultOption()
}
