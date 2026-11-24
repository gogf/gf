// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gspath

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Get_EmptyRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(Get("", false), nil)
	})
}

func Test_Search_CacheIndexFileAtRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := filepath.Clean(t.TempDir())
		indexPath := filepath.Join(root, "index.html")
		t.AssertNil(os.WriteFile(indexPath, []byte("index"), 0o644))

		sp := New(root, false)
		sp.cache = gmap.NewStrStrMap(true)
		sp.cache.Set("/", sp.makeCacheValue(root, true))
		sp.cache.Set("/index.html", sp.makeCacheValue(indexPath, false))

		path, isDir := sp.Search("/", "index.html")
		t.Assert(path, indexPath)
		t.Assert(isDir, false)
	})
}

func Test_isPathUnderRoot_MixedPaths(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Absolute and relative paths cannot be compared.
		t.Assert(isPathUnderRoot("relative", "/absolute/file.txt"), false)
	})
}

func Test_isRealPathUnderRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := filepath.Join(t.TempDir(), "root")
		t.AssertNil(os.Mkdir(root, 0o755))
		realRoot, err := filepath.EvalSymlinks(root)
		t.AssertNil(err)

		// Missing paths cannot be resolved.
		t.Assert(isRealPathUnderRoot(realRoot, filepath.Join(root, "missing.txt")), false)
	})
}

func Test_addToCache_SkipsEscapedItems(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		base := t.TempDir()
		root := filepath.Join(base, "root")
		outside := filepath.Join(base, "outside")
		t.AssertNil(os.Mkdir(root, 0o755))
		t.AssertNil(os.Mkdir(outside, 0o755))
		t.AssertNil(os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o644))

		sp := New(root, false)
		sp.cache = gmap.NewStrStrMap(true)

		// An unresolvable root aborts the cache update.
		sp.addToCache(filepath.Join(root, "missing.txt"), filepath.Join(base, "missing-root"))
		t.Assert(sp.cache.Size(), 0)

		// Files resolving outside the root are skipped.
		sp.addToCache(filepath.Join(outside, "secret.txt"), root)
		t.Assert(sp.cache.Size(), 0)
	})
}
