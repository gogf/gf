// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genctrl

import (
	"context"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
)

const (
	CGenCtrlConfig = `gfcli.gen.ctrl`
	CGenCtrlUsage  = `gf gen ctrl [OPTION]`
	CGenCtrlBrief  = `parse api definitions to generate controller/sdk go files`
	CGenCtrlEg     = `
gf gen ctrl
`
	CGenCtrlBriefSrcFolder     = `source folder path to be parsed. default: api`
	CGenCtrlBriefDstFolder     = `destination folder path storing automatically generated go files. default: internal/controller`
	CGenCtrlBriefWatchFile     = `used in file watcher, it re-generates go files only if given file is under srcFolder`
	CGenCtrlBriefSdkPath       = `also generate SDK go files for api definitions to specified directory`
	CGenCtrlBriefSdkStdVersion = `use standard version prefix for generated sdk request path`
	CGenCtrlBriefSdkNoV1       = `do not add version suffix for interface module name if version is v1`
	CGenCtrlBriefClear         = `auto delete generated and unimplemented controller go files if api definitions are missing`
	CGenCtrlControllerMerge    = `generate all controller files into one go file by name of api definition source go file`
)

const (
	PatternCtrlDefinition = `func\s+\(.+?\)\s+\w+\(.+?\*(\w+)\.(\w+)Req\)\s+\(.+?\*(\w+)\.(\w+)Res,\s+\w+\s+error\)\s+{`
)

const (
	genCtrlFileLockSeconds = 10
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenCtrlConfig`:             CGenCtrlConfig,
		`CGenCtrlUsage`:              CGenCtrlUsage,
		`CGenCtrlBrief`:              CGenCtrlBrief,
		`CGenCtrlEg`:                 CGenCtrlEg,
		`CGenCtrlBriefSrcFolder`:     CGenCtrlBriefSrcFolder,
		`CGenCtrlBriefDstFolder`:     CGenCtrlBriefDstFolder,
		`CGenCtrlBriefWatchFile`:     CGenCtrlBriefWatchFile,
		`CGenCtrlBriefSdkPath`:       CGenCtrlBriefSdkPath,
		`CGenCtrlBriefSdkStdVersion`: CGenCtrlBriefSdkStdVersion,
		`CGenCtrlBriefSdkNoV1`:       CGenCtrlBriefSdkNoV1,
		`CGenCtrlBriefClear`:         CGenCtrlBriefClear,
		`CGenCtrlControllerMerge`:    CGenCtrlControllerMerge,
	})
}

type (
	CGenCtrl      struct{}
	CGenCtrlInput struct {
		g.Meta        `name:"ctrl" config:"{CGenCtrlConfig}" usage:"{CGenCtrlUsage}" brief:"{CGenCtrlBrief}" eg:"{CGenCtrlEg}"`
		SrcFolder     string `short:"s" name:"srcFolder"     brief:"{CGenCtrlBriefSrcFolder}" d:"api"`
		DstFolder     string `short:"d" name:"dstFolder"     brief:"{CGenCtrlBriefDstFolder}" d:"internal/controller"`
		WatchFile     string `short:"w" name:"watchFile"     brief:"{CGenCtrlBriefWatchFile}"`
		SdkPath       string `short:"k" name:"sdkPath"       brief:"{CGenCtrlBriefSdkPath}"`
		SdkStdVersion bool   `short:"v" name:"sdkStdVersion" brief:"{CGenCtrlBriefSdkStdVersion}" orphan:"true"`
		SdkNoV1       bool   `short:"n" name:"sdkNoV1"       brief:"{CGenCtrlBriefSdkNoV1}" orphan:"true"`
		Clear         bool   `short:"c" name:"clear"         brief:"{CGenCtrlBriefClear}" orphan:"true"`
		Merge         bool   `short:"m" name:"merge"         brief:"{CGenCtrlControllerMerge}" orphan:"true"`
	}
	CGenCtrlOutput struct{}
)

func (c CGenCtrl) Ctrl(ctx context.Context, in CGenCtrlInput) (out *CGenCtrlOutput, err error) {
	if in.WatchFile != "" {
		err = c.generateByWatchFile(
			in.WatchFile, in.SdkPath, in.SdkStdVersion, in.SdkNoV1, in.Clear, in.Merge,
		)
		mlog.Print(`done!`)
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
		err = c.generateByModule(
			apiModuleFolderPath, dstModuleFolderPath, in.SdkPath,
			in.SdkStdVersion, in.SdkNoV1, in.Clear, in.Merge,
		)
		if err != nil {
			return nil, err
		}
	}

	mlog.Print(`done!`)
	return
}

func (c CGenCtrl) generateByWatchFile(watchFile, sdkPath string, sdkStdVersion, sdkNoV1, clear, merge bool) (err error) {
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
	if gfile.Exists(watchFile) {
		structsInfo, err := c.getStructsNameInSrc(watchFile)
		if err != nil {
			return err
		}
		if len(structsInfo) == 0 {
			return nil
		}
	}

	var (
		projectRootPath     = gfile.Dir(gfile.Dir(apiModuleFolderPath))
		module              = gfile.Basename(apiModuleFolderPath)
		dstModuleFolderPath = gfile.Join(projectRootPath, "internal", "controller", module)
	)
	return c.generateByModule(
		apiModuleFolderPath, dstModuleFolderPath, sdkPath, sdkStdVersion, sdkNoV1, clear, merge,
	)
}

// parseApiModule parses certain api and generate associated go files by certain module, not all api modules.
func (c CGenCtrl) generateByModule(
	apiModuleFolderPath, dstModuleFolderPath, sdkPath string,
	sdkStdVersion, sdkNoV1, clear, merge bool,
) (err error) {
	// parse src and dst folder go files.
	apiItemsInSrc, err := c.getApiItemsInSrc(apiModuleFolderPath)
	if err != nil {
		return err
	}
	apiItemsInDst, err := c.getApiItemsInDst(dstModuleFolderPath)
	if err != nil {
		return err
	}

	// generate api interface go files.
	if err = newApiInterfaceGenerator().Generate(apiModuleFolderPath, apiItemsInSrc); err != nil {
		return
	}

	// generate controller go files.
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
	if len(toBeImplementedApiItems) > 0 {
		err = newControllerGenerator().Generate(dstModuleFolderPath, toBeImplementedApiItems, merge)
		if err != nil {
			return
		}
	}

	// delete unimplemented controllers if api definitions are missing.
	if clear {
		var (
			apiDefinitionSet    = gset.NewStrSet()
			extraApiItemsInCtrl = make([]apiItem, 0)
		)
		for _, item := range apiItemsInSrc {
			apiDefinitionSet.Add(item.String())
		}
		for _, item := range apiItemsInDst {
			if apiDefinitionSet.Contains(item.String()) {
				continue
			}
			extraApiItemsInCtrl = append(extraApiItemsInCtrl, item)
		}
		if len(extraApiItemsInCtrl) > 0 {
			err = newControllerClearer().Clear(dstModuleFolderPath, extraApiItemsInCtrl)
			if err != nil {
				return
			}
		}
	}

	// generate sdk go files.
	if sdkPath != "" {
		if err = newApiSdkGenerator().Generate(apiItemsInSrc, sdkPath, sdkStdVersion, sdkNoV1); err != nil {
			return
		}
	}
	return
}
