// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gchan provides graceful channel for safe operations.
//
// It's safe to call Chan.Push/Close functions repeatedly.
package gchan

import (
    "errors"
    "github.com/gogf/gf/g/container/gtype"
)

type Chan struct {
    channel chan interface{}
    closed  *gtype.Bool
}

// New creates a graceful channel with given limit.
func New(limit int) *Chan {
    return &Chan {
	    channel : make(chan interface{}, limit),
        closed  : gtype.NewBool(),
    }
}

// Push pushes <value> to channel.
// It is safe to be called repeatedly.
func (c *Chan) Push(value interface{}) error {
    if c.closed.Val() {
        return errors.New("closed")
    }
    c.channel <- value
    return nil
}

// Pop pops value from channel.
// If there's no value in channel, it would block to wait.
func (c *Chan) Pop() interface{} {
    return <- c.channel
}

// Close closes the channel.
// It is safe to be called repeatedly.
func (c *Chan) Close() {
    if !c.closed.Set(true) {
        close(c.channel)
    }
}

// See Len.
func (c *Chan) Size() int {
    return c.Len()
}

// Len returns the length of the channel.
func (c *Chan) Len() int {
	return len(c.channel)
}