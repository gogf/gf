// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/test/gtest"
)

func Test_Params_File_Single(t *testing.T) {
	dstDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/upload/single", func(r *ghttp.Request) {
		file := r.GetUploadFile("file")
		if file == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}

		if name, err := file.Save(dstDirPath, r.GetBool("randomlyRename")); err == nil {
			r.Response.WriteExit(name)
		}
		r.Response.WriteExit("upload failed")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		srcPath := gfile.Join(gdebug.TestDataPath(), "upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "file1.txt")
		content := client.PostContent("/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		gtest.AssertNE(content, "")
		gtest.AssertNE(content, "upload file cannot be empty")
		gtest.AssertNE(content, "upload failed")
		gtest.Assert(content, "file1.txt")
		gtest.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
	// randomly rename.
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		srcPath := gfile.Join(gdebug.TestDataPath(), "upload", "file2.txt")
		content := client.PostContent("/upload/single", g.Map{
			"file":           "@file:" + srcPath,
			"randomlyRename": true,
		})
		dstPath := gfile.Join(dstDirPath, content)
		gtest.AssertNE(content, "")
		gtest.AssertNE(content, "upload file cannot be empty")
		gtest.AssertNE(content, "upload failed")
		gtest.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

func Test_Params_File_CustomName(t *testing.T) {
	dstDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/upload/single", func(r *ghttp.Request) {
		file := r.GetUploadFile("file")
		if file == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}
		file.Filename = "my.txt"
		if name, err := file.Save(dstDirPath, r.GetBool("randomlyRename")); err == nil {
			r.Response.WriteExit(name)
		}
		r.Response.WriteExit("upload failed")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		srcPath := gfile.Join(gdebug.TestDataPath(), "upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "my.txt")
		content := client.PostContent("/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		gtest.AssertNE(content, "")
		gtest.AssertNE(content, "upload file cannot be empty")
		gtest.AssertNE(content, "upload failed")
		gtest.Assert(content, "my.txt")
		gtest.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

func Test_Params_File_Batch(t *testing.T) {
	dstDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
	p := ports.PopRand()
	s := g.Server(p)
	s.BindHandler("/upload/batch", func(r *ghttp.Request) {
		files := r.GetUploadFiles("file")
		if files == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}
		if names, err := files.Save(dstDirPath, r.GetBool("randomlyRename")); err == nil {
			r.Response.WriteExit(gstr.Join(names, ","))
		}
		r.Response.WriteExit("upload failed")
	})
	s.SetPort(p)
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		srcPath1 := gfile.Join(gdebug.TestDataPath(), "upload", "file1.txt")
		srcPath2 := gfile.Join(gdebug.TestDataPath(), "upload", "file2.txt")
		dstPath1 := gfile.Join(dstDirPath, "file1.txt")
		dstPath2 := gfile.Join(dstDirPath, "file2.txt")
		content := client.PostContent("/upload/batch", g.Map{
			"file[0]": "@file:" + srcPath1,
			"file[1]": "@file:" + srcPath2,
		})
		gtest.AssertNE(content, "")
		gtest.AssertNE(content, "upload file cannot be empty")
		gtest.AssertNE(content, "upload failed")
		gtest.Assert(content, "file1.txt,file2.txt")
		gtest.Assert(gfile.GetContents(dstPath1), gfile.GetContents(srcPath1))
		gtest.Assert(gfile.GetContents(dstPath2), gfile.GetContents(srcPath2))
	})
	// randomly rename.
	gtest.Case(t, func() {
		client := ghttp.NewClient()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", p))

		srcPath1 := gfile.Join(gdebug.TestDataPath(), "upload", "file1.txt")
		srcPath2 := gfile.Join(gdebug.TestDataPath(), "upload", "file2.txt")
		content := client.PostContent("/upload/batch", g.Map{
			"file[0]":        "@file:" + srcPath1,
			"file[1]":        "@file:" + srcPath2,
			"randomlyRename": true,
		})
		gtest.AssertNE(content, "")
		gtest.AssertNE(content, "upload file cannot be empty")
		gtest.AssertNE(content, "upload failed")

		array := gstr.SplitAndTrim(content, ",")
		gtest.Assert(len(array), 2)
		dstPath1 := gfile.Join(dstDirPath, array[0])
		dstPath2 := gfile.Join(dstDirPath, array[1])
		gtest.Assert(gfile.GetContents(dstPath1), gfile.GetContents(srcPath1))
		gtest.Assert(gfile.GetContents(dstPath2), gfile.GetContents(srcPath2))
	})
}
