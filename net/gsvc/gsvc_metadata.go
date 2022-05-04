// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"github.com/gogf/gf/v2/container/gvar"
)

// Set sets key-value pair into metadata.
func (m Metadata) Set(key string, value string) {
	m[key] = value
}

// Get retrieves and return value of specified key as gvar.
func (m Metadata) Get(key string) *gvar.Var {
	if v, ok := m[key]; ok {
		return gvar.New(v)
	}

	return nil
}
