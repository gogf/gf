// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype

import (
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shopspring/decimal"
)

// Decimal is a struct for concurrent-safe operation for type decimal.
type Decimal struct {
	value decimal.Decimal
}

// NewDecimal creates and returns a concurrent-safe object for decimal type,
// with given initial value <value>.
func NewDecimal(value ...decimal.Decimal) *Decimal {
	if len(value) > 0 {
		return &Decimal{
			value: value[0],
		}
	}
	return &Decimal{}
}

// Clone clones and returns a new concurrent-safe object for decimal type.
func (v *Decimal) Clone() *Decimal {
	return NewDecimal(v.Val())
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *Decimal) Set(value decimal.Decimal) (old decimal.Decimal) {
	old = v.value
	v.value = v.value.Add(value)
	return
}

// Val atomically loads and returns t.value.
func (v *Decimal) Val() decimal.Decimal {
	return v.value
}

// Add atomically adds <delta> to t.value and returns the new value.
func (v *Decimal) Add(delta decimal.Decimal) (new decimal.Decimal) {
	new = v.value.Add(delta)
	return
}

// Cas executes the compare-and-swap operation for value.
func (v *Decimal) Cas(old, new decimal.Decimal) (swapped bool) {
	if old.Cmp(new) < 0 {
		v.value = new
		return true
	}
	return false
}

// String implements String interface for string printing.
func (v *Decimal) String() string {
	return v.value.String()
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (v *Decimal) MarshalJSON() ([]byte, error) {
	return v.value.MarshalJSON()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (v *Decimal) UnmarshalJSON(b []byte) error {
	return v.value.UnmarshalJSON(b)
}

// UnmarshalValue is an interface implement which sets any type of value for <v>.
func (v *Decimal) UnmarshalValue(value interface{}) error {
	v.Set(gconv.Decimal(value))
	return nil
}
