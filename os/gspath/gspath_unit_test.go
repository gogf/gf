// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gspath_test

import (
	"testing"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gspath"
	"github.com/gogf/gf/test/gtest"
)

func TestSPath_Api(t *testing.T) {
	gtest.Case(t, func() {
		pwd := gfile.Pwd()
		root := pwd
		gfile.Create(gfile.Join(root, "gf_tmp", "gf.txt"))
		defer gfile.Remove(gfile.Join(root, "gf_tmp"))
		fp, isDir := gspath.Search(root, "gf_tmp")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp"))
		gtest.Assert(isDir, true)
		fp, isDir = gspath.Search(root, "gf_tmp", "gf.txt")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp", "gf.txt"))
		gtest.Assert(isDir, false)

		fp, isDir = gspath.SearchWithCache(root, "gf_tmp")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp"))
		gtest.Assert(isDir, true)
		fp, isDir = gspath.SearchWithCache(root, "gf_tmp", "gf.txt")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp", "gf.txt"))
		gtest.Assert(isDir, false)
	})
}

func TestSPath_Basic(t *testing.T) {
	gtest.Case(t, func() {
		pwd := gfile.Pwd()
		root := pwd

		gfile.Create(gfile.Join(root, "gf_tmp", "gf.txt"))
		defer gfile.Remove(gfile.Join(root, "gf_tmp"))
		gsp := gspath.New(root, false)
		realPath, err := gsp.Add(gfile.Join(root, "gf_tmp"))
		gtest.Assert(err, nil)
		gtest.Assert(realPath, gfile.Join(root, "gf_tmp"))
		realPath, err = gsp.Add("gf_tmp1")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		realPath, err = gsp.Add(gfile.Join(root, "gf_tmp", "gf.txt"))
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		gsp.Remove("gf_tmp1")
		gtest.Assert(gsp.Size(), 2)
		gtest.Assert(len(gsp.Paths()), 2)
		gtest.Assert(len(gsp.AllPaths()), 0)
		realPath, err = gsp.Set(gfile.Join(root, "gf_tmp1"))
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		realPath, err = gsp.Set(gfile.Join(root, "gf_tmp", "gf.txt"))
		gtest.AssertNE(err, nil)
		gtest.Assert(realPath, "")

		realPath, err = gsp.Set(root)
		gtest.Assert(err, nil)
		gtest.Assert(realPath, root)

		fp, isDir := gsp.Search("gf_tmp")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp"))
		gtest.Assert(isDir, true)
		fp, isDir = gsp.Search("gf_tmp", "gf.txt")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp", "gf.txt"))
		gtest.Assert(isDir, false)
		fp, isDir = gsp.Search("/", "gf.txt")
		gtest.Assert(fp, root)
		gtest.Assert(isDir, true)

		gsp = gspath.New(root, true)
		realPath, err = gsp.Add(gfile.Join(root, "gf_tmp"))
		gtest.Assert(err, nil)
		gtest.Assert(realPath, gfile.Join(root, "gf_tmp"))

		gfile.Mkdir(gfile.Join(root, "gf_tmp1"))
		gfile.Rename(gfile.Join(root, "gf_tmp1"), gfile.Join(root, "gf_tmp2"))
		gfile.Rename(gfile.Join(root, "gf_tmp2"), gfile.Join(root, "gf_tmp1"))
		defer gfile.Remove(gfile.Join(root, "gf_tmp1"))
		realPath, err = gsp.Add("gf_tmp1")
		gtest.Assert(err != nil, false)
		gtest.Assert(realPath, gfile.Join(root, "gf_tmp1"))
		realPath, err = gsp.Add("gf_tmp3")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		gsp.Remove(gfile.Join(root, "gf_tmp"))
		gsp.Remove(gfile.Join(root, "gf_tmp1"))
		gsp.Remove(gfile.Join(root, "gf_tmp3"))
		gtest.Assert(gsp.Size(), 3)
		gtest.Assert(len(gsp.Paths()), 3)
		gsp.AllPaths()
		gsp.Set(root)
		fp, isDir = gsp.Search("gf_tmp")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp"))
		gtest.Assert(isDir, true)
		fp, isDir = gsp.Search("gf_tmp", "gf.txt")
		gtest.Assert(fp, gfile.Join(root, "gf_tmp", "gf.txt"))
		gtest.Assert(isDir, false)
		fp, isDir = gsp.Search("/", "gf.txt")
		gtest.Assert(fp, pwd)
		gtest.Assert(isDir, true)
	})
}
