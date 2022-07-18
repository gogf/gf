package gendao

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

func generateDo(ctx context.Context, db gdb.DB, tableNames, newTableNames []string, in CGenDaoInternalInput) {
	var (
		doDirPath = gfile.Join(in.Path, in.DoPath)
	)
	in.NoJsonTag = true
	in.DescriptionTag = false
	in.NoModelComment = false
	// Model content.
	for i, tableName := range tableNames {
		in.TableName = tableName
		fieldMap, err := db.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
		}
		var (
			newTableName     = newTableNames[i]
			doFilePath       = gfile.Join(doDirPath, gstr.CaseSnake(newTableName)+".go")
			structDefinition = generateStructDefinition(generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				StructName:           gstr.CaseCamel(newTableName),
				FieldMap:             fieldMap,
				IsDo:                 true,
			})
		)
		// replace all types to interface{}.
		structDefinition, _ = gregex.ReplaceStringFuncMatch(
			"([A-Z]\\w*?)\\s+([\\w\\*\\.]+?)\\s+(//)",
			structDefinition,
			func(match []string) string {
				// If the type is already a pointer/slice/map, it does nothing.
				if !gstr.HasPrefix(match[2], "*") && !gstr.HasPrefix(match[2], "[]") && !gstr.HasPrefix(match[2], "map") {
					return fmt.Sprintf(`%s interface{} %s`, match[1], match[3])
				}
				return match[0]
			},
		)
		modelContent := generateDoContent(
			in,
			tableName,
			gstr.CaseCamel(newTableName),
			structDefinition,
		)
		err = gfile.PutContents(doFilePath, strings.TrimSpace(modelContent))
		if err != nil {
			mlog.Fatalf(`writing content to "%s" failed: %v`, doFilePath, err)
		} else {
			utils.GoFmt(doFilePath)
			mlog.Print("generated:", doFilePath)
		}
	}
}

func generateDoContent(in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string) string {
	doContent := gstr.ReplaceByMap(
		getTemplateFromPathOrDefault(in.TplDaoDoPath, consts.TemplateGenDaoDoContent),
		g.MapStrStr{
			tplVarTableName:          tableName,
			tplVarPackageImports:     getImportPartContent(structDefine, true),
			tplVarTableNameCamelCase: tableNameCamelCase,
			tplVarStructDefine:       structDefine,
		})
	doContent = replaceDefaultVar(in, doContent)
	return doContent
}
