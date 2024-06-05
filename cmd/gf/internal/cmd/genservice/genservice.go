// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
)

const (
	CGenServiceConfig = `gfcli.gen.service`
	CGenServiceUsage  = `gf gen service [OPTION]`
	CGenServiceBrief  = `parse struct and associated functions from packages to generate service go file`
	CGenServiceEg     = `
gf gen service
gf gen service -f Snake
`
	CGenServiceBriefSrcFolder    = `source folder path to be parsed. default: internal/logic`
	CGenServiceBriefDstFolder    = `destination folder path storing automatically generated go files. default: internal/service`
	CGenServiceBriefFileNameCase = `
destination file name storing automatically generated go files, cases are as follows:
| Case            | Example            |
|---------------- |--------------------|
| Lower           | anykindofstring    |
| Camel           | AnyKindOfString    |
| CamelLower      | anyKindOfString    |
| Snake           | any_kind_of_string | default
| SnakeScreaming  | ANY_KIND_OF_STRING |
| SnakeFirstUpper | rgb_code_md5       |
| Kebab           | any-kind-of-string |
| KebabScreaming  | ANY-KIND-OF-STRING |
`
	CGenServiceBriefWatchFile    = `used in file watcher, it re-generates all service go files only if given file is under srcFolder`
	CGenServiceBriefStPattern    = `regular expression matching struct name for generating service. default: ^s([A-Z]\\\\w+)$`
	CGenServiceBriefPackages     = `produce go files only for given source packages(source folders)`
	CGenServiceBriefImportPrefix = `custom import prefix to calculate import path for generated importing go file of logic`
	CGenServiceBriefClear        = `delete all generated go files that are not used any further`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenServiceConfig`:            CGenServiceConfig,
		`CGenServiceUsage`:             CGenServiceUsage,
		`CGenServiceBrief`:             CGenServiceBrief,
		`CGenServiceEg`:                CGenServiceEg,
		`CGenServiceBriefSrcFolder`:    CGenServiceBriefSrcFolder,
		`CGenServiceBriefDstFolder`:    CGenServiceBriefDstFolder,
		`CGenServiceBriefFileNameCase`: CGenServiceBriefFileNameCase,
		`CGenServiceBriefWatchFile`:    CGenServiceBriefWatchFile,
		`CGenServiceBriefStPattern`:    CGenServiceBriefStPattern,
		`CGenServiceBriefPackages`:     CGenServiceBriefPackages,
		`CGenServiceBriefImportPrefix`: CGenServiceBriefImportPrefix,
		`CGenServiceBriefClear`:        CGenServiceBriefClear,
	})
}

type (
	CGenService      struct{}
	CGenServiceInput struct {
		g.Meta          `name:"service" config:"{CGenServiceConfig}" usage:"{CGenServiceUsage}" brief:"{CGenServiceBrief}" eg:"{CGenServiceEg}"`
		SrcFolder       string   `short:"s" name:"srcFolder" brief:"{CGenServiceBriefSrcFolder}" d:"internal/logic"`
		DstFolder       string   `short:"d" name:"dstFolder" brief:"{CGenServiceBriefDstFolder}" d:"internal/service"`
		DstFileNameCase string   `short:"f" name:"dstFileNameCase" brief:"{CGenServiceBriefFileNameCase}" d:"Snake"`
		WatchFile       string   `short:"w" name:"watchFile" brief:"{CGenServiceBriefWatchFile}"`
		StPattern       string   `short:"a" name:"stPattern" brief:"{CGenServiceBriefStPattern}" d:"^s([A-Z]\\w+)$"`
		Packages        []string `short:"p" name:"packages" brief:"{CGenServiceBriefPackages}"`
		ImportPrefix    string   `short:"i" name:"importPrefix" brief:"{CGenServiceBriefImportPrefix}"`
		Clear           bool     `short:"l" name:"clear" brief:"{CGenServiceBriefClear}" orphan:"true"`
	}
	CGenServiceOutput struct{}
)

const (
	genServiceFileLockSeconds = 10
)

func (c CGenService) Service(ctx context.Context, in CGenServiceInput) (out *CGenServiceOutput, err error) {
	in.SrcFolder = filepath.ToSlash(in.SrcFolder)
	in.SrcFolder = gstr.TrimRight(in.SrcFolder, `/`)
	in.WatchFile = filepath.ToSlash(in.WatchFile)
	in.WatchFile = gstr.TrimRight(in.WatchFile, `/`)

	// Watch file handling.
	if in.WatchFile != "" {
		// File lock to avoid multiple processes.
		var (
			flockFilePath = gfile.Temp("gf.cli.gen.service.lock")
			flockContent  = gfile.GetContents(flockFilePath)
		)
		if flockContent != "" {
			if gtime.Timestamp()-gconv.Int64(flockContent) < genServiceFileLockSeconds {
				// If another "gen service" process is running, it just exits.
				mlog.Debug(`another "gen service" process is running, exit`)
				return
			}
		}
		defer gfile.Remove(flockFilePath)
		_ = gfile.PutContents(flockFilePath, gtime.TimestampStr())

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

		in.WatchFile = ""
		in.Packages = []string{gfile.Basename(watchFileDir)}
		return c.Service(ctx, in)
	}

	if !gfile.Exists(in.SrcFolder) {
		mlog.Fatalf(`source folder path "%s" does not exist`, in.SrcFolder)
	}

	if in.ImportPrefix == "" {
		in.ImportPrefix = utils.GetImportPath(in.SrcFolder)
	}

	var (
		isDirty                 atomic.Value                                 // Temp boolean.
		files                   []string                                     // Temp file array.
		initImportSrcPackages   []string                                     // Used for generating logic.go.
		inputPackages           = in.Packages                                // Custom packages.
		dstPackageName          = gstr.ToLower(gfile.Basename(in.DstFolder)) // Package name for generated go files.
		generatedDstFilePathSet = gset.NewStrSet()                           // All generated file path set.
	)
	isDirty.Store(false)

	// The first level folders.
	srcFolderPaths, err := gfile.ScanDir(in.SrcFolder, "*", false)
	if err != nil {
		return nil, err
	}
	// it will use goroutine to generate service files for each package.
	var wg = sync.WaitGroup{}
	for _, srcFolderPath := range srcFolderPaths {
		if !gfile.IsDir(srcFolderPath) {
			continue
		}
		// Only retrieve sub files, no recursively.
		if files, err = gfile.ScanDir(srcFolderPath, "*.go", false); err != nil {
			return nil, err
		}
		if len(files) == 0 {
			continue
		}
		// Parse single logic package folder.
		var (
			srcPackageName      = gfile.Basename(srcFolderPath)
			srcImportedPackages = garray.NewSortedStrArray().SetUnique(true)
			srcStructFunctions  = gmap.NewListMap()
			dstFilePath         = gfile.Join(in.DstFolder,
				c.getDstFileNameCase(srcPackageName, in.DstFileNameCase)+".go",
			)
		)
		generatedDstFilePathSet.Add(dstFilePath)
		// if it were to use goroutine,
		// it would cause the order of the generated functions in the file to be disordered.
		for _, file := range files {
			pkgItems, funcItems, err := c.parseItemsInSrc(file)
			if err != nil {
				return nil, err
			}

			// Calculate imported packages for service generating.
			err = c.calculateImportedItems(in, pkgItems, funcItems, srcImportedPackages)
			if err != nil {
				return nil, err
			}

			// Calculate functions and interfaces for service generating.
			err = c.calculateFuncItems(in, funcItems, srcStructFunctions)
			if err != nil {
				return nil, err
			}
		}

		initImportSrcPackages = append(
			initImportSrcPackages,
			fmt.Sprintf(`%s/%s`, in.ImportPrefix, srcPackageName),
		)
		// Ignore source packages if input packages given.
		if len(inputPackages) > 0 && !gstr.InArray(inputPackages, srcPackageName) {
			mlog.Debugf(
				`ignore source package "%s" as it is not in desired packages: %+v`,
				srcPackageName, inputPackages,
			)
			continue
		}

		// Generating service go file for single logic package.
		wg.Add(1)
		go func(generateServiceFilesInput generateServiceFilesInput) {
			defer wg.Done()
			ok, err := c.generateServiceFile(generateServiceFilesInput)
			if err != nil {
				mlog.Printf(`error generating service file "%s": %v`, generateServiceFilesInput.DstFilePath, err)
			}
			if !isDirty.Load().(bool) && ok {
				isDirty.Store(true)
			}
		}(generateServiceFilesInput{
			CGenServiceInput:    in,
			SrcPackageName:      srcPackageName,
			SrcImportedPackages: srcImportedPackages.Slice(),
			SrcStructFunctions:  srcStructFunctions,
			DstPackageName:      dstPackageName,
			DstFilePath:         dstFilePath,
		})
	}
	wg.Wait()

	if in.Clear {
		files, err = gfile.ScanDirFile(in.DstFolder, "*.go", false)
		if err != nil {
			return nil, err
		}
		var relativeFilePath string
		for _, file := range files {
			relativeFilePath = gstr.SubStrFromR(file, in.DstFolder)
			if !generatedDstFilePathSet.Contains(relativeFilePath) &&
				utils.IsFileDoNotEdit(relativeFilePath) {

				mlog.Printf(`remove no longer used service file: %s`, relativeFilePath)
				if err = gfile.Remove(file); err != nil {
					return nil, err
				}
			}
		}
	}

	if isDirty.Load().(bool) {
		// Generate initialization go file.
		if len(initImportSrcPackages) > 0 {
			if err = c.generateInitializationFile(in, initImportSrcPackages); err != nil {
				return
			}
		}

		// Replace v1 to v2 for GoFrame.
		if err = utils.ReplaceGeneratedContentGFV2(in.DstFolder); err != nil {
			return nil, err
		}
		mlog.Printf(`gofmt go files in "%s"`, in.DstFolder)
		utils.GoFmt(in.DstFolder)
	}

	// auto update main.go.
	if err = c.checkAndUpdateMain(in.SrcFolder); err != nil {
		return nil, err
	}

	mlog.Print(`done!`)
	return
}

func (c CGenService) checkAndUpdateMain(srcFolder string) (err error) {
	var (
		logicPackageName = gstr.ToLower(gfile.Basename(srcFolder))
		logicFilePath    = gfile.Join(srcFolder, logicPackageName+".go")
		importPath       = utils.GetImportPath(logicFilePath)
		importStr        = fmt.Sprintf(`_ "%s"`, importPath)
		mainFilePath     = gfile.Join(gfile.Dir(gfile.Dir(gfile.Dir(logicFilePath))), "main.go")
		mainFileContent  = gfile.GetContents(mainFilePath)
	)
	// No main content found.
	if mainFileContent == "" {
		return nil
	}
	if gstr.Contains(mainFileContent, importStr) {
		return nil
	}
	match, err := gregex.MatchString(`import \(([\s\S]+?)\)`, mainFileContent)
	if err != nil {
		return err
	}
	// No match.
	if len(match) < 2 {
		return nil
	}
	lines := garray.NewStrArrayFrom(gstr.Split(match[1], "\n"))
	for i, line := range lines.Slice() {
		line = gstr.Trim(line)
		if len(line) == 0 {
			continue
		}
		if line[0] == '_' {
			continue
		}
		// Insert the logic import into imports.
		if err = lines.InsertBefore(i, fmt.Sprintf("\t%s\n\n", importStr)); err != nil {
			return err
		}
		break
	}
	mainFileContent, err = gregex.ReplaceString(
		`import \(([\s\S]+?)\)`,
		fmt.Sprintf(`import (%s)`, lines.Join("\n")),
		mainFileContent,
	)
	if err != nil {
		return err
	}
	mlog.Print(`update main.go`)
	err = gfile.PutContents(mainFilePath, mainFileContent)
	utils.GoFmt(mainFilePath)
	return
}
