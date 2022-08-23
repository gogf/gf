package gendao

import (
	"context"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
)

func generateEntity(ctx context.Context, db gdb.DB, tableNames, newTableNames []string, in CGenDaoInternalInput) {
	var entityDirPath = gfile.Join(in.Path, in.EntityPath)
	// Model content.
	for i, tableName := range tableNames {
		fieldMap, err := db.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
		}
		var (
			newTableName   = newTableNames[i]
			entityFilePath = gfile.Join(entityDirPath, gstr.CaseSnake(newTableName)+".go")
			entityContent  = generateEntityContent(
				in,
				newTableName,
				gstr.CaseCamel(newTableName),
				generateStructDefinition(ctx, generateStructDefinitionInput{
					CGenDaoInternalInput: in,
					DB:                   db,
					StructName:           gstr.CaseCamel(newTableName),
					FieldMap:             fieldMap,
					IsDo:                 false,
				}),
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

func generateEntityContent(in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string) string {
	entityContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoEntitylPath, consts.TemplateGenDaoEntityContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(structDefine, false),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
		})
	entityContent = replaceDefaultVar(in, entityContent)
	return entityContent
}
