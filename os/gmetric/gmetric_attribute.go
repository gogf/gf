// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"fmt"
)

// localAttribute implements interface Attribute.
type localAttribute struct {
	key   string
	value any
}

// NewAttribute creates and returns an Attribute by given `key` and `value`.
func NewAttribute(key string, value any) Attribute {
	return &localAttribute{
		key:   key,
		value: value,
	}
}

// Key returns the key of the attribute.
func (l *localAttribute) Key() string {
	return l.key
}

// Value returns the value of the attribute.
func (l *localAttribute) Value() any {
	return l.value
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (l *localAttribute) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{"%s":%#v}`, l.key, l.value)), nil
}
