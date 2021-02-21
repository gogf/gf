// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"github.com/gogf/gf/container/gvar"
)

// NewImmutableAnyAnyMap would create a immutable map from the given map.
func NewImmutableAnyAnyMap(data map[interface{}]interface{}) ImmutableAnyAnyMap {
	mm := make(map[interface{}]interface{}, len(data))
	for k, v := range data {
		mm[k] = v
	}

	return NewAnyAnyMapFrom(mm, false)
}

// ImmutableAnyAnyMap wrap the AnyAnyMap and expose the read function.
type ImmutableAnyAnyMap interface {

	// Iterator iterates the hash map readonly with custom callback function <f>.
	// If <f> returns true, then it continues iterating; or false to stop.
	Iterator(f func(k interface{}, v interface{}) bool)

	// MapCopy returns a copy of underlying data map.
	MapCopy() map[interface{}]interface{}

	// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
	MapStrAny() map[string]interface{}

	// Search searches the map with given <key>.
	// Second return parameter <found> is true if key was found, otherwise false.
	Search(key interface{}) (interface{}, bool)

	// Get returns the value by given <key>.
	Get(key interface{}) interface{}

	// GetVar returns a Var with the value by given <key>.
	// The returned Var is un-concurrent safe.
	GetVar(key interface{}) *gvar.Var

	// Keys returns all keys of the map as a slice.
	Keys() []interface{}

	// Values returns all values of the map as a slice.
	Values() []interface{}

	// Contains checks whether a key exists.
	// It returns true if the <key> exists, or else false.
	Contains(key interface{}) bool

	// Size returns the size of the map.
	Size() int

	// IsEmpty checks whether the map is empty.
	// It returns true if map is empty, or else false.
	IsEmpty() bool

	// String returns the map as a string.
	String() string

	// MarshalJSON implements the interface MarshalJSON for json.Marshal.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
	UnmarshalJSON(b []byte) error

	// UnmarshalValue is an interface implement which sets any type of value for map.
	UnmarshalValue(value interface{}) error
}
