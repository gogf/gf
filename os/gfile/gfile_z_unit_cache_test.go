// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_GetContentsWithCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var f *os.File
		var err error
		fileName := "test"
		strTest := "123"

		if !gfile.Exists(fileName) {
			f, err = os.CreateTemp("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if gfile.Exists(f.Name()) {
			f, err = gfile.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}

			err = gfile.PutContents(f.Name(), strTest)
			if err != nil {
				t.Error("write error", err)
			}

			cache := gfile.GetContentsWithCache(f.Name(), 1)
			t.Assert(cache, strTest)
		}
	})

	gtest.C(t, func(t *gtest.T) {

		var f *os.File
		var err error
		fileName := "test2"
		strTest := "123"

		if !gfile.Exists(fileName) {
			f, err = os.CreateTemp("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if gfile.Exists(f.Name()) {
			cache := gfile.GetContentsWithCache(f.Name())

			f, err = gfile.OpenFile(f.Name(), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				t.Error("file open fail", err)
			}

			err = gfile.PutContents(f.Name(), strTest)
			if err != nil {
				t.Error("write error", err)
			}

			t.Assert(cache, "")

			time.Sleep(100 * time.Millisecond)
		}
	})
}

func Test_GetBytesWithCache(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var f *os.File
		var err error
		fileName := "test_bytes"
		byteContent := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello"

		if !gfile.Exists(fileName) {
			f, err = os.CreateTemp("", fileName)
			if err != nil {
				t.Error("create file fail")
			}
		}

		defer f.Close()
		defer os.Remove(f.Name())

		if gfile.Exists(f.Name()) {
			err = gfile.PutBytes(f.Name(), byteContent)
			if err != nil {
				t.Error("write error", err)
			}

			// Test GetBytesWithCache with custom duration
			cache := gfile.GetBytesWithCache(f.Name(), time.Second*1)
			t.Assert(cache, byteContent)

			// Test cache hit - should return same content
			cache2 := gfile.GetBytesWithCache(f.Name(), time.Second*1)
			t.Assert(cache2, byteContent)
		}
	})

	// Test with non-existent file
	gtest.C(t, func(t *gtest.T) {
		cache := gfile.GetBytesWithCache("/nonexistent_file_12345.txt")
		t.Assert(cache, nil)
	})

	// Test with empty file
	gtest.C(t, func(t *gtest.T) {
		var f *os.File
		var err error
		fileName := "test_bytes_empty"

		f, err = os.CreateTemp("", fileName)
		if err != nil {
			t.Error("create file fail")
		}

		defer f.Close()
		defer os.Remove(f.Name())

		// Read empty file
		cache := gfile.GetBytesWithCache(f.Name(), time.Second*1)
		t.Assert(len(cache), 0)
	})
}
