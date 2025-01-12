// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genservice"
)

func Test_Gen_Service_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path      = gfile.Temp(guid.S())
			dstFolder = path + filepath.FromSlash("/service")
			srvFolder = gtest.DataPath("genservice", "logic")
			in        = genservice.CGenServiceInput{
				SrcFolder:       srvFolder,
				DstFolder:       dstFolder,
				DstFileNameCase: "Snake",
				WatchFile:       "",
				StPattern:       "",
				Packages:        nil,
				ImportPrefix:    "",
				Clear:           false,
			}
		)
		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genservice.CGenService{}.Service(ctx, in)
		t.AssertNil(err)

		// logic file
		var (
			genSrv       = srvFolder + filepath.FromSlash("/logic.go")
			genSrvExpect = srvFolder + filepath.FromSlash("/logic_expect.go")
		)
		defer gfile.Remove(genSrv)
		t.Assert(gfile.GetContents(genSrv), gfile.GetContents(genSrvExpect))

		// files
		files, err := gfile.ScanDir(dstFolder, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			dstFolder + filepath.FromSlash("/article.go"),
			dstFolder + filepath.FromSlash("/base.go"),
			dstFolder + filepath.FromSlash("/delivery.go"),
			dstFolder + filepath.FromSlash("/user.go"),
		})

		// contents
		testPath := gtest.DataPath("genservice", "service")
		expectFiles := []string{
			testPath + filepath.FromSlash("/article.go"),
			testPath + filepath.FromSlash("/base.go"),
			testPath + filepath.FromSlash("/delivery.go"),
			testPath + filepath.FromSlash("/user.go"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

// https://github.com/gogf/gf/issues/3328
func Test_Issue3328(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path        = gfile.Temp(guid.S())
			dstFolder   = path + filepath.FromSlash("/service")
			srvFolder   = gtest.DataPath("issue", "3328", "logic")
			logicGoPath = srvFolder + filepath.FromSlash("/logic.go")
			in          = genservice.CGenServiceInput{
				SrcFolder:       srvFolder,
				DstFolder:       dstFolder,
				DstFileNameCase: "Snake",
				WatchFile:       "",
				StPattern:       "",
				Packages:        nil,
				ImportPrefix:    "",
				Clear:           false,
			}
		)
		gfile.Remove(logicGoPath)
		defer gfile.Remove(logicGoPath)

		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genservice.CGenService{}.Service(ctx, in)
		t.AssertNil(err)

		files, err := gfile.ScanDir(srvFolder, "*", true)
		for _, file := range files {
			if file == logicGoPath {
				if gfile.IsDir(logicGoPath) {
					t.Fatalf("%s should not is folder", logicGoPath)
				}
			}
		}
	})
}

// https://github.com/gogf/gf/issues/3835
func Test_Issue3835(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path      = gfile.Temp(guid.S())
			dstFolder = path + filepath.FromSlash("/service")
			srvFolder = gtest.DataPath("issue", "3835", "logic")
			in        = genservice.CGenServiceInput{
				SrcFolder:       srvFolder,
				DstFolder:       dstFolder,
				DstFileNameCase: "Snake",
				WatchFile:       "",
				StPattern:       "",
				Packages:        nil,
				ImportPrefix:    "",
				Clear:           false,
			}
		)
		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genservice.CGenService{}.Service(ctx, in)
		t.AssertNil(err)

		// contents
		var (
			genFile    = dstFolder + filepath.FromSlash("/issue_3835.go")
			expectFile = gtest.DataPath("issue", "3835", "service", "issue_3835.go")
		)
		t.Assert(gfile.GetContents(genFile), gfile.GetContents(expectFile))
	})
}
