// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genctrl"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_Gen_Ctrl_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path      = gfile.Temp(guid.S())
			apiFolder = gtest.DataPath("genctrl", "api")
			in        = genctrl.CGenCtrlInput{
				SrcFolder:     apiFolder,
				DstFolder:     path,
				WatchFile:     "",
				SdkPath:       "",
				SdkStdVersion: false,
				SdkNoV1:       false,
				Clear:         false,
				Merge:         false,
			}
		)
		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		// apiInterface file
		var (
			genApi       = apiFolder + filepath.FromSlash("/article/article.go")
			genApiExpect = apiFolder + filepath.FromSlash("/article/article_expect.go")
		)
		defer gfile.Remove(genApi)
		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))

		// files
		files, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			path + filepath.FromSlash("/article/article.go"),
			path + filepath.FromSlash("/article/article_new.go"),
			path + filepath.FromSlash("/article/article_v1_create.go"),
			path + filepath.FromSlash("/article/article_v1_get_list.go"),
			path + filepath.FromSlash("/article/article_v1_get_one.go"),
			path + filepath.FromSlash("/article/article_v1_update.go"),
			path + filepath.FromSlash("/article/article_v2_create.go"),
			path + filepath.FromSlash("/article/article_v2_update.go"),
		})

		// content
		testPath := gtest.DataPath("genctrl", "controller")
		expectFiles := []string{
			testPath + filepath.FromSlash("/article/article.go"),
			testPath + filepath.FromSlash("/article/article_new.go"),
			testPath + filepath.FromSlash("/article/article_v1_create.go"),
			testPath + filepath.FromSlash("/article/article_v1_get_list.go"),
			testPath + filepath.FromSlash("/article/article_v1_get_one.go"),
			testPath + filepath.FromSlash("/article/article_v1_update.go"),
			testPath + filepath.FromSlash("/article/article_v2_create.go"),
			testPath + filepath.FromSlash("/article/article_v2_update.go"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}

// https://github.com/gogf/gf/issues/3460
func Test_Gen_Ctrl_UseMerge_Issue3460(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctrlPath = gfile.Temp(guid.S())
			//ctrlPath  = gtest.DataPath("issue", "3460", "controller")
			apiFolder = gtest.DataPath("issue", "3460", "api")
			in        = genctrl.CGenCtrlInput{
				SrcFolder:     apiFolder,
				DstFolder:     ctrlPath,
				WatchFile:     "",
				SdkPath:       "",
				SdkStdVersion: false,
				SdkNoV1:       false,
				Clear:         false,
				Merge:         true,
			}
		)

		err := gfile.Mkdir(ctrlPath)
		t.AssertNil(err)
		defer gfile.Remove(ctrlPath)

		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		if err != nil {
			panic(err)
		}

		files, err := gfile.ScanDir(ctrlPath, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			filepath.Join(ctrlPath, "/hello/hello.go"),
			filepath.Join(ctrlPath, "/hello/hello_new.go"),
			filepath.Join(ctrlPath, "/hello/hello_v1_req.go"),
			filepath.Join(ctrlPath, "/hello/hello_v2_req.go"),
		})

		expectCtrlPath := gtest.DataPath("issue", "3460", "controller")
		expectFiles := []string{
			filepath.Join(expectCtrlPath, "/hello/hello.go"),
			filepath.Join(expectCtrlPath, "/hello/hello_new.go"),
			filepath.Join(expectCtrlPath, "/hello/hello_v1_req.go"),
			filepath.Join(expectCtrlPath, "/hello/hello_v2_req.go"),
		}
		// Line Feed maybe \r\n or \n

		for i, expectFile := range expectFiles {
			val := gfile.GetContents(files[i])
			expect := gfile.GetContents(expectFile)
			t.Assert(val, expect)
		}
	})
}
