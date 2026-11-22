// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress_test

import (
	"archive/zip"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_UnZipContent_PathTraversal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		newZip := func(name string) []byte {
			var buffer bytes.Buffer
			writer := zip.NewWriter(&buffer)
			fileWriter, err := writer.Create(name)
			t.AssertNil(err)
			_, err = fileWriter.Write([]byte("evil"))
			t.AssertNil(err)
			t.AssertNil(writer.Close())
			return buffer.Bytes()
		}

		for _, name := range []string{
			"../evil.txt",
			"a/../../evil.txt",
			`..\evil.txt`,
		} {
			root := t.TempDir()
			dst := gfile.Join(root, "out")
			err := gcompress.UnZipContent(newZip(name), dst)
			t.AssertNE(err, nil)
			t.Assert(gerror.Code(err).Code(), gcode.CodeInvalidParameter.Code())
			t.Assert(gfile.Exists(gfile.Join(root, "evil.txt")), false)
		}
	})
}

func Test_UnZipContent_SymlinkEscape(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := t.TempDir()
		dst := gfile.Join(root, "out")
		outside := gfile.Join(root, "outside")
		t.AssertNil(gfile.Mkdir(dst))
		t.AssertNil(gfile.Mkdir(outside))
		if err := os.Symlink(filepath.Join("..", "outside"), gfile.Join(dst, "link")); err != nil {
			t.Skipf("symlink is not supported: %v", err)
		}

		var buffer bytes.Buffer
		writer := zip.NewWriter(&buffer)
		fileWriter, err := writer.Create("link/evil.txt")
		t.AssertNil(err)
		_, err = fileWriter.Write([]byte("evil"))
		t.AssertNil(err)
		t.AssertNil(writer.Close())

		err = gcompress.UnZipContent(buffer.Bytes(), dst)
		t.AssertNE(err, nil)
		t.Assert(gfile.Exists(gfile.Join(outside, "evil.txt")), false)

		insideDst := gfile.Join(root, "inside")
		t.AssertNil(gfile.Mkdir(insideDst))
		t.AssertNil(gfile.Mkdir(gfile.Join(insideDst, "real")))
		if err = os.Symlink("real", gfile.Join(insideDst, "link")); err != nil {
			t.Skipf("symlink is not supported: %v", err)
		}
		buffer.Reset()
		writer = zip.NewWriter(&buffer)
		fileWriter, err = writer.Create("link/ok.txt")
		t.AssertNil(err)
		_, err = fileWriter.Write([]byte("ok"))
		t.AssertNil(err)
		t.AssertNil(writer.Close())

		t.AssertNil(gcompress.UnZipContent(buffer.Bytes(), insideDst))
		t.Assert(gfile.GetContents(gfile.Join(insideDst, "real", "ok.txt")), "ok")
	})
}

func Test_UnZipContent_NormalAndPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var buffer bytes.Buffer
		writer := zip.NewWriter(&buffer)
		fileWriter, err := writer.Create("dir/file.txt")
		t.AssertNil(err)
		_, err = fileWriter.Write([]byte("ok"))
		t.AssertNil(err)
		t.AssertNil(writer.Close())

		root := t.TempDir()
		dst := gfile.Join(root, "out")
		t.AssertNil(gcompress.UnZipContent(buffer.Bytes(), dst))
		t.Assert(gfile.GetContents(gfile.Join(dst, "dir", "file.txt")), "ok")

		prefixDst := gfile.Join(root, "prefix")
		t.AssertNil(gcompress.UnZipContent(buffer.Bytes(), prefixDst, "dir"))
		t.Assert(gfile.GetContents(gfile.Join(prefixDst, "file.txt")), "ok")
	})
}
