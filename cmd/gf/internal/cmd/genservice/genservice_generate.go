package genservice

import (
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

type generateServiceFilesInput struct {
	CGenServiceInput
	DstFilePath         string // Absolute file path for generated service go file.
	SrcStructFunctions  map[string]*garray.StrArray
	SrcImportedPackages []string
	SrcPackageName      string
	DstPackageName      string
}

func (c CGenService) generateServiceFile(in generateServiceFilesInput) (ok bool, err error) {
	var (
		generatedContent        string
		allFuncArray            = garray.NewStrArray() // Used for check whether interface dirty, going to change file content.
		importedPackagesContent = fmt.Sprintf(
			"import (\n%s\n)", gstr.Join(in.SrcImportedPackages, "\n"),
		)
	)
	generatedContent += gstr.ReplaceByMap(consts.TemplateGenServiceContentHead, g.MapStrStr{
		"{Imports}":     importedPackagesContent,
		"{PackageName}": in.DstPackageName,
	})

	// Type definitions.
	generatedContent += "type("
	generatedContent += "\n"
	for structName, funcArray := range in.SrcStructFunctions {
		allFuncArray.Append(funcArray.Slice()...)
		generatedContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentInterface, g.MapStrStr{
			"{InterfaceName}":  "I" + structName,
			"{FuncDefinition}": funcArray.Join("\n\t"),
		}))
		generatedContent += "\n"
	}
	generatedContent += ")"
	generatedContent += "\n"

	// Generating variable and register definitions.
	var (
		variableContent          string
		generatingInterfaceCheck string
	)
	// Variable definitions.
	for structName, _ := range in.SrcStructFunctions {
		generatingInterfaceCheck = fmt.Sprintf(`[^\w\d]+%s.I%s[^\w\d]`, in.DstPackageName, structName)
		if gregex.IsMatchString(generatingInterfaceCheck, generatedContent) {
			continue
		}
		variableContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentVariable, g.MapStrStr{
			"{StructName}":    structName,
			"{InterfaceName}": "I" + structName,
		}))
		variableContent += "\n"
	}
	if variableContent != "" {
		generatedContent += "var("
		generatedContent += "\n"
		generatedContent += variableContent
		generatedContent += ")"
		generatedContent += "\n"
	}
	// Variable register function definitions.
	for structName, _ := range in.SrcStructFunctions {
		generatingInterfaceCheck = fmt.Sprintf(`[^\w\d]+%s.I%s[^\w\d]`, in.DstPackageName, structName)
		if gregex.IsMatchString(generatingInterfaceCheck, generatedContent) {
			continue
		}
		generatedContent += gstr.Trim(gstr.ReplaceByMap(consts.TemplateGenServiceContentRegister, g.MapStrStr{
			"{StructName}":    structName,
			"{InterfaceName}": "I" + structName,
		}))
		generatedContent += "\n\n"
	}

	// Replace empty braces that have new line.
	generatedContent, _ = gregex.ReplaceString(`{[\s\t]+}`, `{}`, generatedContent)

	// Remove package name calls of `dstPackageName` in produced codes.
	generatedContent, _ = gregex.ReplaceString(fmt.Sprintf(`\*{0,1}%s\.`, in.DstPackageName), ``, generatedContent)

	// Write file content to disk.
	if gfile.Exists(in.DstFilePath) {
		if !utils.IsFileDoNotEdit(in.DstFilePath) {
			mlog.Printf(`ignore file as it is manually maintained: %s`, in.DstFilePath)
			return false, nil
		}
		if !c.isToGenerateServiceGoFile(in.DstPackageName, in.DstFilePath, allFuncArray) {
			mlog.Printf(`not dirty, ignore generating service go file: %s`, in.DstFilePath)
			return false, nil
		}
	}
	mlog.Printf(`generating service go file: %s`, in.DstFilePath)
	if err = gfile.PutContents(in.DstFilePath, generatedContent); err != nil {
		return true, err
	}
	return true, nil
}

// isToGenerateServiceGoFile checks and returns whether the service content dirty.
func (c CGenService) isToGenerateServiceGoFile(dstPackageName, filePath string, funcArray *garray.StrArray) bool {
	var (
		fileContent        = gfile.GetContents(filePath)
		generatedFuncArray = garray.NewSortedStrArrayFrom(funcArray.Slice())
		contentFuncArray   = garray.NewSortedStrArray()
	)
	if fileContent == "" {
		return true
	}
	matches, _ := gregex.MatchAllString(`\s+interface\s+{([\s\S]+?)}`, fileContent)
	for _, match := range matches {
		contentFuncArray.Append(gstr.SplitAndTrim(match[1], "\n")...)
	}
	if generatedFuncArray.Len() != contentFuncArray.Len() {
		mlog.Debugf(
			`dirty, generatedFuncArray.Len()[%d] != contentFuncArray.Len()[%d]`,
			generatedFuncArray.Len(), contentFuncArray.Len(),
		)
		return true
	}
	var funcDefinition string
	for i := 0; i < generatedFuncArray.Len(); i++ {
		funcDefinition, _ = gregex.ReplaceString(
			fmt.Sprintf(`\*{0,1}%s\.`, dstPackageName), ``, generatedFuncArray.At(i),
		)
		if funcDefinition != contentFuncArray.At(i) {
			mlog.Debugf(`dirty, %s != %s`, funcDefinition, contentFuncArray.At(i))
			return true
		}
	}
	return false
}

func (c CGenService) generateInitializationFile(in CGenServiceInput, importSrcPackages []string) (err error) {
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

// getDstFileNameCase call gstr.Case* function to convert the s to specified case.
func (c CGenService) getDstFileNameCase(str, caseStr string) string {
	switch gstr.ToLower(caseStr) {
	case gstr.ToLower("Lower"):
		return gstr.ToLower(str)

	case gstr.ToLower("Camel"):
		return gstr.CaseCamel(str)

	case gstr.ToLower("CamelLower"):
		return gstr.CaseCamelLower(str)

	case gstr.ToLower("Kebab"):
		return gstr.CaseKebab(str)

	case gstr.ToLower("KebabScreaming"):
		return gstr.CaseKebabScreaming(str)

	case gstr.ToLower("SnakeFirstUpper"):
		return gstr.CaseSnakeFirstUpper(str)

	case gstr.ToLower("SnakeScreaming"):
		return gstr.CaseSnakeScreaming(str)
	}
	return gstr.CaseSnake(str)
}
