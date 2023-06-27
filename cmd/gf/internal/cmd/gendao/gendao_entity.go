// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func generateEntity(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathEntity = gfile.Join(in.Path, in.EntityPath)
	if in.Clear {
		doClear(ctx, dirPathEntity, false)
	}
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}

		var (
			newTableName                    = in.NewTableNames[i]
			entityFilePath                  = gfile.Join(dirPathEntity, gstr.CaseSnake(newTableName)+".go")
			structDefinition, appendImports = generateStructDefinition(ctx, generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
				StructName:           gstr.CaseCamel(newTableName),
				FieldMap:             fieldMap,
				IsDo:                 false,
			})
			entityContent = generateEntityContent(
				ctx,
				in,
				newTableName,
				gstr.CaseCamel(newTableName),
				structDefinition,
				appendImports,
			)
		)

		err = gfile.PutContents(entityFilePath, strings.TrimSpace(entityContent))
		if err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", entityFilePath, err)
		} else {
			utils.GoFmt(entityFilePath)
			mlog.Print("generated:", entityFilePath)
		}
	}
}

func generateEntityContent(
	ctx context.Context, in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string, appendImports []string,
) string {
	entityContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoEntityPath, consts.TemplateGenDaoEntityContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(ctx, structDefine, false, appendImports),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
		},
	)
	entityContent = replaceDefaultVar(in, entityContent)
	return entityContent
}
