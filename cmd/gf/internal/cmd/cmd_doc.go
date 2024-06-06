// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

var (
	Doc = cDoc{}
)

type cDoc struct {
	g.Meta `name:"doc" brief:"show current Golang environment variables"`
}

type cDocInput struct {
	g.Meta `name:"doc"`
	Path   string `short:"p"  name:"path"    brief:"download docs directory path, default is \"%temp%/goframe\""`
	Port   int    `short:"o"  name:"port"    brief:"http server port, default is 8080" d:"8080"`
	Update bool   `short:"u"  name:"update"  brief:"clean docs directory and update docs"`
	Clean  bool   `short:"c"  name:"clean"   brief:"clean docs directory"`
}

type cDocOutput struct{}

func (c cDoc) Index(ctx context.Context, in cDocInput) (out *cDocOutput, err error) {
	docs := NewDocSetting(in.Path)
	mlog.Print("下载文档所在目录:", docs.DocDir)
	if in.Clean {
		mlog.Print("清理文档目录")
		err = docs.Clean()
		if err != nil {
			mlog.Print("清理文档目录失败:", err)
			return
		}
		return
	}
	if in.Update {
		mlog.Print("清理旧文档目录")
		err = docs.Clean()
		if err != nil {
			mlog.Print("清理旧文档目录失败:", err)
			return
		}
	}
	err = docs.DownloadDoc()
	if err != nil {
		mlog.Print("下载文档失败:", err)
		return
	}
	s := g.Server()
	s.SetServerRoot(docs.DocDir)
	s.SetPort(in.Port)
	s.SetDumpRouterMap(false)
	mlog.Printf("访问地址 http://127.0.0.1:%d", in.Port)
	s.Run()
	return
}

// DocSetting description
type DocSetting struct {
	TempDir    string
	DocURL     string
	DocDir     string
	DocZipFile string
}

// NewDocSetting description
//
// createTime: 2024-05-14 12:19:55
func NewDocSetting(tempDir string) *DocSetting {
	fileName := "gf-doc-md.zip"
	if tempDir == "" {
		tempDir = gfile.Temp("goframe/docs")
	} else {
		tempDir = gfile.Abs(path.Join(tempDir, "docs"))
	}

	return &DocSetting{
		TempDir:    filepath.FromSlash(tempDir),
		DocDir:     filepath.FromSlash(path.Join(tempDir, "gf-gh-pages")),
		DocURL:     "https://codeload.github.com/gogf/gf/zip/refs/heads/gh-pages",
		DocZipFile: filepath.FromSlash(path.Join(tempDir, fileName)),
	}
}

func (d *DocSetting) Clean() error {
	if _, err := os.Stat(d.TempDir); err == nil {
		err = gfile.Remove(d.TempDir)
		if err != nil {
			mlog.Print("删除临时目录失败:", err)
			return err
		}
	}
	return nil
}

func (d *DocSetting) DownloadDoc() error {
	if _, err := os.Stat(d.TempDir); err != nil {
		err = gfile.Mkdir(d.TempDir)
		if err != nil {
			mlog.Print("创建临时目录失败:", err)
			return nil
		}
	}
	// 判断文件是否存在
	if _, err := os.Stat(d.DocDir); err == nil {
		mlog.Print("文档已存在，无需下载解压缩")
		return nil
	}

	if _, err := os.Stat(d.DocZipFile); err == nil {
		mlog.Print("文件已存在，无需下载")
	} else {
		mlog.Print("文件不存在，开始下载")
		// 下载文件
		resp, err := http.Get(d.DocURL)
		if err != nil {
			mlog.Print("下载文件失败:", err)
			return err
		}
		defer resp.Body.Close()

		// 创建文件
		out, err := os.Create(d.DocZipFile)
		if err != nil {
			mlog.Print("创建文件失败:", err)
			return err
		}
		defer out.Close()

		// 将响应体内容写入文件
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			mlog.Print("写入文件失败:", err)
			return err
		}
	}

	mlog.Print("开始解压缩文件...")
	// 解压缩文件
	err := gcompress.UnZipFile(d.DocZipFile, d.TempDir)
	if err != nil {
		mlog.Print("解压缩文件失败，请重新运行:", err)
		gfile.Remove(d.DocZipFile)
		return err
	}

	mlog.Print("下载并解压缩成功")
	return nil
}
