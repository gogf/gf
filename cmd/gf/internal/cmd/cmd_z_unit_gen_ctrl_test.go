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
		t.AssertNil(err)

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

// gf gen ctrl -m
// In the same module, different API files are added
func Test_Gen_Ctrl_UseMerge_AddNewFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctrlPath = gfile.Temp(guid.S())
			//ctrlPath  = gtest.DataPath("issue", "3460", "controller")
			apiFolder = gtest.DataPath("genctrl-merge", "add_new_file", "api")
			in        = genctrl.CGenCtrlInput{
				SrcFolder: apiFolder,
				DstFolder: ctrlPath,
				Merge:     true,
			}
		)
		const testNewApiFile = `
package v1
import "github.com/gogf/gf/v2/frame/g"
type DictTypeAddReq struct {
	g.Meta
}
type DictTypeAddRes struct {
}
`

		err := gfile.Mkdir(ctrlPath)
		t.AssertNil(err)
		defer gfile.Remove(ctrlPath)

		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		var (
			genApi       = filepath.Join(apiFolder, "/dict/dict.go")
			genApiExpect = filepath.Join(apiFolder, "/dict/dict_expect.go")
		)
		defer gfile.Remove(genApi)
		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))

		genCtrlFiles, err := gfile.ScanDir(ctrlPath, "*.go", true)
		t.AssertNil(err)
		t.Assert(genCtrlFiles, []string{
			filepath.Join(ctrlPath, "/dict/dict.go"),
			filepath.Join(ctrlPath, "/dict/dict_new.go"),
			filepath.Join(ctrlPath, "/dict/dict_v1_dict_type.go"),
		})

		expectCtrlPath := gtest.DataPath("genctrl-merge", "add_new_file", "controller")
		expectFiles := []string{
			filepath.Join(expectCtrlPath, "/dict/dict.go"),
			filepath.Join(expectCtrlPath, "/dict/dict_new.go"),
			filepath.Join(expectCtrlPath, "/dict/dict_v1_dict_type.go"),
		}

		// Line Feed maybe \r\n or \n
		expectFilesContent(t, genCtrlFiles, expectFiles)

		// Add a new API file
		newApiFilePath := filepath.Join(apiFolder, "/dict/v1/test_new.go")
		err = gfile.PutContents(newApiFilePath, testNewApiFile)
		t.AssertNil(err)
		defer gfile.Remove(newApiFilePath)

		// Then execute the command
		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		genApi = filepath.Join(apiFolder, "/dict.go")
		genApiExpect = filepath.Join(apiFolder, "/dict_add_new_ctrl_expect.gotest")

		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))

		genCtrlFiles = append(genCtrlFiles, filepath.Join(ctrlPath, "/dict/dict_v1_test_new.go"))
		// Use the gotest suffix, otherwise the IDE will delete the import
		expectFiles = append(expectFiles, filepath.Join(expectCtrlPath, "/dict/dict_v1_test_new.gotest"))
		// Line Feed maybe \r\n or \n
		expectFilesContent(t, genCtrlFiles, expectFiles)

	})

}

// gf gen ctrl -m
// In the same module, Add the same file to the API
func Test_Gen_Ctrl_UseMerge_AddNewCtrl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctrlPath = gfile.Temp(guid.S())
			//ctrlPath  = gtest.DataPath("issue", "3460", "controller")
			apiFolder = gtest.DataPath("genctrl-merge", "add_new_ctrl", "api")
			in        = genctrl.CGenCtrlInput{
				SrcFolder: apiFolder,
				DstFolder: ctrlPath,
				Merge:     true,
			}
		)

		err := gfile.Mkdir(ctrlPath)
		t.AssertNil(err)
		defer gfile.Remove(ctrlPath)

		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		var (
			genApi       = filepath.Join(apiFolder, "/dict/dict.go")
			genApiExpect = filepath.Join(apiFolder, "/dict/dict_expect.go")
		)
		defer gfile.Remove(genApi)
		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))

		genCtrlFiles, err := gfile.ScanDir(ctrlPath, "*.go", true)
		t.AssertNil(err)
		t.Assert(genCtrlFiles, []string{
			filepath.Join(ctrlPath, "/dict/dict.go"),
			filepath.Join(ctrlPath, "/dict/dict_new.go"),
			filepath.Join(ctrlPath, "/dict/dict_v1_dict_type.go"),
		})

		expectCtrlPath := gtest.DataPath("genctrl-merge", "add_new_ctrl", "controller")
		expectFiles := []string{
			filepath.Join(expectCtrlPath, "/dict/dict.go"),
			filepath.Join(expectCtrlPath, "/dict/dict_new.go"),
			filepath.Join(expectCtrlPath, "/dict/dict_v1_dict_type.go"),
		}

		// Line Feed maybe \r\n or \n
		expectFilesContent(t, genCtrlFiles, expectFiles)

		const testNewApiFile = `

type DictTypeAddReq struct {
	g.Meta
}
type DictTypeAddRes struct {
}
`
		dictModuleFileName := filepath.Join(apiFolder, "/dict/v1/dict_type.go")
		// Save the contents of the file before the changes
		apiFileContents := gfile.GetContents(dictModuleFileName)

		// Add a new API file
		err = gfile.PutContentsAppend(dictModuleFileName, testNewApiFile)
		t.AssertNil(err)

		//==================================
		// Then execute the command
		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		genApi = filepath.Join(apiFolder, "/dict.go")
		genApiExpect = filepath.Join(apiFolder, "/dict_add_new_ctrl_expect.gotest")
		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))

		// Use the gotest suffix, otherwise the IDE will delete the import
		expectFiles[2] = filepath.Join(expectCtrlPath, "/dict/dict_v1_test_new.gotest")
		// Line Feed maybe \r\n or \n
		expectFilesContent(t, genCtrlFiles, expectFiles)

		// Restore the contents of the original API file
		err = gfile.PutContents(dictModuleFileName, apiFileContents)
		t.AssertNil(err)
	})

}

func expectFilesContent(t *gtest.T, paths []string, expectPaths []string) {
	for i, expectFile := range expectPaths {
		val := gfile.GetContents(paths[i])
		expect := gfile.GetContents(expectFile)
		t.Assert(val, expect)
	}
}

// https://github.com/gogf/gf/issues/3569
func Test_Gen_Ctrl_Comments_Issue3569(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctrlPath  = gfile.Temp(guid.S())
			apiFolder = gtest.DataPath("issue", "3569", "api")
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

		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(ctrlPath)
		t.AssertNil(err)
		defer gfile.Remove(ctrlPath)

		_, err = genctrl.CGenCtrl{}.Ctrl(ctx, in)
		t.AssertNil(err)

		//apiInterface file
		var (
			genApi       = apiFolder + filepath.FromSlash("/hello/hello.go")
			genApiExpect = apiFolder + filepath.FromSlash("/hello/hello_expect.go")
		)
		defer gfile.Remove(genApi)
		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))
	})
}
