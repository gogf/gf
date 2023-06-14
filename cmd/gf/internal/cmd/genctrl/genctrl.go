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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
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
	CGenCtrlBriefSdkPath   = `also generate SDK go files to specified directory`
)

const (
	PatternApiDefinition  = `type\s+(\w+)Req\s+struct\s+{`
	PatternCtrlDefinition = `func\s+\(.+?\)\s+\w+\(.+?\*(\w+)\.(\w+)Req\)\s+\(.+?\*(\w+)\.(\w+)Res,\s+\w+\s+error\)\s+{`
)

const (
	genCtrlFileLockSeconds = 10
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
		`CGenCtrlBriefSdkPath`:   CGenCtrlBriefSdkPath,
	})
}

type (
	CGenCtrl      struct{}
	CGenCtrlInput struct {
		g.Meta    `name:"ctrl" config:"{CGenCtrlConfig}" usage:"{CGenCtrlUsage}" brief:"{CGenCtrlBrief}" eg:"{CGenCtrlEg}"`
		SrcFolder string `short:"s" name:"srcFolder" brief:"{CGenCtrlBriefSrcFolder}" d:"api"`
		DstFolder string `short:"d" name:"dstFolder" brief:"{CGenCtrlBriefDstFolder}" d:"internal/controller"`
		WatchFile string `short:"w" name:"watchFile" brief:"{CGenCtrlBriefWatchFile}"`
		SdkPath   string `short:"k" name:"sdkPath"   brief:"{CGenCtrlBriefSdkPath}"`
	}
	CGenCtrlOutput struct{}
)

func (c CGenCtrl) Ctrl(ctx context.Context, in CGenCtrlInput) (out *CGenCtrlOutput, err error) {
	if in.WatchFile != "" {
		err = c.generateByWatchFile(in.WatchFile, in.SdkPath)
		return
	}

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
		err = c.generateByModule(apiModuleFolderPath, dstModuleFolderPath, in.SdkPath)
		if err != nil {
			return nil, err
		}
	}

	mlog.Print(`done!`)
	return
}

func (c CGenCtrl) generateByWatchFile(watchFile, sdkPath string) (err error) {
	// File lock to avoid multiple processes.
	var (
		flockFilePath = gfile.Temp("gf.cli.gen.service.lock")
		flockContent  = gfile.GetContents(flockFilePath)
	)
	if flockContent != "" {
		if gtime.Timestamp()-gconv.Int64(flockContent) < genCtrlFileLockSeconds {
			// If another generating process is running, it just exits.
			mlog.Debug(`another "gen service" process is running, exit`)
			return
		}
	}
	defer gfile.Remove(flockFilePath)
	_ = gfile.PutContents(flockFilePath, gtime.TimestampStr())

	// check this updated file is an api file.
	// watch file should be in standard goframe project structure.
	var (
		apiVersionPath      = gfile.Dir(watchFile)
		apiModuleFolderPath = gfile.Dir(apiVersionPath)
		shouldBeNameOfAPi   = gfile.Basename(gfile.Dir(apiModuleFolderPath))
	)
	if shouldBeNameOfAPi != "api" {
		return nil
	}
	// watch file should have api definitions.
	if !gregex.IsMatchString(PatternApiDefinition, gfile.GetContents(watchFile)) {
		return nil
	}
	var (
		projectRootPath     = gfile.Dir(gfile.Dir(apiModuleFolderPath))
		module              = gfile.Basename(apiModuleFolderPath)
		dstModuleFolderPath = gfile.Join(projectRootPath, "internal", "controller", module)
	)
	return c.generateByModule(apiModuleFolderPath, dstModuleFolderPath, sdkPath)
}

// parseApiModule parses certain api and generate associated go files by certain module, not all api modules.
func (c CGenCtrl) generateByModule(apiModuleFolderPath, dstModuleFolderPath, sdkPath string) (err error) {
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

	// generate sdk go files.
	if sdkPath != "" {
		if err = newApiSdkGenerator().Generate(sdkPath, apiItemsInSrc); err != nil {
			return
		}
	}

	return
}
