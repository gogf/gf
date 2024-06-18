// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genservice

import (
	"bytes"
	"fmt"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

type generateServiceFilesInput struct {
	CGenServiceInput
	SrcPackageName      string
	SrcImportedPackages []string
	SrcStructFunctions  *gmap.ListMap
	DstPackageName      string
	DstFilePath         string // Absolute file path for generated service go file.
}

func (c CGenService) generateServiceFile(in generateServiceFilesInput) (ok bool, err error) {
	var generatedContent bytes.Buffer

	c.generatePackageImports(&generatedContent, in.DstPackageName, in.SrcImportedPackages)
	c.generateType(&generatedContent, in.SrcStructFunctions, in.DstPackageName)
	c.generateVar(&generatedContent, in.SrcStructFunctions)
	c.generateFunc(&generatedContent, in.SrcStructFunctions)

	// Write file content to disk.
	if gfile.Exists(in.DstFilePath) {
		if !utils.IsFileDoNotEdit(in.DstFilePath) {
			mlog.Printf(`ignore file as it is manually maintained: %s`, in.DstFilePath)
			return false, nil
		}
	}
	mlog.Printf(`generating service go file: %s`, in.DstFilePath)
	if err = gfile.PutBytes(in.DstFilePath, generatedContent.Bytes()); err != nil {
		return true, err
	}
	return true, nil
}

// generateInitializationFile generates `logic.go`.
func (c CGenService) generateInitializationFile(in CGenServiceInput, importSrcPackages []string) (err error) {
	var (
		logicPackageName = gstr.ToLower(gfile.Basename(in.SrcFolder))
		logicFilePath    = gfile.Join(in.SrcFolder, logicPackageName+".go")
		logicImports     string
		generatedContent string
	)
	if !utils.IsFileDoNotEdit(logicFilePath) {
		mlog.Debugf(`ignore file as it is manually maintained: %s`, logicFilePath)
		return nil
	}
	for _, importSrcPackage := range importSrcPackages {
		logicImports += fmt.Sprintf(`%s_ "%s"%s`, "\t", importSrcPackage, "\n")
	}
	generatedContent = gstr.ReplaceByMap(consts.TemplateGenServiceLogicContent, g.MapStrStr{
		"{PackageName}": logicPackageName,
		"{Imports}":     logicImports,
	})
	mlog.Printf(`generating init go file: %s`, logicFilePath)
	if err = gfile.PutContents(logicFilePath, generatedContent); err != nil {
		return err
	}
	utils.GoFmt(logicFilePath)
	return nil
}

// getDstFileNameCase call gstr.Case* function to convert the s to specified case.
func (c CGenService) getDstFileNameCase(str, caseStr string) (newStr string) {
	if newStr := gstr.CaseConvert(str, gstr.CaseTypeMatch(caseStr)); newStr != str {
		return newStr
	}
	return gstr.CaseSnake(str)
}
