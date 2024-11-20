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
		tests           = []struct {
			tag    string
			path   string
			fspath fs.FS
			opt    gres.Option
		}{
			{
				tag:    "none@true",
				path:   srcAbsPath,
				fspath: os.DirFS(srcAbsPath),
				opt: gres.Option{
					Prefix:   "",
					KeepPath: true,
				},
			},
			{
				tag:    "none@false",
				path:   srcAbsPath,
				fspath: os.DirFS(srcAbsPath),
				opt: gres.Option{
					Prefix:   "",
					KeepPath: false,
				},
			},
			{
				tag:    "files@true",
				path:   srcAbsPath,
				fspath: os.DirFS(srcAbsPath),
				opt: gres.Option{
					Prefix:   "files",
					KeepPath: true,
				},
			},
			{
				tag:    "files@false",
				path:   srcAbsPath,
				fspath: os.DirFS(srcAbsPath),

				opt: gres.Option{
					Prefix:   "files",
					KeepPath: false,
				},
			},
			{
				tag:    "./testdata/files@true",
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),
				opt: gres.Option{
					Prefix:   "./testdata/files",
					KeepPath: true,
				},
			},
			{
				tag:    "./testdata/files@false",
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),

				opt: gres.Option{
					Prefix:   "./testdata/files",
					KeepPath: false,
				},
			},
			{
				tag:    "./testdata/files@true@/testdata/files",
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),
				opt: gres.Option{
					Prefix:   "/testdata/files",
					KeepPath: true,
				},
			},
			{
				tag:    "./testdata/files@false@/testdata/files",
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),

				opt: gres.Option{
					Prefix:   "/testdata/files",
					KeepPath: false,
				},
			},
			{
				tag:    "./testdata/files@true@/",
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),
				opt: gres.Option{
					Prefix:   "/",
					KeepPath: true,
				},
			},
			{
				tag:    "./testdata/files@false@/",
				path:   srcRelativePath,
				fspath: os.DirFS(srcRelativePath),

				opt: gres.Option{
					Prefix:   "/",
					KeepPath: false,
				},
			},
		}
	)
	for _, testinfo := range tests {
		t.Log(testinfo.tag, ":")
		var (
			t1, err1         = gres.PackWithOption(testinfo.path, testinfo.opt)
			t2, err2         = gres.PackFsWithOption(testinfo.fspath, testinfo.path, testinfo.opt)
			r1, r2           = gres.New(), gres.New()
			err3             = r1.Add(string(t1))
			err4             = r2.Add(string(t2))
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
		// _, _ = r1files, r2files
	}
}
