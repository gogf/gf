// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"bytes"
	"context"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

// generateTable generates dao files for given tables.
func generateTable(ctx context.Context, in CGenDaoInternalInput) {
	dirPathTable := gfile.Join(in.Path, in.TablePath)
	if !in.GenTable {
		if gfile.Exists(dirPathTable) {
			in.genItems.AppendDirPath(dirPathTable)
		}
		return
	}
	in.genItems.AppendDirPath(dirPathTable)
	for i := 0; i < len(in.TableNames); i++ {
		var (
			realTableName = in.TableNames[i]
			newTableName  = in.NewTableNames[i]
		)
		generateTableSingle(ctx, generateTableSingleInput{
			CGenDaoInternalInput: in,
			TableName:            realTableName,
			NewTableName:         newTableName,
			DirPathTable:         dirPathTable,
		})
	}
}

type generateTableSingleInput struct {
	CGenDaoInternalInput
	// TableName specifies the table name of the table.
	TableName string
	// NewTableName specifies the prefix-stripped or custom edited name of the table.
	NewTableName string
	DirPathTable string
}

// generateTableSingle generates dao files for a single table.
func generateTableSingle(ctx context.Context, in generateTableSingleInput) {
	// Generating table data preparing.
	fieldMap, err := in.DB.TableFields(ctx, in.TableName)
	if err != nil {
		mlog.Fatalf(`fetching tables fields failed for table "%s": %+v`, in.TableName, err)
	}

	tableNameSnakeCase := gstr.CaseSnake(in.NewTableName)
	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		fileName += "_table"
	}
	path := filepath.FromSlash(gfile.Join(in.DirPathTable, fileName+".go"))
	in.genItems.AppendGeneratedFilePath(path)
	if in.OverwriteDao || !gfile.Exists(path) {
		var (
			ctx        = context.Background()
			tplContent = getTemplateFromPathOrDefault(
				in.TplDaoTablePath, consts.TemplateGenTableContent,
			)
		)
		tplView.ClearAssigns()
		tplView.Assigns(gview.Params{
			tplVarGroupName:          in.Group,
			tplVarTableName:          in.TableName,
			tplVarTableNameCamelCase: formatFieldName(in.NewTableName, FieldNameCaseCamel),
			tplVarPackageName:        filepath.Base(in.TablePath),
			tplVarTableFields:        generateTableFields(fieldMap),
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

// generateTableFields generates and returns the field definition content for specified table.
func generateTableFields(fields map[string]*gdb.TableField) string {
	var buf bytes.Buffer
	fieldNames := make([]string, 0, len(fields))
	for fieldName := range fields {
		fieldNames = append(fieldNames, fieldName)
	}
	sort.Slice(fieldNames, func(i, j int) bool {
		return fields[fieldNames[i]].Index < fields[fieldNames[j]].Index // 升序
	})
	for index, fieldName := range fieldNames {
		field := fields[fieldName]
		buf.WriteString("    " + strconv.Quote(field.Name) + ": {\n")
		buf.WriteString("        Index:   " + gconv.String(field.Index) + ",\n")
		buf.WriteString("        Name:    " + strconv.Quote(field.Name) + ",\n")
		buf.WriteString("        Type:    " + strconv.Quote(field.Type) + ",\n")
		buf.WriteString("        Null:    " + gconv.String(field.Null) + ",\n")
		buf.WriteString("        Key:     " + strconv.Quote(field.Key) + ",\n")
		buf.WriteString("        Default: " + generateDefaultValue(field.Default) + ",\n")
		buf.WriteString("        Extra:   " + strconv.Quote(field.Extra) + ",\n")
		buf.WriteString("        Comment: " + strconv.Quote(field.Comment) + ",\n")
		buf.WriteString("    },")
		if index != len(fieldNames)-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// generateDefaultValue generates and returns the default value definition for specified field.
func generateDefaultValue(value interface{}) string {
	if value == nil {
		return "nil"
	}
	switch v := value.(type) {
	case string:
		return strconv.Quote(v)
	default:
		return gconv.String(v)
	}
}
