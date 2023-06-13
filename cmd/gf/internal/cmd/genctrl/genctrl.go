// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	CGenCtrlConfig = `gfcli.gen.ctrl`
	CGenCtrlUsage  = `gf gen ctrl [OPTION]`
	CGenCtrlBrief  = `parse struct and associated functions from packages to generate ctrl go file`
	CGenCtrlEg     = `
gf gen ctrl
`
	CGenCtrlBriefSrcFolder = `source folder path to be parsed. default: internal/logic`
	CGenCtrlBriefDstFolder = `destination folder path storing automatically generated go files. default: internal/ctrl`
	CGenCtrlBriefWatchFile = `used in file watcher, it re-generates all ctrl go files only if given file is under srcFolder`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenCtrlConfig`:         CGenCtrlConfig,
		`CGenCtrlUsage`:          CGenCtrlUsage,
		`CGenCtrlBrief`:          CGenCtrlBrief,
		`CGenCtrlEg`:             CGenCtrlEg,
		`CGenCtrlBriefSrcFolder`: CGenCtrlBriefSrcFolder,
		`CGenCtrlBriefDstFolder`: CGenCtrlBriefDstFolder,
		`CGenCtrlBriefWatchFile`: CGenCtrlBriefWatchFile,
	})
}

type (
	CGenCtrl      struct{}
	CGenCtrlInput struct {
		g.Meta    `name:"ctrl" config:"{CGenCtrlConfig}" usage:"{CGenCtrlUsage}" brief:"{CGenCtrlBrief}" eg:"{CGenCtrlEg}"`
		SrcFolder string `short:"s" name:"srcFolder" brief:"{CGenCtrlBriefSrcFolder}" d:"api"`
		DstFolder string `short:"d" name:"dstFolder" brief:"{CGenCtrlBriefDstFolder}" d:"internal/controller"`
		WatchFile string `short:"w" name:"watchFile" brief:"{CGenCtrlBriefWatchFile}"`
	}
	CGenCtrlOutput struct{}
)

const (
	genCtrlFileLockSeconds = 1
)

func (c CGenCtrl) Ctrl(ctx context.Context, in CGenCtrlInput) (out *CGenCtrlOutput, err error) {
	in.SrcFolder = "/Users/txqiangguo/Workspace/eros/app/khaos-shark/api"
	in.DstFolder = "/Users/txqiangguo/Workspace/eros/app/khaos-shark/internal/controller"
	// File lock to avoid multiple processes.
	var (
		flockFilePath = gfile.Temp("gf.cli.gen.ctrl.lock")
		flockContent  = gfile.GetContents(flockFilePath)
	)
	if flockContent != "" {
		if gtime.Timestamp()-gconv.Int64(flockContent) < genCtrlFileLockSeconds {
			// If another "gen ctrl" process is running, it just exits.
			mlog.Debug(`another "gen ctrl" process is running, exit`)
			return
		}
	}
	defer gfile.Remove(flockFilePath)
	_ = gfile.PutContents(flockFilePath, gtime.TimestampStr())

	in.SrcFolder = gstr.TrimRight(in.SrcFolder, `\/`)
	in.SrcFolder = gstr.Replace(in.SrcFolder, "\\", "/")
	in.WatchFile = gstr.TrimRight(in.WatchFile, `\/`)
	in.WatchFile = gstr.Replace(in.WatchFile, "\\", "/")

	// Watch file handling.
	if in.WatchFile != "" {
		// It works only if given WatchFile is in SrcFolder.
		var (
			watchFileDir = gfile.Dir(in.WatchFile)
			srcFolderDir = gfile.Dir(watchFileDir)
		)
		mlog.Debug("watchFileDir:", watchFileDir)
		mlog.Debug("logicFolderDir:", srcFolderDir)
		if !gstr.HasSuffix(gstr.Replace(srcFolderDir, `\`, `/`), in.SrcFolder) {
			mlog.Printf(`ignore watch file "%s", not in source path "%s"`, in.WatchFile, in.SrcFolder)
			return
		}
		var newWorkingDir = gfile.Dir(gfile.Dir(srcFolderDir))
		if err = gfile.Chdir(newWorkingDir); err != nil {
			mlog.Fatalf(`%+v`, err)
		}
		mlog.Debug("Chdir:", newWorkingDir)
		_ = gfile.Remove(flockFilePath)
		var command = fmt.Sprintf(
			`%s gen ctrl -packages=%s`,
			gfile.SelfName(), gfile.Basename(watchFileDir),
		)
		err = gproc.ShellRun(ctx, command)
		return
	}

	if !gfile.Exists(in.SrcFolder) {
		mlog.Fatalf(`source folder path "%s" does not exist`, in.SrcFolder)
	}

	apiItemsInSrc, err := c.getApiItemsInSrc(in.SrcFolder)
	if err != nil {
		return nil, err
	}
	apiItemsInDst, err := c.getApiItemsInDst(in.DstFolder)
	if err != nil {
		return nil, err
	}

	// generate api interface.
	if err = newApiInterfaceGenerator().Generate(in.SrcFolder, apiItemsInSrc); err != nil {
		return
	}

	// api filtering for already implemented api controllers.
	var (
		alreadyImplementedCtrlSet = gset.NewStrSet()
		toBeImplementedApiItems   = make([]apiItem, 0)
	)
	for _, item := range apiItemsInDst {
		alreadyImplementedCtrlSet.Add(item.String())
	}
	for _, item := range apiItemsInSrc {
		if alreadyImplementedCtrlSet.Contains(item.String()) {
			continue
		}
		toBeImplementedApiItems = append(toBeImplementedApiItems, item)
	}

	// generate go files.
	if len(toBeImplementedApiItems) > 0 {
		err = newControllerGenerator().Generate(in.DstFolder, toBeImplementedApiItems)
		if err != nil {
			return
		}
	}

	mlog.Print(`done!`)
	return
}
