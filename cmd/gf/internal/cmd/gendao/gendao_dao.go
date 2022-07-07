package gendao

import (
	"bytes"
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
	"github.com/olekukonko/tablewriter"
)

// generateDaoContentFile generates the dao and model content of given table.
func generateDao(ctx context.Context, db gdb.DB, in CGenDaoInternalInput) {
	// Generating table data preparing.
	fieldMap, err := db.TableFields(ctx, in.TableName)
	if err != nil {
		mlog.Fatalf(`fetching tables fields failed for table "%s": %+v`, in.TableName, err)
	}
	var (
		dirRealPath             = gfile.RealPath(in.Path)
		dirPathDao              = gfile.Join(in.Path, in.DaoPath)
		tableNameCamelCase      = gstr.CaseCamel(in.NewTableName)
		tableNameCamelLowerCase = gstr.CaseCamelLower(in.NewTableName)
		tableNameSnakeCase      = gstr.CaseSnake(in.NewTableName)
		importPrefix            = in.ImportPrefix
	)
	if importPrefix == "" {
		if dirRealPath == "" {
			dirRealPath = in.Path
			importPrefix = dirRealPath
			importPrefix = gstr.Trim(dirRealPath, "./")
		} else {
			importPrefix = gstr.Replace(dirRealPath, gfile.Pwd(), "")
		}
		importPrefix = gstr.Replace(importPrefix, gfile.Separator, "/")
		importPrefix = gstr.Join(g.SliceStr{in.ModName, importPrefix, in.DaoPath}, "/")
		importPrefix, _ = gregex.ReplaceString(`\/{2,}`, `/`, gstr.Trim(importPrefix, "/"))
	} else {
		importPrefix = gstr.Join(g.SliceStr{importPrefix, in.DaoPath}, "/")
	}

	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		fileName += "_table"
	}

	// dao - index
	generateDaoIndex(in, tableNameCamelCase, tableNameCamelLowerCase, importPrefix, dirPathDao, fileName)

	// dao - internal
	generateDaoInternal(in, tableNameCamelCase, tableNameCamelLowerCase, importPrefix, dirPathDao, fileName, fieldMap)
}

func generateDaoIndex(in CGenDaoInternalInput, tableNameCamelCase, tableNameCamelLowerCase, importPrefix, dirPathDao, fileName string) {
	path := gfile.Join(dirPathDao, fileName+".go")
	if in.OverwriteDao || !gfile.Exists(path) {
		indexContent := gstr.ReplaceByMap(getTplDaoIndexContent(""), g.MapStrStr{
			tplVarImportPrefix:            importPrefix,
			tplVarTableName:               in.TableName,
			tplVarTableNameCamelCase:      tableNameCamelCase,
			tplVarTableNameCamelLowerCase: tableNameCamelLowerCase,
		})
		indexContent = replaceDefaultVar(in, indexContent)
		if err := gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", path)
		}
	}
}

func generateDaoInternal(
	in CGenDaoInternalInput,
	tableNameCamelCase, tableNameCamelLowerCase, importPrefix string,
	dirPathDao, fileName string,
	fieldMap map[string]*gdb.TableField,
) {
	path := gfile.Join(dirPathDao, "internal", fileName+".go")
	modelContent := gstr.ReplaceByMap(getTplDaoInternalContent(""), g.MapStrStr{
		tplVarImportPrefix:            importPrefix,
		tplVarTableName:               in.TableName,
		tplVarGroupName:               in.Group,
		tplVarTableNameCamelCase:      tableNameCamelCase,
		tplVarTableNameCamelLowerCase: tableNameCamelLowerCase,
		tplVarColumnDefine:            gstr.Trim(generateColumnDefinitionForDao(fieldMap)),
		tplVarColumnNames:             gstr.Trim(generateColumnNamesForDao(fieldMap)),
	})
	modelContent = replaceDefaultVar(in, modelContent)
	if err := gfile.PutContents(path, strings.TrimSpace(modelContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", path)
	}
}

func getTplDaoIndexContent(tplDaoIndexPath string) string {
	if tplDaoIndexPath != "" {
		return gfile.GetContents(tplDaoIndexPath)
	}
	return consts.TemplateDaoDaoIndexContent
}

func getTplDaoInternalContent(tplDaoInternalPath string) string {
	if tplDaoInternalPath != "" {
		return gfile.GetContents(tplDaoInternalPath)
	}
	return consts.TemplateDaoDaoInternalContent
}

// generateColumnNamesForDao generates and returns the column names assignment content of column struct
// for specified table.
func generateColumnNamesForDao(fieldMap map[string]*gdb.TableField) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)
	for index, name := range names {
		field := fieldMap[name]
		array[index] = []string{
			"            #" + gstr.CaseCamel(field.Name) + ":",
			fmt.Sprintf(` #"%s",`, field.Name),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	namesContent := buffer.String()
	// Let's do this hack of table writer for indent!
	namesContent = gstr.Replace(namesContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(namesContent)
	return buffer.String()
}

// generateColumnDefinitionForDao generates and returns the column names definition for specified table.
func generateColumnDefinitionForDao(fieldMap map[string]*gdb.TableField) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)
	for index, name := range names {
		var (
			field   = fieldMap[name]
			comment = gstr.Trim(gstr.ReplaceByArray(field.Comment, g.SliceStr{
				"\n", " ",
				"\r", " ",
			}))
		)
		array[index] = []string{
			"    #" + gstr.CaseCamel(field.Name),
			" # " + "string",
			" #" + fmt.Sprintf(`// %s`, comment),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	defineContent := buffer.String()
	// Let's do this hack of table writer for indent!
	defineContent = gstr.Replace(defineContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(defineContent)
	return buffer.String()
}
