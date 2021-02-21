// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gset

// NewImmutableIntSet would create a immutable int set.
func NewImmutableIntSet(items ...int) ImmutableIntSet {
	return NewIntSetFrom(items, false)
}

// ImmutableIntSet wrap the IntSet and expose the read method.
type ImmutableIntSet interface {

	// Iterator iterates the set readonly with given callback function <f>,
	// if <f> returns true then continue iterating; or false to stop.
	Iterator(f func(v int) bool)

	// Contains checks whether the set contains <item>.
	Contains(item int) bool

	// Size returns the size of the set.
	Size() int

	// Slice returns the a of items of the set as slice.
	Slice() []int

	// Join joins items with a string <glue>.
	Join(glue string) string

	// String returns items as a string, which implements like json.Marshal does.
	String() string

	// Equal checks whether the two sets equal.
	Equal(other *IntSet) bool

	// IsSubsetOf checks whether the current set is a sub-set of <other>.
	IsSubsetOf(other *IntSet) bool

	// Union returns a new set which is the union of <set> and <other>.
	// Which means, all the items in <newSet> are in <set> or in <other>.
	Union(others ...*IntSet) *IntSet

	// Diff returns a new set which is the difference set from <set> to <other>.
	// Which means, all the items in <newSet> are in <set> but not in <other>.
	Diff(others ...*IntSet) *IntSet

	// Intersect returns a new set which is the intersection from <set> to <other>.
	// Which means, all the items in <newSet> are in <set> and also in <other>.
	Intersect(others ...*IntSet) *IntSet

	// Complement returns a new set which is the complement from <set> to <full>.
	// Which means, all the items in <newSet> are in <full> and not in <set>.
	//
	// It returns the difference between <full> and <set>
	// if the given set <full> is not the full set of <set>.
	Complement(full *IntSet) *IntSet

	// Sum sums items.
	// Note: The items should be converted to int type,
	// or you'd get a result that you unexpected.
	Sum() int

	// MarshalJSON implements the interface MarshalJSON for json.Marshal.
	MarshalJSON() ([]byte, error)

	// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
	UnmarshalJSON(b []byte) error

	// UnmarshalValue is an interface implement which sets any type of value for set.
	UnmarshalValue(value interface{}) error
}
