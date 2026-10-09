//go:build windows

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file checks that open readers do not block Windows session writes.

package gsession_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_StorageFileSharedReader preserves direct writes with an open read handle.
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
