// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

type localAttribute struct {
	key   string
	value any
}

func NewAttribute(key string, value any) Attribute {
	return &localAttribute{
		key:   key,
		value: value,
	}
}

func (l *localAttribute) Key() string {
	return l.key
}

func (l *localAttribute) Value() any {
	return l.value
}
