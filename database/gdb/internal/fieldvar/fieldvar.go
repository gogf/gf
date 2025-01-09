// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package fieldvar provides a variable wrapper for field value in database operations.
// It is used for field value storage and conversion in database operations.
package fieldvar

import (
	"github.com/gogf/gf/v2/database/gdb/internal/defines"
	"github.com/gogf/gf/v2/internal/json"
)

// Var is a wrapper for any type of value, which is used for field variable.
// Note that, do not embed *gvar.Var into Var but use it as an attribute, as there issue in nil pointer receiver
// when calling methods that is not defined directly on Var.
type Var struct {
	value     any
	localType defines.LocalType
}

// New creates and returns a new Var object.
func New(value any) *Var {
	return &Var{
		value: value,
	}
}

// NewWithType creates and returns a new Var object with specified local type.
func NewWithType(value any, localType defines.LocalType) *Var {
	return &Var{
		value:     value,
		localType: localType,
	}
}

// LocalType returns the local type of the value.
func (v *Var) LocalType() defines.LocalType {
	return v.localType
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Var) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Val())
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Var) UnmarshalJSON(b []byte) error {
	var i interface{}
	if err := json.UnmarshalUseNumber(b, &i); err != nil {
		return err
	}
	v.Set(i)
	return nil
}

// UnmarshalValue is an interface implement which sets any type of value for Var.
func (v *Var) UnmarshalValue(value interface{}) error {
	v.Set(value)
	return nil
}
