// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gmap

import (
	"github.com/gogf/gf/container/gvar"
)

// NewImmutableListMap would create a immutable map from the given map.
func NewImmutableListMap(data map[interface{}]interface{}) ImmutableListMap {
	return NewListMapFrom(data, false)
}

// ImmutableListMap wrap the ListMap and expose the read function.
type ImmutableListMap interface {

	// Iterator is alias of IteratorAsc.
	Iterator(f func(key, value interface{}) bool)

	// IteratorAsc iterates the map readonly in ascending order with given callback function <f>.
	// If <f> returns true, then it continues iterating; or false to stop.
	IteratorAsc(f func(key interface{}, value interface{}) bool)

	// IteratorDesc iterates the map readonly in descending order with given callback function <f>.
	// If <f> returns true, then it continues iterating; or false to stop.
	IteratorDesc(f func(key interface{}, value interface{}) bool)

	// Map returns a copy of the underlying data of the map.
	Map() map[interface{}]interface{}

	// MapStrAny returns a copy of the underlying data of the map as map[string]interface{}.
	MapStrAny() map[string]interface{}

	// Search searches the map with given <key>.
	// Second return parameter <found> is true if key was found, otherwise false.
	Search(key interface{}) (value interface{}, found bool)

	// Get returns the value by given <key>.
	Get(key interface{}) (value interface{})

	// GetVar returns a Var with the value by given <key>.
	// The returned Var is un-concurrent safe.
	GetVar(key interface{}) *gvar.Var

	// Keys returns all keys of the map as a slice in ascending order.
	Keys() []interface{}

	// Values returns all values of the map as a slice.
	Values() []interface{}

	// Contains checks whether a key exists.
	// It returns true if the <key> exists, or else false.
	Contains(key interface{}) (ok bool)

	// Size returns the size of the map.
	Size() (size int)

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
	UnmarshalValue(value interface{}) (err error)
}
