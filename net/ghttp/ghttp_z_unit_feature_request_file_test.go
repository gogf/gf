// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Params_File_Single(t *testing.T) {
	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/single", func(r *ghttp.Request) {
		file := r.GetUploadFile("file")
		if file == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}

		if name, err := file.Save(dstDirPath, r.Get("randomlyRename").Bool()); err == nil {
			r.Response.WriteExit(name)
		}
		r.Response.WriteExit("upload failed")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "file1.txt")
		content := client.PostContent(ctx, "/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "file1.txt")
		t.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
	// randomly rename.
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file2.txt")
		content := client.PostContent(ctx, "/upload/single", g.Map{
			"file":           "@file:" + srcPath,
			"randomlyRename": true,
		})
		dstPath := gfile.Join(dstDirPath, content)
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

func Test_Params_File_CustomName(t *testing.T) {
	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/single", func(r *ghttp.Request) {
		file := r.GetUploadFile("file")
		if file == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}
		file.Filename = "my.txt"
		if name, err := file.Save(dstDirPath, r.Get("randomlyRename").Bool()); err == nil {
			r.Response.WriteExit(name)
		}
		r.Response.WriteExit("upload failed")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "my.txt")
		content := client.PostContent(ctx, "/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "my.txt")
		t.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

func Test_Params_File_Batch(t *testing.T) {
	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/batch", func(r *ghttp.Request) {
		files := r.GetUploadFiles("file")
		if files == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}
		if names, err := files.Save(dstDirPath, r.Get("randomlyRename").Bool()); err == nil {
			r.Response.WriteExit(gstr.Join(names, ","))
		}
		r.Response.WriteExit("upload failed")
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath1 := gtest.DataPath("upload", "file1.txt")
		srcPath2 := gtest.DataPath("upload", "file2.txt")
		dstPath1 := gfile.Join(dstDirPath, "file1.txt")
		dstPath2 := gfile.Join(dstDirPath, "file2.txt")
		content := client.PostContent(ctx, "/upload/batch", g.Map{
			"file[0]": "@file:" + srcPath1,
			"file[1]": "@file:" + srcPath2,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "file1.txt,file2.txt")
		t.Assert(gfile.GetContents(dstPath1), gfile.GetContents(srcPath1))
		t.Assert(gfile.GetContents(dstPath2), gfile.GetContents(srcPath2))
	})
	// randomly rename.
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath1 := gtest.DataPath("upload", "file1.txt")
		srcPath2 := gtest.DataPath("upload", "file2.txt")
		content := client.PostContent(ctx, "/upload/batch", g.Map{
			"file[0]":        "@file:" + srcPath1,
			"file[1]":        "@file:" + srcPath2,
			"randomlyRename": true,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")

		array := gstr.SplitAndTrim(content, ",")
		t.Assert(len(array), 2)
		dstPath1 := gfile.Join(dstDirPath, array[0])
		dstPath2 := gfile.Join(dstDirPath, array[1])
		t.Assert(gfile.GetContents(dstPath1), gfile.GetContents(srcPath1))
		t.Assert(gfile.GetContents(dstPath2), gfile.GetContents(srcPath2))
	})
}

func Test_Params_Strict_Route_File_Single_Ptr_Attrr(t *testing.T) {
	type Req struct {
		gmeta.Meta `method:"post" mime:"multipart/form-data"`
		File       *ghttp.UploadFile `type:"file"`
	}
	type Res struct{}

	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/single", func(ctx context.Context, req *Req) (res *Res, err error) {
		var (
			r    = g.RequestFromCtx(ctx)
			file = req.File
		)
		if file == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}
		name, err := file.Save(dstDirPath)
		if err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(name)
		return
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "file1.txt")
		content := client.PostContent(ctx, "/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "file1.txt")
		t.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

func Test_Params_Strict_Route_File_Single_Struct_Attr(t *testing.T) {
	type Req struct {
		gmeta.Meta `method:"post" mime:"multipart/form-data"`
		File       ghttp.UploadFile `type:"file"`
	}
	type Res struct{}

	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/single", func(ctx context.Context, req *Req) (res *Res, err error) {
		var (
			r    = g.RequestFromCtx(ctx)
			file = req.File
		)
		name, err := file.Save(dstDirPath)
		if err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(name)
		return
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "file1.txt")
		content := client.PostContent(ctx, "/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "file1.txt")
		t.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

func Test_Params_File_Upload_Required(t *testing.T) {
	type Req struct {
		gmeta.Meta `method:"post" mime:"multipart/form-data"`
		File       *ghttp.UploadFile `type:"file" v:"required#upload file is required"`
	}
	type Res struct{}

	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.BindHandler("/upload/required", func(ctx context.Context, req *Req) (res *Res, err error) {
		return
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// file is empty
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		content := client.PostContent(ctx, "/upload/required")
		t.Assert(content, `{"code":51,"message":"upload file is required","data":null}`)
	})
}

func Test_Params_File_MarshalJSON(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/upload/single", func(r *ghttp.Request) {
		file := r.GetUploadFile("file")
		if file == nil {
			r.Response.WriteExit("upload file cannot be empty")
		}

		if bytes, err := json.Marshal(file); err != nil {
			r.Response.WriteExit(err)
		} else {
			r.Response.WriteExit(bytes)
		}
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file1.txt")
		content := client.PostContent(ctx, "/upload/single", g.Map{
			"file": "@file:" + srcPath,
		})
		t.Assert(strings.Contains(content, "file1.txt"), true)
	})
}

// Select only one file when batch uploading
func Test_Params_Strict_Route_File_Batch_Up_One(t *testing.T) {
	type Req struct {
		gmeta.Meta `method:"post" mime:"multipart/form-data"`
		Files      ghttp.UploadFiles `type:"file"`
	}
	type Res struct{}

	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/batch", func(ctx context.Context, req *Req) (res *Res, err error) {
		var (
			r     = g.RequestFromCtx(ctx)
			files = req.Files
		)
		if len(files) == 0 {
			r.Response.WriteExit("upload file cannot be empty")
		}
		names, err := files.Save(dstDirPath)
		if err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(gstr.Join(names, ","))
		return
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath := gtest.DataPath("upload", "file1.txt")
		dstPath := gfile.Join(dstDirPath, "file1.txt")
		content := client.PostContent(ctx, "/upload/batch", g.Map{
			"files": "@file:" + srcPath,
		})
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "file1.txt")
		t.Assert(gfile.GetContents(dstPath), gfile.GetContents(srcPath))
	})
}

// Select multiple files during batch upload
func Test_Params_Strict_Route_File_Batch_Up_Multiple(t *testing.T) {
	type Req struct {
		gmeta.Meta `method:"post" mime:"multipart/form-data"`
		Files      ghttp.UploadFiles `type:"file"`
	}
	type Res struct{}

	dstDirPath := gfile.Temp(gtime.TimestampNanoStr())
	s := g.Server(guid.S())
	s.BindHandler("/upload/batch", func(ctx context.Context, req *Req) (res *Res, err error) {
		var (
			r     = g.RequestFromCtx(ctx)
			files = req.Files
		)
		if len(files) == 0 {
			r.Response.WriteExit("upload file cannot be empty")
		}
		names, err := files.Save(dstDirPath)
		if err != nil {
			r.Response.WriteExit(err)
		}
		r.Response.WriteExit(gstr.Join(names, ","))
		return
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()
	time.Sleep(100 * time.Millisecond)
	// normal name
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))

		srcPath1 := gtest.DataPath("upload", "file1.txt")
		srcPath2 := gtest.DataPath("upload", "file2.txt")
		dstPath1 := gfile.Join(dstDirPath, "file1.txt")
		dstPath2 := gfile.Join(dstDirPath, "file2.txt")
		content := client.PostContent(ctx, "/upload/batch",
			"files=@file:"+srcPath1+
				"&files=@file:"+srcPath2,
		)
		t.AssertNE(content, "")
		t.AssertNE(content, "upload file cannot be empty")
		t.AssertNE(content, "upload failed")
		t.Assert(content, "file1.txt,file2.txt")
		t.Assert(gfile.GetContents(dstPath1), gfile.GetContents(srcPath1))
		t.Assert(gfile.GetContents(dstPath2), gfile.GetContents(srcPath2))
	})
}
