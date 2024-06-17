// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"path/filepath"
	"testing"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genservice"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_Gen_Service_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path      = gfile.Temp(guid.S())
			dstFolder = path + filepath.FromSlash("/service")
			apiFolder = gtest.DataPath("genservice", "logic")
			in        = genservice.CGenServiceInput{
				SrcFolder:       apiFolder,
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
			genApi       = apiFolder + filepath.FromSlash("/logic.go")
			genApiExpect = apiFolder + filepath.FromSlash("/logic_expect.go")
		)
		defer gfile.Remove(genApi)
		t.Assert(gfile.GetContents(genApi), gfile.GetContents(genApiExpect))

		// files
		files, err := gfile.ScanDir(dstFolder, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			dstFolder + filepath.FromSlash("/article.go"),
			dstFolder + filepath.FromSlash("/delivery.go"),
			dstFolder + filepath.FromSlash("/user.go"),
		})

		// contents
		testPath := gtest.DataPath("genservice", "service")
		expectFiles := []string{
			testPath + filepath.FromSlash("/article.go"),
			testPath + filepath.FromSlash("/delivery.go"),
			testPath + filepath.FromSlash("/user.go"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
