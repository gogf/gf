// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gkvdb provides a lightweight, embeddable and persistent key-value database.
package gkvdb

import (
	"github.com/dgraph-io/badger"
	"github.com/gogf/gf/container/gmap"
)

type Options = badger.Options

const (
	// Default instance name.
	DEFAULT_NAME = "default"
)

var (
	// Instance map
	instances = gmap.NewStrAnyMap(true)
)

// DefaultOptions returns the default options for gkvdb.
func DefaultOptions(path string) Options {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	options.SyncWrites = false
	return options
}

// Instance returns an instance of DB with specified name.
// The <name> param is unnecessary, if <name> is not passed,
// it returns a instance with default name.
func Instance(name ...string) *DB {
	instanceName := DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	v := instances.GetOrSetFuncLock(instanceName, func() interface{} {
		return New()
	})
	if v != nil {
		return v.(*DB)
	}
	return nil
}
