package gendao

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

func generateDo(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathDo = gfile.Join(in.Path, in.DoPath)
	if in.Clear {
		doClear(ctx, dirPathDo)
	}
	in.NoJsonTag = true
	in.DescriptionTag = false
	in.NoModelComment = false
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}
		var (
			newTableName     = in.NewTableNames[i]
			doFilePath       = gfile.Join(dirPathDo, gstr.CaseSnake(newTableName)+".go")
			structDefinition = generateStructDefinition(ctx, generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
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
