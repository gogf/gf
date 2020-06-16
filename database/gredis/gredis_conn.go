// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"encoding/json"
	"github.com/gogf/gf/container/gvar"
	"reflect"
)

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *Conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	for k, v := range args {
		reflectValue = reflect.ValueOf(v)
		reflectKind = reflectValue.Kind()
		if reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		case
			reflect.Struct,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			// Ignore slice type of: []byte.
			if _, ok := v.([]byte); !ok {
				if args[k], err = json.Marshal(v); err != nil {
					return nil, err
				}
			}
		}
	}
	return c.Conn.Do(commandName, args...)
}

// DoVar retrieves and returns the result from command as gvar.Var.
func (c *Conn) DoVar(command string, args ...interface{}) (gvar.Var, error) {
	v, err := c.Do(command, args...)
	return gvar.New(v), err
}

// ReceiveVar receives a single reply as gvar.Var from the Redis server.
func (c *Conn) ReceiveVar() (gvar.Var, error) {
	v, err := c.Receive()
	return gvar.New(v), err
}
