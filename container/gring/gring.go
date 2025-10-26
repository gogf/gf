// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gring provides a concurrent-safe/unsafe ring(circular lists).
//
// Deprecated.
package gring

// Ring is a struct of ring structure.
//
// Deprecated.
type Ring struct {
	*TRing[any]
}

// New creates and returns a Ring structure of `cap` elements.
// The optional parameter `safe` specifies whether using this structure in concurrent safety,
// which is false in default.
//
// Deprecated.
func New(cap int, safe ...bool) *Ring {
	return &Ring{
		TRing: NewTRing[any](cap, safe...),
	}
}

// Val returns the item's value of current position.
func (r *Ring) Val() any {
	return r.TRing.Val()
}

// Len returns the size of ring.
func (r *Ring) Len() int {
	return r.TRing.Len()
}

// Cap returns the capacity of ring.
func (r *Ring) Cap() int {
	return r.TRing.Cap()
}

// Set sets value to the item of current position.
func (r *Ring) Set(value any) *Ring {
	r.TRing.Set(value)
	return r
}

// Put sets `value` to current item of ring and moves position to next item.
func (r *Ring) Put(value any) *Ring {
	r.TRing.Put(value)
	return r
}

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ring and returns that ring element. r must not be empty.
func (r *Ring) Move(n int) *Ring {
	r.TRing.Move(n)
	return r
}

// Prev returns the previous ring element. r must not be empty.
func (r *Ring) Prev() *Ring {
	r.TRing.Prev()
	return r
}

// Next returns the next ring element. r must not be empty.
func (r *Ring) Next() *Ring {
	r.TRing.Next()
	return r
}

// Link connects ring r with ring s such that r.Next()
// becomes s and returns the original value for r.Next().
// r must not be empty.
//
// If r and s point to the same ring, linking
// them removes the elements between r and s from the ring.
// The removed elements form a sub-ring and the result is a
// reference to that sub-ring (if no elements were removed,
// the result is still the original value for r.Next(),
// and not nil).
//
// If r and s point to different rings, linking
// them creates a single ring with the elements of s inserted
// after r. The result points to the element following the
// last element of s after insertion.
func (r *Ring) Link(s *Ring) *Ring {
	r.TRing.Link(s.TRing)
	return r
}

// Unlink removes n % r.Len() elements from the ring r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed sub-ring. r must not be empty.
func (r *Ring) Unlink(n int) *Ring {
	return &Ring{
		TRing: r.TRing.Unlink(n),
	}
}

// RLockIteratorNext iterates and locks reading forward
// with given callback function `f` within RWMutex.RLock.
// If `f` returns true, then it continues iterating; or false to stop.
func (r *Ring) RLockIteratorNext(f func(value any) bool) {
	r.TRing.RLockIteratorNext(f)
}

// RLockIteratorPrev iterates and locks writing backward
// with given callback function `f` within RWMutex.RLock.
// If `f` returns true, then it continues iterating; or false to stop.
func (r *Ring) RLockIteratorPrev(f func(value any) bool) {
	r.TRing.RLockIteratorPrev(f)
}

// SliceNext returns a copy of all item values as slice forward from current position.
func (r *Ring) SliceNext() []any {
	return r.TRing.SliceNext()
}

// SlicePrev returns a copy of all item values as slice backward from current position.
func (r *Ring) SlicePrev() []any {
	return r.TRing.SlicePrev()
}
