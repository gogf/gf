// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT License was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file defines the Cmd type used as a future-result container for
// commands queued in a Pipeline or TxPipeline.

package gredis

import (
	"github.com/gogf/gf/v2/container/gvar"
)

// Cmd holds the future result of a command queued in a Pipeline or TxPipeline.
// The result is populated after Pipeliner.Exec() or Tx.Exec().
// Before Exec, Cmd.Val() returns nil and Cmd.Result() returns (nil, nil).
type Cmd struct {
	// val holds the populated result after Exec. It is nil until Exec populates it.
	val *gvar.Var

	// err holds the populated error after Exec. It is nil before Exec or if no error occurred.
	err error
}

// Result returns the populated value and error after Exec.
// Before Exec is called, it returns (nil, nil).
func (c *Cmd) Result() (*gvar.Var, error) {
	return c.val, c.err
}

// Val returns the populated *gvar.Var value.
// It returns nil before Exec has been called or if the command returned no value.
func (c *Cmd) Val() *gvar.Var {
	return c.val
}

// Err returns the populated error after Exec.
// It returns nil before Exec has been called or if no error occurred.
func (c *Cmd) Err() error {
	return c.err
}

// SetVal sets the result value of the Cmd. Used by driver implementations
// to populate the result after Exec.
func (c *Cmd) SetVal(val *gvar.Var) {
	c.val = val
}

// SetErr sets the result error of the Cmd. Used by driver implementations
// to populate the error after Exec.
func (c *Cmd) SetErr(err error) {
	c.err = err
}
