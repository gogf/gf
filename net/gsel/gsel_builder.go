// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

// defaultBuilder is the default Builder for globally used purpose.
var defaultBuilder = NewBuilderRoundRobin()

// SetBuilder sets the default builder for globally used purpose.
func SetBuilder(builder Builder) {
	defaultBuilder = builder
}

// GetBuilder returns the default builder for globally used purpose.
func GetBuilder() Builder {
	return defaultBuilder
}
