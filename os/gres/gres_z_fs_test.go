package gres_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_PackFS(t *testing.T) {
	var (
		srcAbsPath      = gtest.DataPath("files")
		srcRelativePath = "./testdata/files"
		tests           = map[string]struct {
			path   string
			fspath fs.FS
			opt    gres.Option
		}{
			"files@true": {
				path:   srcAbsPath,
				fspath: os.DirFS(srcAbsPath),
				opt: gres.Option{
					Prefix:   "files",
					KeepPath: true,
				},
			},
			"files@false": {
				path:   srcAbsPath,
				fspath: os.DirFS(srcAbsPath),

				opt: gres.Option{
					Prefix:   "files",
					KeepPath: false,
				},
			},
			"./testdata/files@true": {
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),
				opt: gres.Option{
					Prefix:   "./testdata/files",
					KeepPath: true,
				},
			},
			"./testdata/files@false": {
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),

				opt: gres.Option{
					Prefix:   "./testdata/files",
					KeepPath: false,
				},
			},
			"./testdata/files@true1": {
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),
				opt: gres.Option{
					Prefix:   "/testdata/files",
					KeepPath: true,
				},
			},
			"./testdata/files@false1": {
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),

				opt: gres.Option{
					Prefix:   "/testdata/files",
					KeepPath: false,
				},
			},
			"./testdata/files@true2": {
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),
				opt: gres.Option{
					Prefix:   "/",
					KeepPath: true,
				},
			},
			"./testdata/files@false2": {
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),

				opt: gres.Option{
					Prefix:   "/",
					KeepPath: false,
				},
			},
		}
	)
	for key, testinfo := range tests {
		t.Log(key, ":")
		var (
			t1, err1         = gres.PackWithOption(testinfo.path, testinfo.opt)
			t2, err2         = gres.PackFsWithOption(testinfo.fspath, "files", testinfo.opt)
			r1, r2           = gres.New(), gres.New()
			err3             = r1.Add(string(t1))
			err4             = r2.Add(string(t1))
			r1files, r2files = r1.ScanDir(".", "*", true), r2.ScanDir(".", "*", true)
		)
		gtest.AssertNil(err1)
		gtest.AssertNil(err2)
		gtest.AssertNil(err3)
		gtest.AssertNil(err4)
		gtest.AssertEQ(t1, t2)
		gtest.Assert(r1files, r2files)
		// t.Log("r1:")
		// r1.Dump()
		// t.Log("r2:")
		// r2.Dump()
	}
}
