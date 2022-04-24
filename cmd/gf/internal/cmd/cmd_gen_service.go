package cmd

import (
	"context"

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
		g.Meta    `name:"service" brief:"parse logic struct and associated functions to generate service go file"`
		Logic     string `short:"l" name:"logic" brief:"logic folder path to be parsed" d:"internal/logic"`
		Path      string `short:"p" name:"path" brief:"folder path storing automatically generated go files" d:"internal/service"`
		Pattern   string `short:"a" name:"pattern" brief:"regular expression matching struct name for generating service" d:"s(\\w+)"`
		WatchFile string `short:"w" name:"watchFile" brief:"used in file watcher, it generates service go files only if given file is under Logic folder"`
	}
	cGenServiceOutput struct{}
)

func (c cGen) Service(ctx context.Context, in cGenServiceInput) (out *cGenPbOutput, err error) {
	in.Logic = gstr.Trim(in.Logic, `\/`)
	in.WatchFile = gstr.Trim(in.WatchFile, `\/`)
	if !gfile.Exists(in.Logic) {
		mlog.Fatalf(`logic folder path "%s" does not exist`, in.Logic)
	}
	if in.WatchFile != "" {
		// It works only if given WatchFile is in Logic folder.
		if !gstr.Contains(gstr.Replace(in.WatchFile, "\\", "/"), gstr.Replace(in.Logic, "\\", "/")) {
			mlog.Printf(`ignore watch file "%s", not in logic path "%s"`, in.WatchFile, in.Logic)
			return
		}
	}
	var (
		files       []string
		fileContent string
		matches     [][]string
		packageName = gstr.ToLower(gfile.Basename(in.Path))
	)
	logicFolders, err := gfile.ScanDir(in.Logic, "*", false)
	if err != nil {
		return nil, err
	}
	for _, logicFolder := range logicFolders {
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
				structMatch, err = gregex.MatchString(in.Pattern, structName)
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
		// Generating go files for service.
		for structName, funcArray := range interfaceMap {
			var (
				filePath         = gfile.Join(in.Path, gstr.ToLower(structName)+".go")
				generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceContent, g.MapStrStr{
					"{StructName}":     structName,
					"{PackageName}":    packageName,
					"{FuncDefinition}": funcArray.Join("\n\t"),
				})
			)
			mlog.Printf(`generating service go file: %s`, filePath)
			if err = gfile.PutContents(filePath, generatedContent); err != nil {
				return nil, err
			}
		}
	}
	mlog.Printf(`goimports go files in "%s", it may take seconds...`, in.Path)
	utils.GoImports(in.Path)

	// Replica v1 to v2 for GoFrame.
	err = gfile.ReplaceDirFunc(func(path, content string) string {
		if gstr.Contains(content, `"github.com/gogf/gf`) && !gstr.Contains(content, `"github.com/gogf/gf/v2`) {
			content = gstr.Replace(content, `"github.com/gogf/gf"`, `"github.com/gogf/gf/v2"`)
			content = gstr.Replace(content, `"github.com/gogf/gf/`, `"github.com/gogf/gf/v2/`)
			return content
		}
		return content
	}, in.Path, "*.go", false)
	if err != nil {
		return nil, err
	}
	mlog.Printf(`gofmt go files in "%s"`, in.Path)
	utils.GoFmt(in.Path)
	mlog.Print(`done!`)
	return
}
