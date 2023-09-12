// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstructs

// Signature returns a unique string as this type.
func (t Type) Signature() string {
	return t.PkgPath() + "/" + t.String()
}

// FieldKeys returns the keys of current struct/map.
func (t Type) FieldKeys() []string {
	keys := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		keys[i] = t.Field(i).Name
	}
	return keys
}
