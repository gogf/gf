// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

// generateDao generates dao files (index + internal) for all tables in the input.
// It creates the dao directory structure and iterates over each table to generate
// individual dao files via generateDaoSingle.
func generateDao(ctx context.Context, in CGenDaoInternalInput) {
	var (
		dirPathDao         = gfile.Join(in.Path, in.DaoPath)
		dirPathDaoInternal = gfile.Join(dirPathDao, "internal")
	)
	in.genItems.AppendDirPath(dirPathDao)
	for i := 0; i < len(in.TableNames); i++ {
		var (
			realTableName = in.TableNames[i]
			newTableName  = in.NewTableNames[i]
		)
		generateDaoSingle(ctx, generateDaoSingleInput{
			CGenDaoInternalInput: in,
			TableName:            realTableName,
			NewTableName:         newTableName,
			DirPathDao:           dirPathDao,
			DirPathDaoInternal:   dirPathDaoInternal,
			IsSharding:           in.ShardingTableSet.Contains(newTableName),
		})
	}
}

// generateDaoSingleInput holds all parameters needed to generate dao files for a single table.
type generateDaoSingleInput struct {
	CGenDaoInternalInput
	TableName          string // Original table name as it exists in the database.
	NewTableName       string // Processed table name after prefix removal and sharding.
	DirPathDao         string // Directory path for the dao index files.
	DirPathDaoInternal string // Directory path for the dao internal implementation files.
	IsSharding         bool   // Whether this table is a sharding table (merged from multiple physical tables).
}

// generateDaoSingle generates the dao and model content of given table.
func generateDaoSingle(ctx context.Context, in generateDaoSingleInput) {
	// Generating table data preparing.
	fieldMap, err := getTableFields(ctx, in.CGenDaoInternalInput, in.TableName)
	if err != nil {
		mlog.Fatalf(`fetching tables fields failed for table "%s": %+v`, in.TableName, err)
	}
	var (
		tableNameCamelCase      = formatFieldName(in.NewTableName, FieldNameCaseCamel)
		tableNameCamelLowerCase = formatFieldName(in.NewTableName, FieldNameCaseCamelLower)
		tableNameSnakeCase      = gstr.CaseSnake(in.NewTableName)
		importPrefix            = in.ImportPrefix
	)
	if importPrefix == "" {
		importPrefix = utils.GetImportPath(gfile.Join(in.Path, in.DaoPath))
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
	generateDaoIndex(generateDaoIndexInput{
		generateDaoSingleInput:  in,
		TableNameCamelCase:      tableNameCamelCase,
		TableNameCamelLowerCase: tableNameCamelLowerCase,
		ImportPrefix:            importPrefix,
		FileName:                fileName,
	})

	// dao - internal
	generateDaoInternal(generateDaoInternalInput{
		generateDaoSingleInput:  in,
		TableNameCamelCase:      tableNameCamelCase,
		TableNameCamelLowerCase: tableNameCamelLowerCase,
		ImportPrefix:            importPrefix,
		FileName:                fileName,
		FieldMap:                fieldMap,
	})
}

// generateDaoIndexInput holds parameters for generating the dao index file.
// The index file provides the public API (exported struct and constructor)
// for accessing the DAO, delegating to the internal implementation.
type generateDaoIndexInput struct {
	generateDaoSingleInput
	TableNameCamelCase      string // CamelCase version of the table name (e.g., "UserDetail").
	TableNameCamelLowerCase string // camelCase version of the table name (e.g., "userDetail").
	ImportPrefix            string // Go import path prefix for the dao package.
	FileName                string // Output file name (without extension).
}

// generateDaoIndex generates the dao index file for a single table.
// The index file is the public-facing dao file that users import directly.
// It will NOT overwrite an existing file unless OverwriteDao is enabled,
// allowing users to customize the index file without losing changes.
func generateDaoIndex(in generateDaoIndexInput) {
	path := filepath.FromSlash(gfile.Join(in.DirPathDao, in.FileName+".go"))
	// It should add path to result slice whenever it would generate the path file or not.
	in.genItems.AppendGeneratedFilePath(path)
	if in.OverwriteDao || !gfile.Exists(path) {
		var (
			ctx        = context.Background()
			tplContent = getTemplateFromPathOrDefault(
				in.TplDaoIndexPath, consts.TemplateGenDaoIndexContent,
			)
		)
		tplView.ClearAssigns()
		tplView.Assigns(gview.Params{
			tplVarTableSharding:           in.IsSharding,
			tplVarTableShardingPrefix:     in.NewTableName + "_",
			tplVarImportPrefix:            in.ImportPrefix,
			tplVarTableName:               in.TableName,
			tplVarTableNameCamelCase:      in.TableNameCamelCase,
			tplVarTableNameCamelLowerCase: in.TableNameCamelLowerCase,
			tplVarPackageName:             filepath.Base(in.DaoPath),
		})
		indexContent, err := tplView.ParseContent(ctx, tplContent)
		if err != nil {
			mlog.Fatalf("parsing template content failed: %v", err)
		}
		if err = gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", gfile.RealPath(path))
		}
	}
}

// generateDaoInternalInput holds parameters for generating the dao internal file.
// The internal file contains the actual DAO implementation with column definitions
// and is always overwritten on regeneration.
type generateDaoInternalInput struct {
	generateDaoSingleInput
	TableNameCamelCase      string                     // CamelCase version of the table name.
	TableNameCamelLowerCase string                     // camelCase version of the table name.
	ImportPrefix            string                     // Go import path prefix for the dao package.
	FileName                string                     // Output file name (without extension).
	FieldMap                map[string]*gdb.TableField // Map of column name to field metadata.
}

// generateDaoInternal generates the dao internal implementation file for a single table.
// This file is always regenerated (overwritten) and contains the Columns struct definition
// with column name constants and their string value assignments.
func generateDaoInternal(in generateDaoInternalInput) {
	var (
		ctx                    = context.Background()
		removeFieldPrefixArray = gstr.SplitAndTrim(in.RemoveFieldPrefix, ",")
		tplContent             = getTemplateFromPathOrDefault(
			in.TplDaoInternalPath, consts.TemplateGenDaoInternalContent,
		)
	)
	tplView.ClearAssigns()
	tplView.Assigns(gview.Params{
		tplVarImportPrefix:            in.ImportPrefix,
		tplVarTableName:               in.TableName,
		tplVarGroupName:               in.Group,
		tplVarTableNameCamelCase:      in.TableNameCamelCase,
		tplVarTableNameCamelLowerCase: in.TableNameCamelLowerCase,
		tplVarColumnDefine:            gstr.Trim(generateColumnDefinitionForDao(in.FieldMap, removeFieldPrefixArray)),
		tplVarColumnNames:             gstr.Trim(generateColumnNamesForDao(in.FieldMap, removeFieldPrefixArray)),
	})
	assignDefaultVar(tplView, in.CGenDaoInternalInput)
	modelContent, err := tplView.ParseContent(ctx, tplContent)
	if err != nil {
		mlog.Fatalf("parsing template content failed: %v", err)
	}
	path := filepath.FromSlash(gfile.Join(in.DirPathDaoInternal, in.FileName+".go"))
	in.genItems.AppendGeneratedFilePath(path)
	if err := gfile.PutContents(path, strings.TrimSpace(modelContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", gfile.RealPath(path))
	}
}

// generateColumnNamesForDao generates and returns the column names assignment content of column struct
// for specified table.
func generateColumnNamesForDao(fieldMap map[string]*gdb.TableField, removeFieldPrefixArray []string) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)

	for index, name := range names {
		field := fieldMap[name]

		newFiledName := field.Name
		for _, v := range removeFieldPrefixArray {
			newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
		}

		array[index] = []string{
			"            #" + formatFieldName(newFiledName, FieldNameCaseCamel) + ":",
			fmt.Sprintf(` #"%s",`, field.Name),
		}
	}
	table := tablewriter.NewTable(buffer, twRenderer, twConfig)
	table.Bulk(array)
	table.Render()
	namesContent := buffer.String()
	// Let's do this hack of table writer for indent!
	namesContent = gstr.Replace(namesContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(namesContent)
	return buffer.String()
}

// generateColumnDefinitionForDao generates and returns the column names definition for specified table.
func generateColumnDefinitionForDao(fieldMap map[string]*gdb.TableField, removeFieldPrefixArray []string) string {
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
		newFiledName := field.Name
		for _, v := range removeFieldPrefixArray {
			newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
		}
		array[index] = []string{
			"    #" + formatFieldName(newFiledName, FieldNameCaseCamel),
			" # " + "string",
			" #" + fmt.Sprintf(`// %s`, comment),
		}
	}
	table := tablewriter.NewTable(buffer, twRenderer, twConfig)
	table.Bulk(array)
	table.Render()
	defineContent := buffer.String()
	// Let's do this hack of table writer for indent!
	defineContent = gstr.Replace(defineContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(defineContent)
	return buffer.String()
}
