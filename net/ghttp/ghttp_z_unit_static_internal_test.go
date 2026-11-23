// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_OpenStaticFile_MissingRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rootPath := gfile.Join(t.TempDir(), "missing")
		file, err := openStaticFile(context.Background(), &staticFile{
			RootPath: rootPath,
			Path:     gfile.Join(rootPath, "file.txt"),
		})
		t.AssertNil(file)
		t.AssertNE(err, nil)
	})
}
