// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "github.com/gogf/gf/v2/database/gdb/internal/fieldvar"

// NewValue creates and returns a new Value object.
func NewValue(value any) Value {
    return fieldvar.New(value)
}

// NewValueWithType creates and returns a new Value object with specified local type.
func NewValueWithType(value any, localType LocalType) Value {
    return fieldvar.NewWithType(value, localType)
}
