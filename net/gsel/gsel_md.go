// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

// MD is a mapping from metadata keys to value array.
// Users should use the following two convenience functions New and Pairs to generate MD.
type MD interface {
	// Len returns the number of items in md.
	Len() int

	// Get obtains the values for a given key.
	//
	// k is converted to lowercase before searching in md.
	Get(k string) []string

	// Set sets the value of a given key with a slice of values.
	//
	// k is converted to lowercase before storing in md.
	Set(key string, values ...string)

	// Append adds the values to key k, not overwriting what was already stored at
	// that key.
	//
	// k is converted to lowercase before storing in md.
	Append(k string, values ...string)

	// Delete removes the values for a given key k which is converted to lowercase
	// before removing it from md.
	Delete(k string)
}
