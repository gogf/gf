package cmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

type (
	cGenServiceInput struct {
		g.Meta       `name:"service" config:"gfcli.gen.service" brief:"parse struct and associated functions from packages to generate service go file"`
		SrcFolder    string `short:"s" name:"srcFolder" brief:"source folder path to be parsed. default: internal/logic" d:"internal/logic"`
		DstFolder    string `short:"d" name:"dstFolder" brief:"destination folder path storing automatically generated go files. default: internal/service" d:"internal/service"`
		WatchFile    string `short:"w" name:"watchFile" brief:"used in file watcher, it generates service go files only if given file is under srcFolder"`
		StPattern    string `short:"a" name:"stPattern" brief:"regular expression matching struct name for generating service. default: s(\\w+)" d:"s(\\w+)"`
		Packages     string `short:"p" name:"packages" brief:"produce go files only for given source packages, multiple packages joined with char ','"`
		ImportPrefix string `short:"i" name:"importPrefix" brief:"custom import prefix to calculate import path for generated go files"`
		OverWrite    bool   `short:"o" name:"overwrite" brief:"overwrite files that already exist in generating folder. default: true" d:"true" orphan:"true"`
	}
	cGenServiceOutput struct{}
)

const (
	genServiceFileLockSeconds = 10
)

func (c cGen) Service(ctx context.Context, in cGenServiceInput) (out *cGenServiceOutput, err error) {
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
		if !gstr.HasSuffix(srcFolderDir, in.SrcFolder) {
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
		err = gproc.ShellRun(command)
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
		files             []string
		fileContent       string
		matches           [][]string
		importSrcPackages []string
		inputPackages     = gstr.SplitAndTrim(in.Packages, ",")
		dstPackageName    = gstr.ToLower(gfile.Basename(in.DstFolder))
	)
	srcFolders, err := gfile.ScanDir(in.SrcFolder, "*", false)
	if err != nil {
		return nil, err
	}
	for _, srcFolder := range srcFolders {
		if !gfile.IsDir(srcFolder) {
			continue
		}
		if files, err = gfile.ScanDir(srcFolder, "*.go", false); err != nil {
			return nil, err
		}
		if len(files) == 0 {
			continue
		}
		var (
			// StructName => FunctionDefinitions
			interfaceMap       = make(map[string]*garray.StrArray)
			interfaceFuncArray *garray.StrArray
			ok                 bool
		)
		for _, file := range files {
			fileContent = gfile.GetContents(file)
			matches, err = gregex.MatchAllString(`func \(\w+ (.+?)\) ([\s\S]+?) {`, fileContent)
			if err != nil {
				return nil, err
			}
			for _, match := range matches {
				var (
					structMatch  []string
					structName   = gstr.Trim(match[1], "*")
					functionHead = gstr.Trim(gstr.Replace(match[2], "\n", ""))
				)
				if !gstr.IsLetterUpper(functionHead[0]) {
					continue
				}
				if structMatch, err = gregex.MatchString(in.StPattern, structName); err != nil {
					return nil, err
				}
				if len(structMatch) < 1 {
					continue
				}
				structName = gstr.CaseCamel(structMatch[1])
				if interfaceFuncArray, ok = interfaceMap[structName]; !ok {
					interfaceMap[structName] = garray.NewStrArray()
					interfaceFuncArray = interfaceMap[structName]
				}
				// Remove package name calls of `dstPackageName` in produced codes.
				functionHead, _ = gregex.ReplaceString(fmt.Sprintf(`\*{0,1}%s\.`, dstPackageName), ``, functionHead)
				interfaceFuncArray.Append(functionHead)
			}
		}
		importSrcPackages = append(
			importSrcPackages,
			fmt.Sprintf(`%s/%s`, in.ImportPrefix, gfile.Basename(srcFolder)),
		)
		// Ignore source packages if input packages given.
		if len(inputPackages) > 0 && !gstr.InArray(inputPackages, gfile.Basename(srcFolder)) {
			mlog.Debugf(
				`ignore source package "%s" as it is not in desired packages: %+v`,
				gfile.Basename(srcFolder), inputPackages,
			)
			continue
		}
		// Generating go files for service.
		if err = c.generateServiceFiles(in, interfaceMap, dstPackageName); err != nil {
			return
		}
	}
	// Generate initialization go file.
	if len(importSrcPackages) > 0 {
		if err = c.generateInitializationFile(in, importSrcPackages); err != nil {
			return
		}
	}

	// Go imports updating.
	mlog.Printf(`goimports go files in "%s", it may take seconds...`, in.DstFolder)
	utils.GoImports(in.DstFolder)

	// Replica v1 to v2 for GoFrame.
	err = gfile.ReplaceDirFunc(func(path, content string) string {
		if gstr.Contains(content, `"github.com/gogf/gf`) && !gstr.Contains(content, `"github.com/gogf/gf/v2`) {
			content = gstr.Replace(content, `"github.com/gogf/gf"`, `"github.com/gogf/gf/v2"`)
			content = gstr.Replace(content, `"github.com/gogf/gf/`, `"github.com/gogf/gf/v2/`)
			return content
		}
		return content
	}, in.DstFolder, "*.go", false)
	if err != nil {
		return nil, err
	}
	mlog.Printf(`gofmt go files in "%s"`, in.DstFolder)
	utils.GoFmt(in.DstFolder)
	mlog.Print(`done!`)
	return
}

func (c cGen) generateServiceFiles(
	in cGenServiceInput, interfaceMap map[string]*garray.StrArray, dstPackageName string,
) (err error) {
	for structName, funcArray := range interfaceMap {
		var (
			filePath         = gfile.Join(in.DstFolder, gstr.ToLower(structName)+".go")
			generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceContent, g.MapStrStr{
				"{StructName}":     structName,
				"{PackageName}":    dstPackageName,
				"{FuncDefinition}": funcArray.Join("\n\t"),
			})
		)
		if gfile.Exists(filePath) {
			if !in.OverWrite {
				mlog.Printf(`ignore generating service go file: %s`, filePath)
				continue
			}
		}

		mlog.Printf(`generating service go file: %s`, filePath)
		if err = gfile.PutContents(filePath, generatedContent); err != nil {
			return err
		}
	}
	return nil
}

func (c cGen) generateInitializationFile(in cGenServiceInput, importSrcPackages []string) (err error) {
	var (
		srcPackageName   = gstr.ToLower(gfile.Basename(in.SrcFolder))
		srcFilePath      = gfile.Join(in.SrcFolder, srcPackageName+".go")
		srcImports       string
		generatedContent string
	)
	for _, importSrcPackage := range importSrcPackages {
		srcImports += fmt.Sprintf(`%s_ "%s"%s`, "\t", importSrcPackage, "\n")
	}
	generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceLogicContent, g.MapStrStr{
		"{PackageName}": srcPackageName,
		"{Imports}":     srcImports,
	})
	mlog.Printf(`generating init go file: %s`, srcFilePath)
	if err = gfile.PutContents(srcFilePath, generatedContent); err != nil {
		return err
	}
	utils.GoFmt(srcFilePath)
	return nil
}
