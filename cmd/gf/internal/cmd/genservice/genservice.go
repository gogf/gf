package genservice

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
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
	CGenServiceBriefPackages     = `produce go files only for given source packages`
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
			`%s gen service -packages=%s`,
			gfile.SelfName(), gfile.Basename(watchFileDir),
		)
		err = gproc.ShellRun(ctx, command)
		return
	}

	if !gfile.Exists(in.SrcFolder) {
		mlog.Fatalf(`source folder path "%s" does not exist`, in.SrcFolder)
	}

	if in.ImportPrefix == "" {
		if !gfile.Exists("go.mod") {
			mlog.Fatal("ImportPrefix is empty and go.mod does not exist in current working directory")
		}
		var (
			goModContent = gfile.GetContents("go.mod")
			match, _     = gregex.MatchString(`^module\s+(.+)\s*`, goModContent)
		)
		if len(match) > 1 {
			in.ImportPrefix = fmt.Sprintf(`%s/%s`, gstr.Trim(match[1]), gstr.Replace(in.SrcFolder, `\`, `/`))
		}
	}

	var (
		isDirty                 bool                                         // Temp boolean.
		files                   []string                                     // Temp file array.
		fileContent             string                                       // Temp file content for handling go file.
		initImportSrcPackages   []string                                     // Used for generating logic.go.
		inputPackages           = in.Packages                                // Custom packages.
		dstPackageName          = gstr.ToLower(gfile.Basename(in.DstFolder)) // Package name for generated go files.
		generatedDstFilePathSet = gset.NewStrSet()                           // All generated file path set.
	)
	// The first level folders.
	srcFolderPaths, err := gfile.ScanDir(in.SrcFolder, "*", false)
	if err != nil {
		return nil, err
	}
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
		var (
			// StructName => FunctionDefinitions
			srcPkgInterfaceMap  = make(map[string]*garray.StrArray)
			srcImportedPackages = garray.NewSortedStrArray().SetUnique(true)
			srcPackageName      = gfile.Basename(srcFolderPath)
			ok                  bool
			dstFilePath         = gfile.Join(in.DstFolder,
				c.getDstFileNameCase(srcPackageName, in.DstFileNameCase)+".go",
			)
		)
		generatedDstFilePathSet.Add(dstFilePath)
		for _, file := range files {
			fileContent = gfile.GetContents(file)
			fileContent, err := gregex.ReplaceString(`/[/|\*](.+)`, "", fileContent)
			if err != nil {
				return nil, err
			}
			// Calculate imported packages of source go files.
			err = c.calculateImportedPackages(fileContent, srcImportedPackages)
			if err != nil {
				return nil, err
			}
			// Calculate functions and interfaces for service generating.
			err = c.calculateInterfaceFunctions(in, fileContent, srcPkgInterfaceMap, dstPackageName)
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
		// Generating service go file for logic.
		if ok, err = c.generateServiceFile(generateServiceFilesInput{
			CGenServiceInput:    in,
			SrcStructFunctions:  srcPkgInterfaceMap,
			SrcImportedPackages: srcImportedPackages.Slice(),
			SrcPackageName:      srcPackageName,
			DstPackageName:      dstPackageName,
			DstFilePath:         dstFilePath,
		}); err != nil {
			return
		}
		if ok {
			isDirty = true
		}
	}

	if in.Clear {
		files, err = gfile.ScanDirFile(in.DstFolder, "*.go", false)
		if err != nil {
			return nil, err
		}
		var relativeFilePath string
		for _, file := range files {
			relativeFilePath = gstr.SubStrFromR(file, in.DstFolder)
			if !generatedDstFilePathSet.Contains(relativeFilePath) && utils.IsFileDoNotEdit(relativeFilePath) {
				mlog.Printf(`remove no longer used service file: %s`, relativeFilePath)
				if err = gfile.Remove(file); err != nil {
					return nil, err
				}
			}
		}
	}

	if isDirty {
		// Generate initialization go file.
		if len(initImportSrcPackages) > 0 {
			if err = c.generateInitializationFile(in, initImportSrcPackages); err != nil {
				return
			}
		}

		// Replace v1 to v2 for GoFrame.
		if err = c.replaceGeneratedServiceContentGFV2(in); err != nil {
			return nil, err
		}
		mlog.Printf(`gofmt go files in "%s"`, in.DstFolder)
		utils.GoFmt(in.DstFolder)
	}

	mlog.Print(`done!`)
	return
}

func (c CGenService) replaceGeneratedServiceContentGFV2(in CGenServiceInput) (err error) {
	return gfile.ReplaceDirFunc(func(path, content string) string {
		if gstr.Contains(content, `"github.com/gogf/gf`) && !gstr.Contains(content, `"github.com/gogf/gf/v2`) {
			content = gstr.Replace(content, `"github.com/gogf/gf"`, `"github.com/gogf/gf/v2"`)
			content = gstr.Replace(content, `"github.com/gogf/gf/`, `"github.com/gogf/gf/v2/`)
			return content
		}
		return content
	}, in.DstFolder, "*.go", false)
}
