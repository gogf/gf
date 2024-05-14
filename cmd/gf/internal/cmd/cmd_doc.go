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
}

type cDocOutput struct{}

func (c cDoc) Index(ctx context.Context, in cDocInput) (out *cDocOutput, err error) {

	docDir := NewDocSetting().DownloadDoc()
	mlog.Print("文档目录:", docDir)
	s := g.Server()
	s.SetServerRoot(docDir)
	s.SetPort(8199)
	mlog.Print("http://127.0.0.1:8199")
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
func NewDocSetting() *DocSetting {
	fileName := "gf-doc-md.zip"
	tempDir := gfile.Temp("goframe")
	return &DocSetting{
		TempDir:    tempDir,
		DocDir:     path.Join(tempDir, "gf-gh-pages"),
		DocURL:     "https://codeload.github.com/gogf/gf/zip/refs/heads/gh-pages",
		DocZipFile: path.Join(tempDir, fileName),
	}
}

func (d *DocSetting) DownloadDoc() (docDir string) {
	docDir = d.DocDir
	// 判断文件是否存在
	if _, err := os.Stat(docDir); err == nil {
		mlog.Print("目录已存在，无需下载解压缩")
		return
	}

	if _, err := os.Stat(d.DocZipFile); err == nil {
		mlog.Print("文件已存在，无需下载")
	} else {
		mlog.Print("文件不存在，开始下载")
		// 下载文件
		resp, err := http.Get(d.DocURL)
		if err != nil {
			mlog.Print("下载文件失败:", err)
			return
		}
		defer resp.Body.Close()

		// 创建文件
		out, err := os.Create(d.DocZipFile)
		if err != nil {
			mlog.Print("创建文件失败:", err)
			return
		}
		defer out.Close()

		// 将响应体内容写入文件
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			mlog.Print("写入文件失败:", err)
			return
		}
	}

	mlog.Print("开始解压缩文件...")
	// 解压缩文件
	err := gcompress.UnZipFile(d.DocZipFile, d.TempDir)
	if err != nil {
		mlog.Print("解压缩文件失败:", err)
		return
	}

	mlog.Print("下载并解压缩成功")
	return
}
