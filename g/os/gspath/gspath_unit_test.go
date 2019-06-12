package gspath_test

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gspath"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestSPath_Api(t *testing.T) {
	gtest.Case(t, func() {
		pwd := gfile.Pwd()
		root := pwd + gfile.Separator
		gfile.Create(root + "gf_tmp" + gfile.Separator + "gf.txt")
		defer gfile.Remove(root + "gf_tmp")
		fp, isDir := gspath.Search(root, "gf_tmp")
		gtest.Assert(fp, root+"gf_tmp")
		gtest.Assert(isDir, true)
		fp, isDir = gspath.Search(root, "gf_tmp", "gf.txt")
		gtest.Assert(fp, root+"gf_tmp"+gfile.Separator+"gf.txt")
		gtest.Assert(isDir, false)

		fp, isDir = gspath.SearchWithCache(root, "gf_tmp")
		gtest.Assert(fp, root+"gf_tmp")
		gtest.Assert(isDir, true)
		fp, isDir = gspath.SearchWithCache(root, "gf_tmp", "gf.txt")
		gtest.Assert(fp, root+"gf_tmp"+gfile.Separator+"gf.txt")
		gtest.Assert(isDir, false)
	})
}

func TestSPath_Basic(t *testing.T) {
	gtest.Case(t, func() {
		pwd := gfile.Pwd()
		root := pwd + gfile.Separator
		gfile.Create(root + "gf_tmp" + gfile.Separator + "gf.txt")
		defer gfile.Remove(root + "gf_tmp")
		gsp := gspath.New(root, false)
		realPath, err := gsp.Add(root + "gf_tmp")
		gtest.Assert(err, nil)
		gtest.Assert(realPath, root+"gf_tmp")
		realPath, err = gsp.Add("gf_tmp1")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		realPath, err = gsp.Add(root + "gf_tmp" + gfile.Separator + "gf.txt")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		gsp.Remove("gf_tmp1")
		gtest.Assert(gsp.Size(), 2)
		gtest.Assert(len(gsp.Paths()), 2)
		gtest.Assert(len(gsp.AllPaths()), 0)
		realPath, err = gsp.Set(root + "gf_tmp1")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		realPath, err = gsp.Set(root + "gf_tmp" + gfile.Separator + "gf.txt")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		gsp.Set(root)
		fp, isDir := gsp.Search("gf_tmp")
		gtest.Assert(fp, root+"gf_tmp")
		gtest.Assert(isDir, true)
		fp, isDir = gsp.Search("gf_tmp", "gf.txt")
		gtest.Assert(fp, root+"gf_tmp"+gfile.Separator+"gf.txt")
		gtest.Assert(isDir, false)
		fp, isDir = gsp.Search("/", "gf.txt")
		gtest.Assert(fp, root+gfile.Separator)
		gtest.Assert(isDir, true)

		gsp = gspath.New(root, true)
		realPath, err = gsp.Add(root + "gf_tmp")
		gtest.Assert(err, nil)
		gtest.Assert(realPath, root+"gf_tmp")

		gfile.Mkdir(root + "gf_tmp1")
		gfile.Rename(root+"gf_tmp1", root+"gf_tmp2")
		gfile.Rename(root+"gf_tmp2", root+"gf_tmp1")
		defer gfile.Remove(root + "gf_tmp1")
		realPath, err = gsp.Add("gf_tmp1")
		gtest.Assert(err != nil, false)
		gtest.Assert(realPath, root+"gf_tmp1")
		realPath, err = gsp.Add("gf_tmp3")
		gtest.Assert(err != nil, true)
		gtest.Assert(realPath, "")
		gsp.Remove(root + "gf_tmp")
		gsp.Remove(root + "gf_tmp1")
		gsp.Remove(root + "gf_tmp3")
		gtest.Assert(gsp.Size(), 3)
		gtest.Assert(len(gsp.Paths()), 3)
		gsp.AllPaths()
		gsp.Set(root)
		fp, isDir = gsp.Search("gf_tmp")
		gtest.Assert(fp, root+"gf_tmp")
		gtest.Assert(isDir, true)
		fp, isDir = gsp.Search("gf_tmp", "gf.txt")
		gtest.Assert(fp, root+"gf_tmp"+gfile.Separator+"gf.txt")
		gtest.Assert(isDir, false)
		fp, isDir = gsp.Search("/", "gf.txt")
		gtest.Assert(fp, pwd)
		gtest.Assert(isDir, true)
	})
}
