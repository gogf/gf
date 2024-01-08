// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genpb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
	"path/filepath"
	"testing"
)

func Test_Gen_Pb_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			path              = gfile.Temp(guid.S())
			genApiPath        = path + filepath.FromSlash("/api")
			genControllerPath = path + filepath.FromSlash("/controller")

			in = genpb.CGenPbInput{
				Path:       gtest.DataPath("genpb", "protobuf"),
				OutputApi:  genApiPath,
				OutputCtrl: genControllerPath,
			}
		)
		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		err = gfile.Mkdir(path)
		t.AssertNil(err)
		err = gfile.Mkdir(genApiPath)
		t.AssertNil(err)
		err = gfile.Mkdir(genControllerPath)
		t.AssertNil(err)
		defer gfile.Remove(path)

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		t.AssertNil(err)

		// files
		files, err := gfile.ScanDir(path, "*.go", true)
		t.AssertNil(err)
		t.Assert(files, []string{
			path + filepath.FromSlash("/api/article/v1/article.pb.go"),
			path + filepath.FromSlash("/api/article/v1/article_grpc.pb.go"),
			path + filepath.FromSlash("/controller/article/article.go"),
		})

		// content
		testPath := gtest.DataPath("genpb")
		expectFiles := []string{
			testPath + filepath.FromSlash("/api/article/v1/article.pb.go.txt"),
			testPath + filepath.FromSlash("/api/article/v1/article_grpc.pb.go.txt"),
			testPath + filepath.FromSlash("/controller/article/article.go.txt"),
		}
		for i := range files {
			t.Assert(gfile.GetContents(files[i]), gfile.GetContents(expectFiles[i]))
		}
	})
}
