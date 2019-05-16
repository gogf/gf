// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import "github.com/gogf/gf/g/container/gvar"

// DoVar returns value from Do as gvar.Var.
func (c *Conn) DoVar(command string, args ...interface{}) (*gvar.Var, error) {
	v, err := c.Do(command, args...)
	return gvar.New(v, true), err
}

// ReceiveVar receives a single reply as gvar.Var from the Redis server.
func (c *Conn) ReceiveVar() (*gvar.Var, error) {
	v, err := c.Receive()
	return gvar.New(v, true), err
}