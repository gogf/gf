// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gproc"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	Fix = cFix{}
)

type cFix struct {
	g.Meta `name:"fix" brief:"auto fixing codes after upgrading to new GoFrame version" usage:"gf fix" `
}

type cFixInput struct {
	g.Meta  `name:"fix"`
	Path    string `name:"path"    short:"p" brief:"directory path, it uses current working directory in default"`
	Version string `name:"version" short:"v" brief:"custom specified version to fix, leave it empty to auto detect"`
}

type cFixOutput struct{}

type cFixItem struct {
	Version string
	Func    func(version string) error
}

func (c cFix) Index(ctx context.Context, in cFixInput) (out *cFixOutput, err error) {

	if in.Path == "" {
		in.Path = gfile.Pwd()
	}
	if in.Version == "" {
		in.Version, err = c.autoDetectVersion(in)
		if err != nil {
			mlog.Fatal(err)
		}
		if in.Version == "" {
			mlog.Print(`no GoFrame usage found, exit fixing`)
			return
		}
		mlog.Debugf(`current GoFrame version auto detect "%s"`, in.Version)
	}

	if !gproc.IsChild() {
		mlog.Printf(`start auto fixing directory path "%s"...`, in.Path)
		defer mlog.Print(`done!`)
	}

	err = c.doFix(in)
	return
}

func (c cFix) doFix(in cFixInput) (err error) {

	var items = []cFixItem{
		{Version: "v2.3", Func: c.doFixV23},
	}
	for _, item := range items {
		if gstr.CompareVersionGo(in.Version, item.Version) < 0 {
			mlog.Debugf(
				`current GoFrame or contrib package version "%s" is lesser than "%s", nothing to do`,
				in.Version, item.Version,
			)
			continue
		}
		if err = item.Func(in.Version); err != nil {
			return
		}
	}
	return
}

// doFixV23 fixes code when upgrading to GoFrame v2.3.
func (c cFix) doFixV23(version string) error {
	replaceFunc := func(path, content string) string {
		// gdb.TX from struct to interface.
		content = gstr.Replace(content, "*gdb.TX", "gdb.TX")
		// function name changes for package gtcp/gudp.
		if gstr.Contains(content, "/gf/v2/net/gtcp") || gstr.Contains(content, "/gf/v2/net/gudp") {
			content = gstr.ReplaceByMap(content, g.MapStrStr{
				".SetSendDeadline":      ".SetDeadlineSend",
				".SetReceiveDeadline":   ".SetDeadlineRecv",
				".SetReceiveBufferWait": ".SetBufferWaitRecv",
			})
		}
		return content
	}
	return gfile.ReplaceDirFunc(replaceFunc, ".", "*.go", true)
}

func (c cFix) autoDetectVersion(in cFixInput) (string, error) {
	var (
		err     error
		path    = gfile.Join(in.Path, "go.mod")
		version string
	)
	if !gfile.Exists(path) {
		return "", gerror.Newf(`"%s" not found in current working directory`, path)
	}
	err = gfile.ReadLines(path, func(line string) error {
		array := gstr.SplitAndTrim(line, " ")
		if len(array) > 0 {
			if gstr.HasPrefix(array[0], gfPackage) {
				version = array[1]
			}
		}
		return nil
	})
	if err != nil {
		mlog.Fatal(err)
	}
	return version, nil
}
