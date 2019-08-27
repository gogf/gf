// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gkvdb_test

import (
	"strings"
	"testing"
	"time"

	"github.com/dgraph-io/badger/options"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/garray"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/database/gkvdb"

	"github.com/gogf/gf/test/gtest"
)

func init() {

}

func Test_New(t *testing.T) {
	gtest.Case(t, func() {
		name := gconv.String(gtime.Nanosecond())
		path := "/tmp/gkvdb/" + name
		key := []byte("key")
		value := []byte("value")

		db := gkvdb.Instance(name)
		// https://github.com/dgraph-io/badger#memory-usage
		db.Options().ValueLogLoadingMode = options.FileIO
		err := db.SetPath(path)
		gtest.Assert(err, nil)

		err = db.Set(key, value)
		gtest.Assert(err, nil)

		gtest.Assert(db.Get(key), value)
		gtest.Assert(db.Delete(key), nil)
		gtest.Assert(db.Get(key), nil)
	})
}

func Test_Set(t *testing.T) {
	gtest.Case(t, func() {
		name := gconv.String(gtime.Nanosecond())
		path := "/tmp/gkvdb/" + name
		key := []byte("key")
		value := []byte("value")

		db := gkvdb.Instance(name)
		// https://github.com/dgraph-io/badger#memory-usage
		db.Options().ValueLogLoadingMode = options.FileIO
		err := db.SetPath(path)
		gtest.Assert(err, nil)

		err = db.Set(key, value, 1000*time.Millisecond)
		gtest.Assert(err, nil)

		gtest.Assert(db.Get(key), value)
		time.Sleep(1500 * time.Millisecond)
		gtest.Assert(db.Get(key), nil)
	})
}

func Test_Iterate(t *testing.T) {
	gtest.Case(t, func() {
		name := gconv.String(gtime.Nanosecond())
		path := "/tmp/gkvdb/" + name
		db := gkvdb.Instance(name)
		// https://github.com/dgraph-io/badger#memory-usage
		db.Options().ValueLogLoadingMode = options.FileIO
		err := db.SetPath(path)
		gtest.Assert(err, nil)

		strArray := garray.NewSortedStringArray()
		strArrayReverse := garray.NewSortedStringArrayComparator(func(a, b string) int {
			switch strings.Compare(a, b) {
			case 0:
				return 0
			case 1:
				return -1
			case -1:
				return 1
			}
			return 0
		})
		for i := 1; i <= 10; i++ {
			key := []byte("key_" + gconv.String(i))
			strArray.Add(string(key))
			strArrayReverse.Add(string(key))
			db.Set(key, key)
		}

		array := garray.New()
		db.Iterate(nil, func(key, value []byte) bool {
			array.Append(string(key))
			return true
		})
		gtest.Assert(array.Slice(), strArray.Slice())

		array = garray.New()
		db.IterateDesc(nil, func(key, value []byte) bool {
			array.Append(string(key))
			return true
		})
		gtest.Assert(array.Slice(), strArrayReverse.Slice())

		array = garray.New()
		db.Iterate([]byte("key_1"), func(key, value []byte) bool {
			array.Append(key)
			return true
		})
		gtest.Assert(array.Slice(), g.Slice{[]byte("key_1"), []byte("key_10")})

		array = garray.New()
		db.IterateAsc([]byte("key_1"), func(key, value []byte) bool {
			array.Append(key)
			return true
		})
		gtest.Assert(array.Slice(), g.Slice{[]byte("key_1"), []byte("key_10")})

	})
}
