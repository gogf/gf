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
		StPattern    string `short:"a" name:"stPattern" brief:"regular expression matching struct name for generating service. default: s([A-Z]\\\w+)" d:"s([A-Z]\\w+)"`
		Packages     string `short:"p" name:"packages" brief:"produce go files only for given source packages, multiple packages joined with char ','"`
		ImportPrefix string `short:"i" name:"importPrefix" brief:"custom import prefix to calculate import path for generated importing go file of logic"`
		OverWrite    bool   `short:"o" name:"overwrite" brief:"overwrite service go files that already exist in generating folder. default: true" d:"true" orphan:"true"`
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
		isDirty           bool
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
			srcPkgInterfaceMap       = make(map[string]*garray.StrArray)
			srcPkgInterfaceFuncArray *garray.StrArray
			ok                       bool
		)
		for _, file := range files {
			fileContent = gfile.GetContents(file)
			matches, err = gregex.MatchAllString(`func \((.+?)\) ([\s\S]+?) {`, fileContent)
			if err != nil {
				return nil, err
			}
			for _, match := range matches {
				var (
					structName    string
					structMatch   []string
					funcReceiver  = gstr.Trim(match[1])
					receiverArray = gstr.SplitAndTrim(funcReceiver, " ")
					functionHead  = gstr.Trim(gstr.Replace(match[2], "\n", ""))
				)
				if len(receiverArray) > 1 {
					structName = receiverArray[1]
				} else {
					structName = receiverArray[0]
				}
				structName = gstr.Trim(structName, "*")

				// Xxx(\n    ctx context.Context, req *v1.XxxReq,\n) -> Xxx(ctx context.Context, req *v1.XxxReq)
				functionHead = gstr.Replace(functionHead, `,)`, `)`)
				functionHead, _ = gregex.ReplaceString(`\(\s+`, `(`, functionHead)
				functionHead, _ = gregex.ReplaceString(`\s{2,}`, ` `, functionHead)
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
				if srcPkgInterfaceFuncArray, ok = srcPkgInterfaceMap[structName]; !ok {
					srcPkgInterfaceMap[structName] = garray.NewStrArray()
					srcPkgInterfaceFuncArray = srcPkgInterfaceMap[structName]
				}
				// Remove package name calls of `dstPackageName` in produced codes.
				functionHead, _ = gregex.ReplaceString(fmt.Sprintf(`\*{0,1}%s\.`, dstPackageName), ``, functionHead)
				srcPkgInterfaceFuncArray.Append(functionHead)
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
		if ok, err = c.generateServiceFiles(in, srcPkgInterfaceMap, dstPackageName); err != nil {
			return
		}
		if ok {
			isDirty = true
		}
	}

	if isDirty {
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
		if err = c.replaceGeneratedServiceContentGFV2(in); err != nil {
			return nil, err
		}
		mlog.Printf(`gofmt go files in "%s"`, in.DstFolder)
		utils.GoFmt(in.DstFolder)
	}

	mlog.Print(`done!`)
	return
}

func (c cGen) generateServiceFiles(
	in cGenServiceInput, srcPkgInterfaceMap map[string]*garray.StrArray, dstPackageName string,
) (ok bool, err error) {
	for structName, funcArray := range srcPkgInterfaceMap {
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
				mlog.Printf(`not overwrite, ignore generating service go file: %s`, filePath)
				continue
			}
			if !c.isToGenerateServiceGoFile(filePath, funcArray) {
				mlog.Printf(`not dirty, ignore generating service go file: %s`, filePath)
				continue
			}
		}
		ok = true
		mlog.Printf(`generating service go file: %s`, filePath)
		if err = gfile.PutContents(filePath, generatedContent); err != nil {
			return ok, err
		}
	}
	return ok, nil
}

// isToGenerateServiceGoFile checks and returns whether the service content dirty.
func (c cGen) isToGenerateServiceGoFile(filePath string, funcArray *garray.StrArray) bool {
	if !utils.IsFileDoNotEdit(filePath) {
		mlog.Debugf(`ignore file as it is manually maintained: %s`, filePath)
		return false
	}
	var (
		fileContent        = gfile.GetContents(filePath)
		generatedFuncArray = garray.NewSortedStrArrayFrom(funcArray.Slice())
		contentFuncArray   = garray.NewSortedStrArray()
	)
	if fileContent == "" {
		return true
	}
	match, _ := gregex.MatchString(`interface\s+{([\s\S]+?)}`, fileContent)
	if len(match) != 2 {
		return false
	}
	contentFuncArray.Append(gstr.SplitAndTrim(match[1], "\n")...)
	if generatedFuncArray.Len() != contentFuncArray.Len() {
		return true
	}
	for i := 0; i < generatedFuncArray.Len(); i++ {
		if generatedFuncArray.At(i) != contentFuncArray.At(i) {
			mlog.Debugf(`dirty, %s != %s`, generatedFuncArray.At(i), contentFuncArray.At(i))
			return true
		}
	}
	return false
}

func (c cGen) generateInitializationFile(in cGenServiceInput, importSrcPackages []string) (err error) {
	var (
		srcPackageName   = gstr.ToLower(gfile.Basename(in.SrcFolder))
		srcFilePath      = gfile.Join(in.SrcFolder, srcPackageName+".go")
		srcImports       string
		generatedContent string
	)
	if !utils.IsFileDoNotEdit(srcFilePath) {
		mlog.Debugf(`ignore file as it is manually maintained: %s`, srcFilePath)
		return nil
	}
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

func (c cGen) replaceGeneratedServiceContentGFV2(in cGenServiceInput) (err error) {
	return gfile.ReplaceDirFunc(func(path, content string) string {
		if gstr.Contains(content, `"github.com/gogf/gf`) && !gstr.Contains(content, `"github.com/gogf/gf/v2`) {
			content = gstr.Replace(content, `"github.com/gogf/gf"`, `"github.com/gogf/gf/v2"`)
			content = gstr.Replace(content, `"github.com/gogf/gf/`, `"github.com/gogf/gf/v2/`)
			return content
		}
		return content
	}, in.DstFolder, "*.go", false)
}
