// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gins provides instances and core components management.
package gins

import (
	"github.com/gogf/gf/container/gmap"
)

var (
	// instances is the instance map for common used components.
	instances = gmap.NewStrAnyMap(true)
)

// Get returns the instance by given name.
func Get(name string) interface{} {
	return instances.Get(name)
}

// Set sets a instance object to the instance manager with given name.
func Set(name string, instance interface{}) {
	instances.Set(name, instance)
}

// GetOrSet returns the instance by name,
// or set instance to the instance manager if it does not exist and returns this instance.
func GetOrSet(name string, instance interface{}) interface{} {
	return instances.GetOrSet(name, instance)
}

// GetOrSetFunc returns the instance by name,
// or sets instance with returned value of callback function <f> if it does not exist
// and then returns this instance.
func GetOrSetFunc(name string, f func() interface{}) interface{} {
	return instances.GetOrSetFunc(name, f)
}

// GetOrSetFuncLock returns the instance by name,
// or sets instance with returned value of callback function <f> if it does not exist
// and then returns this instance.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func GetOrSetFuncLock(name string, f func() interface{}) interface{} {
	return instances.GetOrSetFuncLock(name, f)
}

// SetIfNotExist sets <instance> to the map if the <name> does not exist, then returns true.
// It returns false if <name> exists, and <instance> would be ignored.
func SetIfNotExist(name string, instance interface{}) bool {
	return instances.SetIfNotExist(name, instance)
}
