//go:build windows

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file checks Windows-specific session file behavior.

package gsession_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_StorageFileSharedReader checks that an open read handle does not block
// writes: the rename path falls back to direct write so SetSession never errors.
func Test_StorageFileSharedReader(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprintf("encrypted=%v", encrypted), func(t *testing.T) {
			gtest.C(t, func(gt *gtest.T) {
				var (
					ctx       = context.Background()
					path      = t.TempDir()
					sessionId = "shared-reader"
					ttl       = time.Minute
					storage   = gsession.NewStorageFile(path, ttl)
					data      = gmap.NewStrAnyMapFrom(g.Map{"value": "first"}, true)
				)
				storage.SetCryptoEnabled(encrypted)
				gt.AssertNil(storage.SetSession(ctx, sessionId, data, ttl))
				reader, err := os.Open(gfile.Join(path, sessionId+".session"))
				gt.AssertNil(err)
				defer func() {
					gt.AssertNil(reader.Close())
				}()
				data.Set("value", "second")
				gt.AssertNil(storage.SetSession(ctx, sessionId, data, ttl))
				loaded, err := storage.GetSession(ctx, sessionId, ttl)
				gt.AssertNil(err)
				gt.AssertNE(loaded, nil)
				gt.Assert(loaded.Get("value"), "second")
			})
		})
	}
}

// Test_Issue4792_WindowsConcurrentReadWrite is a concurrent regression for
// https://github.com/gogf/gf/issues/4792 on Windows. It runs two writers and
// four readers against the same session ID for two seconds and asserts that no
// reader ever observes an empty (nil) session or a read error, and that both
// writes and reads actually happened.
func Test_Issue4792_WindowsConcurrentReadWrite(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx       = context.Background()
			path      = t.TempDir()
			sessionId = "race-session"
			ttl       = time.Hour
			storage   = gsession.NewStorageFile(path, ttl)
			data      = gmap.NewStrAnyMapFrom(g.Map{"userId": 1}, true)
		)
		gt.AssertNil(storage.SetSession(ctx, sessionId, data, ttl))

		var (
			stop      int32
			wg        sync.WaitGroup
			nilReads  int64
			readErrs  int64
			writeErrs int64
			writes    int64
			reads     int64
		)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for atomic.LoadInt32(&stop) == 0 {
					if err := storage.SetSession(ctx, sessionId, data, ttl); err != nil {
						atomic.AddInt64(&writeErrs, 1)
						continue
					}
					atomic.AddInt64(&writes, 1)
				}
			}()
		}
		var rg sync.WaitGroup
		for i := 0; i < 4; i++ {
			rg.Add(1)
			go func() {
				defer rg.Done()
				for atomic.LoadInt32(&stop) == 0 {
					got, err := storage.GetSession(ctx, sessionId, ttl)
					atomic.AddInt64(&reads, 1)
					if err != nil {
						atomic.AddInt64(&readErrs, 1)
						continue
					}
					if got == nil {
						atomic.AddInt64(&nilReads, 1)
					}
				}
			}()
		}
		time.Sleep(2 * time.Second)
		atomic.StoreInt32(&stop, 1)
		wg.Wait()
		rg.Wait()

		gt.AssertEQ(writeErrs, int64(0))
		gt.AssertEQ(readErrs, int64(0))
		gt.AssertEQ(nilReads, int64(0))
		gt.AssertGT(writes, int64(0))
		gt.AssertGT(reads, int64(0))
	})
}

// Test_StorageFileWriteError asserts that SetSession returns an error when the
// session file path is occupied by a directory (OpenWithFlagPerm fails) and does
// not panic.
func Test_StorageFileWriteError(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx       = context.Background()
			path      = t.TempDir()
			sessionId = "err-session"
			ttl       = time.Minute
			storage   = gsession.NewStorageFile(path, ttl)
			data      = gmap.NewStrAnyMapFrom(g.Map{"x": 1}, true)
		)
		// Place a directory at the session file path so that OpenWithFlagPerm fails.
		sessionPath := gfile.Join(path, sessionId+".session")
		gt.AssertNil(os.Mkdir(sessionPath, os.ModePerm))

		err := storage.SetSession(ctx, sessionId, data, ttl)
		gt.AssertNE(err, nil)
	})
}
