// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genv

// MustSet performs as Set, but it panics if any error occurs.
func MustSet(key, value string) {
	if err := Set(key, value); err != nil {
		panic(err)
	}
}

// MustRemove performs as Remove, but it panics if any error occurs.
func MustRemove(key ...string) {
	if err := Remove(key...); err != nil {
		panic(err)
	}
}
