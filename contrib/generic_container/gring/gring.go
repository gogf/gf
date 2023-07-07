// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gring provides a concurrent-safe/unsafe ring(circular lists).
package gring

import (
	"container/ring"

	"github.com/gogf/gf/contrib/generic_container/v2/internal/rwmutex"
	"github.com/gogf/gf/v2/container/gtype"
)

// Ring is a struct of ring structure.
type Ring[T any] struct {
	mu    *rwmutex.RWMutex
	ring  *ring.Ring  // Underlying ring.
	len   *gtype.Int  // Length(already used size).
	cap   *gtype.Int  // Capability(>=len).
	dirty *gtype.Bool // Dirty, which means the len and cap should be recalculated. It's marked dirty when the size of ring changes.
}

// internalRingItem stores the ring element value.
type internalRingItem[T any] struct {
	Value T
}

// New creates and returns a Ring structure of `cap` elements.
// The optional parameter `safe` specifies whether using this structure in concurrent safety,
// which is false in default.
func New[T any](cap int, safe ...bool) *Ring[T] {
	return &Ring[T]{
		mu:    rwmutex.New(safe...),
		ring:  ring.New(cap),
		len:   gtype.NewInt(),
		cap:   gtype.NewInt(cap),
		dirty: gtype.NewBool(),
	}
}

// Val returns the item's value of current position.
func (r *Ring[T]) Val() T {
	var value T
	r.mu.RLock()
	if r.ring.Value != nil {
		value = r.ring.Value.(internalRingItem[T]).Value
	}
	r.mu.RUnlock()
	return value
}

// Len returns the size of ring.
func (r *Ring[T]) Len() int {
	r.checkAndUpdateLenAndCap()
	return r.len.Val()
}

// Cap returns the capacity of ring.
func (r *Ring[T]) Cap() int {
	r.checkAndUpdateLenAndCap()
	return r.cap.Val()
}

// Checks and updates the len and cap of ring when ring is dirty.
func (r *Ring[T]) checkAndUpdateLenAndCap() {
	if !r.dirty.Val() {
		return
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	totalLen := 0
	emptyLen := 0
	if r.ring != nil {
		if r.ring.Value == nil {
			emptyLen++
		}
		totalLen++
		for p := r.ring.Next(); p != r.ring; p = p.Next() {
			if p.Value == nil {
				emptyLen++
			}
			totalLen++
		}
	}
	r.cap.Set(totalLen)
	r.len.Set(totalLen - emptyLen)
	r.dirty.Set(false)
}

// Set sets value to the item of current position.
func (r *Ring[T]) Set(value T) *Ring[T] {
	r.mu.Lock()
	if r.ring.Value == nil {
		r.len.Add(1)
	}
	r.ring.Value = internalRingItem[T]{Value: value}
	r.mu.Unlock()
	return r
}

// Put sets `value` to current item of ring and moves position to next item.
func (r *Ring[T]) Put(value T) *Ring[T] {
	r.mu.Lock()
	if r.ring.Value == nil {
		r.len.Add(1)
	}
	r.ring.Value = internalRingItem[T]{Value: value}
	r.ring = r.ring.Next()
	r.mu.Unlock()
	return r
}

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ring and returns that ring element. r must not be empty.
func (r *Ring[T]) Move(n int) *Ring[T] {
	r.mu.Lock()
	r.ring = r.ring.Move(n)
	r.mu.Unlock()
	return r
}

// Prev returns the previous ring element. r must not be empty.
func (r *Ring[T]) Prev() *Ring[T] {
	r.mu.Lock()
	r.ring = r.ring.Prev()
	r.mu.Unlock()
	return r
}

// Next returns the next ring element. r must not be empty.
func (r *Ring[T]) Next() *Ring[T] {
	r.mu.Lock()
	r.ring = r.ring.Next()
	r.mu.Unlock()
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
func (r *Ring[T]) Link(s *Ring[T]) *Ring[T] {
	r.mu.Lock()
	s.mu.Lock()
	r.ring.Link(s.ring)
	s.mu.Unlock()
	r.mu.Unlock()
	r.dirty.Set(true)
	s.dirty.Set(true)
	return r
}

// Unlink removes n % r.Len() elements from the ring r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed sub-ring. r must not be empty.
func (r *Ring[T]) Unlink(n int) *Ring[T] {
	r.mu.Lock()
	resultRing := r.ring.Unlink(n)
	r.dirty.Set(true)
	r.mu.Unlock()
	resultGRing := New[T](resultRing.Len())
	resultGRing.ring = resultRing
	resultGRing.dirty.Set(true)
	return resultGRing
}

// RLockIteratorNext iterates and locks reading forward
// with given callback function `f` within RWMutex.RLock.
// If `f` returns true, then it continues iterating; or false to stop.
func (r *Ring[T]) RLockIteratorNext(f func(value T) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.ring.Value != nil && !f(r.ring.Value.(internalRingItem[T]).Value) {
		return
	}
	for p := r.ring.Next(); p != r.ring; p = p.Next() {
		if p.Value == nil || !f(p.Value.(internalRingItem[T]).Value) {
			break
		}
	}
}

// RLockIteratorPrev iterates and locks writing backward
// with given callback function `f` within RWMutex.RLock.
// If `f` returns true, then it continues iterating; or false to stop.
func (r *Ring[T]) RLockIteratorPrev(f func(value T) bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.ring.Value != nil && !f(r.ring.Value.(internalRingItem[T]).Value) {
		return
	}
	for p := r.ring.Prev(); p != r.ring; p = p.Prev() {
		if p.Value == nil || !f(p.Value.(internalRingItem[T]).Value) {
			break
		}
	}
}

// SliceNext returns a copy of all item values as slice forward from current position.
func (r *Ring[T]) SliceNext() []T {
	s := make([]T, 0)
	r.mu.RLock()
	if r.ring.Value != nil {
		s = append(s, r.ring.Value.(internalRingItem[T]).Value)
	}
	for p := r.ring.Next(); p != r.ring; p = p.Next() {
		if p.Value == nil {
			break
		}
		s = append(s, p.Value.(internalRingItem[T]).Value)
	}
	r.mu.RUnlock()
	return s
}

// SlicePrev returns a copy of all item values as slice backward from current position.
func (r *Ring[T]) SlicePrev() []T {
	s := make([]T, 0)
	r.mu.RLock()
	if r.ring.Value != nil {
		s = append(s, r.ring.Value.(internalRingItem[T]).Value)
	}
	for p := r.ring.Prev(); p != r.ring; p = p.Prev() {
		if p.Value == nil {
			break
		}
		s = append(s, p.Value.(internalRingItem[T]).Value)
	}
	r.mu.RUnlock()
	return s
}
