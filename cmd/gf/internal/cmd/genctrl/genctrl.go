// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"context"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	CGenCtrlConfig = `gfcli.gen.ctrl`
	CGenCtrlUsage  = `gf gen ctrl [OPTION]`
	CGenCtrlBrief  = `parse api definitions to generate ctrl go file`
	CGenCtrlEg     = `
gf gen ctrl
`
	CGenCtrlBriefSrcFolder = `source folder path to be parsed. default: api`
	CGenCtrlBriefDstFolder = `destination folder path storing automatically generated go files. default: internal/controller`
	CGenCtrlBriefWatchFile = `used in file watcher, it re-generates go files only if given file is under srcFolder`
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

	if !gfile.Exists(in.SrcFolder) {
		mlog.Fatalf(`source folder path "%s" does not exist`, in.SrcFolder)
	}
	// retrieve all api modules.
	apiModuleFolderPaths, err := gfile.ScanDir(in.SrcFolder, "*", false)
	if err != nil {
		return nil, err
	}
	for _, apiModuleFolderPath := range apiModuleFolderPaths {
		if !gfile.IsDir(apiModuleFolderPath) {
			continue
		}
		// generate go files by api module.
		var (
			module              = gfile.Basename(apiModuleFolderPath)
			dstModuleFolderPath = gfile.Join(in.DstFolder, module)
		)
		err = c.generateByModule(ctx, apiModuleFolderPath, dstModuleFolderPath)
		if err != nil {
			return nil, err
		}
	}

	mlog.Print(`done!`)
	return
}

// parseApiModule parses certain api and generate associated go files by certain module, not all api modules.
func (c CGenCtrl) generateByModule(ctx context.Context, apiModuleFolderPath, dstModuleFolderPath string) (err error) {
	// parse src and dst folder go files.
	apiItemsInSrc, err := c.getApiItemsInSrc(apiModuleFolderPath)
	if err != nil {
		return err
	}
	apiItemsInDst, err := c.getApiItemsInDst(dstModuleFolderPath)
	if err != nil {
		return err
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

	// generate api interface go files.
	if err = newApiInterfaceGenerator().Generate(apiModuleFolderPath, apiItemsInSrc); err != nil {
		return
	}

	// generate controller go files.
	if len(toBeImplementedApiItems) > 0 {
		err = newControllerGenerator().Generate(dstModuleFolderPath, toBeImplementedApiItems)
		if err != nil {
			return
		}
	}

	return
}
