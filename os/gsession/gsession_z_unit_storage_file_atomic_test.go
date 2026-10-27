//go:build unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies publication of complete file session snapshots.

package gsession_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_Issue4792 checks same-id writers and readers across storage instances.
func Test_Issue4792(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprintf("encrypted=%v", encrypted), func(t *testing.T) {
			gtest.C(t, func(gt *gtest.T) {
				var (
					ctx       = context.Background()
					path      = t.TempDir()
					sessionId = "concurrent-session"
					ttl       = time.Minute
					storages  = []*gsession.StorageFile{
						gsession.NewStorageFile(path, ttl),
						gsession.NewStorageFile(path, ttl),
					}
					snapshots = make(map[string]string)
					start     = make(chan struct{})
					results   = make(chan error, 8)
					wg        sync.WaitGroup
				)
				for _, storage := range storages {
					storage.SetCryptoEnabled(encrypted)
				}
				for i := 0; i < 4; i++ {
					writer := fmt.Sprintf("writer-%d", i)
					snapshots[writer] = strings.Repeat(writer, 2048*(i+1))
				}
				gt.AssertNil(storages[0].SetSession(ctx, sessionId, gmap.NewStrAnyMapFrom(g.Map{
					"writer": "writer-0", "payload": snapshots["writer-0"],
				}, true), ttl))

				for i := 0; i < 4; i++ {
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						<-start
						writer := fmt.Sprintf("writer-%d", index)
						data := gmap.NewStrAnyMapFrom(g.Map{
							"writer": writer, "payload": snapshots[writer],
						}, true)
						for j := 0; j < 100; j++ {
							if err := storages[index%len(storages)].SetSession(ctx, sessionId, data, ttl); err != nil {
								results <- fmt.Errorf("writer %d: %w", index, err)
								return
							}
						}
					}(i)
					wg.Add(1)
					go func(index int) {
						defer wg.Done()
						<-start
						for j := 0; j < 400; j++ {
							data, err := storages[index%len(storages)].GetSession(ctx, sessionId, ttl)
							if err != nil {
								results <- fmt.Errorf("reader %d: %w", index, err)
								return
							}
							if data == nil {
								results <- fmt.Errorf("reader %d observed an empty session", index)
								return
							}
							writer, ok := data.Get("writer").(string)
							payload, exists := snapshots[writer]
							if !ok || !exists || data.Get("payload") != payload {
								results <- fmt.Errorf("reader %d observed an incomplete snapshot", index)
								return
							}
						}
					}(i)
				}
				close(start)
				wg.Wait()
				close(results)
				for err := range results {
					gt.AssertNil(err)
				}
				files, err := os.ReadDir(path)
				gt.AssertNil(err)
				gt.Assert(len(files), 1)
				gt.Assert(files[0].Name(), sessionId+".session")
			})
		})
	}
}

// Test_StorageFileAtomicWriteLongId preserves the maximum session filename length.
func Test_StorageFileAtomicWriteLongId(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			path      = t.TempDir()
			ttl       = time.Minute
			storage   = gsession.NewStorageFile(path, ttl)
			sessionId = strings.Repeat("a", 247)
			data      = gmap.NewStrAnyMapFrom(g.Map{"value": "long-id"}, true)
		)
		gt.AssertNil(storage.SetSession(context.Background(), sessionId, data, ttl))
		loaded, err := storage.GetSession(context.Background(), sessionId, ttl)
		gt.AssertNil(err)
		gt.Assert(loaded.Get("value"), "long-id")
		entries, err := os.ReadDir(path)
		gt.AssertNil(err)
		gt.Assert(len(entries), 1)
		gt.Assert(entries[0].Name(), sessionId+".session")
	})
}

// Test_StorageFileAtomicWriteLifecycle checks the file format, permissions and expiry.
func Test_StorageFileAtomicWriteLifecycle(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		t.Run(fmt.Sprintf("encrypted=%v", encrypted), func(t *testing.T) {
			gtest.C(t, func(gt *gtest.T) {
				var (
					ctx       = context.Background()
					path      = t.TempDir()
					sessionId = "lifecycle-session"
					filePath  = gfile.Join(path, sessionId+".session")
					ttl       = time.Minute
					storage   = gsession.NewStorageFile(path, ttl)
					data      = gmap.NewStrAnyMapFrom(g.Map{"value": "first"}, true)
					before    = gtime.TimestampMilli()
				)
				storage.SetCryptoEnabled(encrypted)
				gt.AssertNil(storage.SetSession(ctx, sessionId, data, ttl))
				content, err := os.ReadFile(filePath)
				gt.AssertNil(err)
				gt.Assert(len(content) > 8, true)
				timestamp := gbinary.DecodeToInt64(content[:8])
				gt.Assert(timestamp >= before && timestamp <= gtime.TimestampMilli(), true)
				info, statErr := os.Stat(filePath)
				gt.AssertNil(statErr)
				gt.Assert(info.Mode().Perm(), os.FileMode(0o600))
				gt.AssertNil(os.Chmod(filePath, 0o640))
				data.Set("value", "second")
				gt.AssertNil(storage.SetSession(ctx, sessionId, data, ttl))
				loaded, err := storage.GetSession(ctx, sessionId, ttl)
				gt.AssertNil(err)
				gt.Assert(loaded.Get("value"), "second")
				info, statErr = os.Stat(filePath)
				gt.AssertNil(statErr)
				gt.Assert(info.Mode().Perm(), os.FileMode(0o640))
				time.Sleep(10 * time.Millisecond)
				expired, err := storage.GetSession(ctx, sessionId, time.Millisecond)
				gt.AssertNil(err)
				gt.AssertNil(expired)
				gt.AssertNil(storage.RemoveAll(ctx, sessionId))
				removed, err := storage.GetSession(ctx, sessionId, ttl)
				gt.AssertNil(err)
				gt.AssertNil(removed)
				entries, err := os.ReadDir(path)
				gt.AssertNil(err)
				gt.Assert(len(entries), 0)
			})
		})
	}
}

// Test_StorageFileAtomicWriteFailure checks failed publications leave no temporary files.
func Test_StorageFileAtomicWriteFailure(t *testing.T) {
	t.Run("stat", func(t *testing.T) {
		gtest.C(t, func(gt *gtest.T) {
			var (
				path     = t.TempDir()
				ttl      = time.Minute
				storage  = gsession.NewStorageFile(path, ttl)
				data     = gmap.NewStrAnyMapFrom(g.Map{"value": "test"}, true)
				filePath = gfile.Join(path, "loop.session")
			)
			gt.AssertNil(os.Symlink(filePath, filePath))
			gt.AssertNE(storage.SetSession(context.Background(), "loop", data, ttl), nil)
			entries, err := os.ReadDir(path)
			gt.AssertNil(err)
			gt.Assert(len(entries), 1)
			gt.Assert(entries[0].Name(), "loop.session")
			gt.Assert(entries[0].Type()&os.ModeSymlink != 0, true)
		})
	})
	t.Run("rename", func(t *testing.T) {
		gtest.C(t, func(gt *gtest.T) {
			var (
				path    = t.TempDir()
				ttl     = time.Minute
				storage = gsession.NewStorageFile(path, ttl)
				data    = gmap.NewStrAnyMapFrom(g.Map{"value": "test"}, true)
			)
			gt.AssertNil(os.Mkdir(gfile.Join(path, "blocked.session"), 0o700))
			gt.AssertNE(storage.SetSession(context.Background(), "blocked", data, ttl), nil)
			entries, err := os.ReadDir(path)
			gt.AssertNil(err)
			gt.Assert(len(entries), 1)
			gt.Assert(entries[0].Name(), "blocked.session")
			gt.Assert(entries[0].IsDir(), true)
		})
	})
	t.Run("create", func(t *testing.T) {
		gtest.C(t, func(gt *gtest.T) {
			var (
				path    = t.TempDir()
				ttl     = time.Minute
				storage = gsession.NewStorageFile(path, ttl)
				data    = gmap.NewStrAnyMapFrom(g.Map{"value": "test"}, true)
			)
			gt.AssertNil(os.Remove(path))
			gt.AssertNE(storage.SetSession(context.Background(), "missing", data, ttl), nil)
			_, err := os.Stat(path)
			gt.Assert(os.IsNotExist(err), true)
		})
	})
}
