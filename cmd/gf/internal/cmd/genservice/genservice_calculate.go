package genservice

import (
	"go/parser"
	"go/token"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func (c CGenService) calculateImportedPackages(fileContent string, srcImportedPackages *garray.SortedStrArray) (err error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", fileContent, parser.ImportsOnly)
	if err != nil {
		return err
	}
	for _, s := range f.Imports {
		if s.Path != nil {
			if s.Name != nil {
				// If it has alias, and it is not `_`.
				if pkgAlias := s.Name.String(); pkgAlias != "_" {
					srcImportedPackages.Add(pkgAlias + " " + s.Path.Value)
				}
			} else {
				// no alias
				srcImportedPackages.Add(s.Path.Value)
			}
		}
	}
	return nil
}

func (c CGenService) calculateInterfaceFunctions(
	in CGenServiceInput, fileContent string, srcPkgInterfaceMap map[string]*garray.StrArray, dstPackageName string,
) (err error) {
	var (
		ok                       bool
		matches                  [][]string
		srcPkgInterfaceFuncArray *garray.StrArray
	)
	// calculate struct name and its functions according function definitions.
	fileContent2 := ""
	fileContentLines := gstr.Split(fileContent, "\n")
	for _, line := range fileContentLines {
		if gstr.HasPrefix(line, "//") {
			continue
		}
		fileContent2 += line + "\n"
	}
	matches, err := gregex.MatchAllString(`func \((.+?)\) ([\s\S]+?) {`, fileContent2)
	if err != nil {
		return err
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

		// Case of:
		// Xxx(\n    ctx context.Context, req *v1.XxxReq,\n) -> Xxx(ctx context.Context, req *v1.XxxReq)
		functionHead = gstr.Replace(functionHead, `,)`, `)`)
		functionHead, _ = gregex.ReplaceString(`\(\s+`, `(`, functionHead)
		functionHead, _ = gregex.ReplaceString(`\s{2,}`, ` `, functionHead)
		if !gstr.IsLetterUpper(functionHead[0]) {
			continue
		}
		// Match and pick the struct name from receiver.
		if structMatch, err = gregex.MatchString(in.StPattern, structName); err != nil {
			return err
		}
		if len(structMatch) < 1 {
			continue
		}
		structName = gstr.CaseCamel(structMatch[1])
		if srcPkgInterfaceFuncArray, ok = srcPkgInterfaceMap[structName]; !ok {
			srcPkgInterfaceMap[structName] = garray.NewStrArray()
			srcPkgInterfaceFuncArray = srcPkgInterfaceMap[structName]
		}
		srcPkgInterfaceFuncArray.Append(functionHead)
	}
	// calculate struct name according type definitions.
	matches, err = gregex.MatchAllString(`type (.+) struct\s*{`, fileContent)
	if err != nil {
		return err
	}
	for _, match := range matches {
		var (
			structName  string
			structMatch []string
		)
		if structMatch, err = gregex.MatchString(in.StPattern, match[1]); err != nil {
			return err
		}
		if len(structMatch) < 1 {
			continue
		}
		structName = gstr.CaseCamel(structMatch[1])
		if srcPkgInterfaceFuncArray, ok = srcPkgInterfaceMap[structName]; !ok {
			srcPkgInterfaceMap[structName] = garray.NewStrArray()
		}
	}
	return nil
}
