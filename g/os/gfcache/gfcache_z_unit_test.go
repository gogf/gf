// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gfcache_test

import (
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/g/os/gfcache"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/test/gtest"
)

func TestGetContents(t *testing.T) {
	gtest.Case(t, func() {

		var f *os.File
		var err error
		fileName := "test.txt"
		strTest := "123"

		if !gfile.Exists(fileName) {
			f, err = gfile.Create(fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		cache := gfcache.GetContents(fileName, 2)

		if gfile.Exists(fileName) {
			f, err = gfile.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}
		}

		defer f.Close()

		_, err = f.Write([]byte(strTest))
		if err != nil {
			t.Error("write error", err)
		}

		cache = gfcache.GetContents(fileName)
		gtest.Assert(cache, "")

		time.Sleep(time.Duration(4) * time.Second)

		if gfile.Exists(fileName) {
			err = gfile.Remove(fileName)
			if err != nil {
				t.Error("file remove fail", err)
			}
		}
	})

}
