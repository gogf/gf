// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func generateEntity(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathEntity = gfile.Join(in.Path, in.EntityPath)
	in.genItems.AppendDirPath(dirPathEntity)
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}

		var (
			newTableName                    = in.NewTableNames[i]
			entityFilePath                  = filepath.FromSlash(gfile.Join(dirPathEntity, gstr.CaseSnake(newTableName)+".go"))
			structDefinition, appendImports = generateStructDefinition(ctx, generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
				StructName:           formatFieldName(newTableName, FieldNameCaseCamel),
				FieldMap:             fieldMap,
				IsDo:                 false,
			})
			entityContent = generateEntityContent(
				ctx,
				in,
				newTableName,
				formatFieldName(newTableName, FieldNameCaseCamel),
				structDefinition,
				appendImports,
			)
		)
		in.genItems.AppendGeneratedFilePath(entityFilePath)
		err = gfile.PutContents(entityFilePath, strings.TrimSpace(entityContent))
		if err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", entityFilePath, err)
		} else {
			utils.GoFmt(entityFilePath)
			mlog.Print("generated:", gfile.RealPath(entityFilePath))
		}
	}
}

func generateEntityContent(
	ctx context.Context, in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string, appendImports []string,
) string {
	var (
		tplContent = getTemplateFromPathOrDefault(
			in.TplDaoEntityPath, consts.TemplateGenDaoEntityContent,
		)
	)
	tplView.ClearAssigns()
	tplView.Assigns(gview.Params{
		tplVarTableName:          tableName,
		tplVarPackageImports:     getImportPartContent(ctx, structDefine, false, appendImports),
		tplVarTableNameCamelCase: tableNameCamelCase,
		tplVarStructDefine:       structDefine,
		tplVarPackageName:        filepath.Base(in.EntityPath),
	})
	assignDefaultVar(tplView, in)
	entityContent, err := tplView.ParseContent(ctx, tplContent)
	if err != nil {
		mlog.Fatalf("parsing template content failed: %v", err)
	}
	return entityContent
}
