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
	"time"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	// DocURL is the download address of the document
	DocURL = "https://github.com/gogf/gf/archive/refs/heads/gh-pages.zip"
)

var (
	Doc = cDoc{}
)

type cDoc struct {
	g.Meta `name:"doc" brief:"download https://pages.goframe.org/ to run locally"`
}

type cDocInput struct {
	g.Meta `name:"doc" config:"gfcli.doc"`
	Path   string `short:"p"  name:"path"    brief:"download docs directory path, default is \"%temp%/goframe\""`
	Port   int    `short:"o"  name:"port"    brief:"http server port, default is 8080" d:"8080"`
	Update bool   `short:"u"  name:"update"  brief:"clean docs directory and update docs"`
	Clean  bool   `short:"c"  name:"clean"   brief:"clean docs directory"`
	Proxy  string `short:"x"  name:"proxy"   brief:"proxy for download, such as https://hub.gitmirror.com/;https://ghproxy.com/;https://ghproxy.net/;https://ghps.cc/"`
}

type cDocOutput struct{}

func (c cDoc) Index(ctx context.Context, in cDocInput) (out *cDocOutput, err error) {
	docs := NewDocSetting(ctx, in)
	mlog.Print("Directory where the document is downloaded:", docs.TempDir)
	if in.Clean {
		mlog.Print("Cleaning document directory")
		err = docs.Clean()
		if err != nil {
			mlog.Print("Failed to clean document directory:", err)
			return
		}
		return
	}
	if in.Update {
		mlog.Print("Cleaning old document directory")
		err = docs.Clean()
		if err != nil {
			mlog.Print("Failed to clean old document directory:", err)
			return
		}
	}
	err = docs.DownloadDoc()
	if err != nil {
		mlog.Print("Failed to download document:", err)
		return
	}
	s := g.Server()
	s.SetServerRoot(docs.DocDir)
	s.SetPort(in.Port)
	s.SetDumpRouterMap(false)
	mlog.Printf("Access address http://127.0.0.1:%d", in.Port)
	s.Run()
	return
}

// DocSetting doc setting
type DocSetting struct {
	TempDir    string
	DocURL     string
	DocDir     string
	DocZipFile string
}

// NewDocSetting new DocSetting
func NewDocSetting(ctx context.Context, in cDocInput) *DocSetting {
	fileName := "gf-doc-md.zip"
	tempDir := in.Path
	if tempDir == "" {
		tempDir = gfile.Temp("goframe/docs")
	} else {
		tempDir = gfile.Abs(path.Join(tempDir, "docs"))
	}

	return &DocSetting{
		TempDir:    filepath.FromSlash(tempDir),
		DocDir:     filepath.FromSlash(path.Join(tempDir, "gf-gh-pages")),
		DocURL:     in.Proxy + DocURL,
		DocZipFile: filepath.FromSlash(path.Join(tempDir, fileName)),
	}

}

// Clean clean the temporary directory
func (d *DocSetting) Clean() error {
	if _, err := os.Stat(d.TempDir); err == nil {
		err = gfile.Remove(d.TempDir)
		if err != nil {
			mlog.Print("Failed to delete temporary directory:", err)
			return err
		}
	}
	return nil
}

// DownloadDoc download the document
func (d *DocSetting) DownloadDoc() error {
	if _, err := os.Stat(d.TempDir); err != nil {
		err = gfile.Mkdir(d.TempDir)
		if err != nil {
			mlog.Print("Failed to create temporary directory:", err)
			return nil
		}
	}
	// Check if the file exists
	if _, err := os.Stat(d.DocDir); err == nil {
		mlog.Print("Document already exists, no need to download and unzip")
		return nil
	}

	if _, err := os.Stat(d.DocZipFile); err == nil {
		mlog.Print("File already exists, no need to download")
	} else {
		mlog.Printf("File does not exist, start downloading: %s", d.DocURL)
		startTime := time.Now()
		// Download the file
		resp, err := http.Get(d.DocURL)
		if err != nil {
			mlog.Print("Failed to download file:", err)
			return err
		}
		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(d.DocZipFile)
		if err != nil {
			mlog.Print("Failed to create file:", err)
			return err
		}
		defer out.Close()

		// Write the response body to the file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			mlog.Print("Failed to write file:", err)
			return err
		}
		mlog.Printf("Download successful, time-consuming: %v", time.Since(startTime))
	}

	mlog.Print("Start unzipping the file...")
	// Unzip the file
	err := gcompress.UnZipFile(d.DocZipFile, d.TempDir)
	if err != nil {
		mlog.Print("Failed to unzip the file, please run again:", err)
		gfile.Remove(d.DocZipFile)
		return err
	}

	mlog.Print("Download and unzip successful")
	return nil
}
