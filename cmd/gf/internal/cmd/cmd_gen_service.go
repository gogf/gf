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
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type (
	cGenServiceInput struct {
		g.Meta       `name:"service" brief:"parse struct and associated functions from packages to generate service go file"`
		SrcFolder    string `short:"s" name:"srcFolder" brief:"source folder path to be parsed" d:"internal/logic"`
		DstFolder    string `short:"d" name:"dstFolder" brief:"destination folder path storing automatically generated go files" d:"internal/service"`
		StPattern    string `short:"a" name:"stPattern" brief:"regular expression matching struct name for generating service" d:"s(\\w+)"`
		ImportPrefix string `short:"p" name:"importPrefix" brief:"custom import prefix to calculate import path for generated go files"`
		WatchFile    string `short:"w" name:"watchFile" brief:"used in file watcher, it generates service go files only if given file is under Logic folder"`
		OverWrite    bool   `short:"o" name:"overwrite" brief:"overwrite files that already exist in generating folder" d:"true" orphan:"true"`
	}
	cGenServiceOutput struct{}
)

func (c cGen) Service(ctx context.Context, in cGenServiceInput) (out *cGenServiceOutput, err error) {
	in.SrcFolder = gstr.Trim(in.SrcFolder, `\/`)
	in.WatchFile = gstr.Trim(in.WatchFile, `\/`)
	if !gfile.Exists(in.SrcFolder) {
		mlog.Fatalf(`logic folder path "%s" does not exist`, in.SrcFolder)
	}
	if in.WatchFile != "" {
		// It works only if given WatchFile is in Logic folder.
		if !gstr.Contains(gstr.Replace(in.WatchFile, "\\", "/"), gstr.Replace(in.SrcFolder, "\\", "/")) {
			mlog.Printf(`ignore watch file "%s", not in source path "%s"`, in.WatchFile, in.SrcFolder)
			return
		}
	}
	if in.ImportPrefix == "" {
		if !gfile.Exists("go.mod") {
			mlog.Fatal("go.mod does not exist in current working directory")
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
		files          []string
		fileContent    string
		matches        [][]string
		srcPackages    []string
		dstPackageName = gstr.ToLower(gfile.Basename(in.DstFolder))
	)
	logicFolders, err := gfile.ScanDir(in.SrcFolder, "*", false)
	if err != nil {
		return nil, err
	}
	for _, logicFolder := range logicFolders {
		if !gfile.IsDir(logicFolder) {
			continue
		}
		if files, err = gfile.ScanDir(logicFolder, "*.go", false); err != nil {
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
			matches, err = gregex.MatchAllString(`func \(\w+ (.+?)\) (.+?) {`, fileContent)
			if err != nil {
				return nil, err
			}
			for _, match := range matches {
				var (
					structMatch []string
					structName  = gstr.Trim(match[1], "*")
					funcTitle   = gstr.Trim(gstr.Replace(match[2], "\n", ""))
				)
				if !gstr.IsLetterUpper(funcTitle[0]) {
					continue
				}
				structMatch, err = gregex.MatchString(in.StPattern, structName)
				if err != nil {
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
				interfaceFuncArray.Append(funcTitle)
			}
		}
		srcPackages = append(srcPackages, fmt.Sprintf(`%s/%s`, in.ImportPrefix, gfile.Basename(logicFolder)))
		// Generating go files for service.
		for structName, funcArray := range interfaceMap {
			var (
				filePath         = gfile.Join(in.DstFolder, gstr.ToLower(structName)+".go")
				generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceContent, g.MapStrStr{
					"{StructName}":     structName,
					"{PackageName}":    dstPackageName,
					"{FuncDefinition}": funcArray.Join("\n\t"),
				})
			)
			if !in.OverWrite && gfile.Exists(filePath) {
				mlog.Printf(`ignore generating service go file: %s`, filePath)
				continue
			}
			mlog.Printf(`generating service go file: %s`, filePath)
			if err = gfile.PutContents(filePath, generatedContent); err != nil {
				return nil, err
			}
		}
	}
	// Generate initialization go file.
	if len(srcPackages) > 0 {
		var (
			srcPackageName   = gstr.ToLower(gfile.Basename(in.SrcFolder))
			srcFilePath      = gfile.Join(in.SrcFolder, srcPackageName+".go")
			srcImports       string
			generatedContent string
		)
		for _, srcPackage := range srcPackages {
			srcImports += fmt.Sprintf(`%s_ "%s"%s`, "\t", srcPackage, "\n")
		}
		generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceLogicContent, g.MapStrStr{
			"{PackageName}": srcPackageName,
			"{Imports}":     srcImports,
		})
		mlog.Printf(`generating init go file: %s`, srcFilePath)
		if err = gfile.PutContents(srcFilePath, generatedContent); err != nil {
			return nil, err
		}
		utils.GoFmt(srcFilePath)
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
