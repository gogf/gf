// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_KV(t *testing.T) {
	config := `
[kvdb]
    default = "path=/tmp/gkvdb&sync=false"
    cache   = "path=/tmp/gkvdb-cache&sync=true"
`
	path := "config.toml"
	err := gfile.PutContents(path, config)
	gtest.Assert(err, nil)
	defer gfile.Remove(path)
	defer gins.Config().Clear()

	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	gtest.Case(t, func() {
		kvDefault := gins.KV()
		kvCache := gins.KV("cache")
		key := []byte("key")
		value := []byte("value")
		err := kvDefault.Set(key, value)
		gtest.Assert(err, nil)

		gtest.Assert(kvDefault.Get(key), value)
		gtest.Assert(kvCache.Get(key), nil)
	})
}
