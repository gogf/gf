// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file provides tests for the dao index file rendered from its template, including the
// imports that GoFmt adds for the sharding handler.

package gendao

import (
	"context"
	"go/parser"
	"go/token"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

// daoIndexGdbImportLine is the gdb import line of the formatted dao index file.
const daoIndexGdbImportLine = `"github.com/gogf/gf/v2/database/gdb"`

// renderDaoIndex renders the dao index template for the given table.
func renderDaoIndex(t *gtest.T, tableName string, sharding bool) string {
	var (
		view    = gview.New()
		content string
		err     error
	)
	view.Assigns(gview.Params{
		"TplPackageName":             "dao",
		"TplImportPrefix":            "github.com/gogf/gf/cmd/gf/v2/internal/cmd/testdata/gendao",
		"TplTableName":               tableName,
		"TplTableNameCamelCase":      gstr.UcFirst(tableName),
		"TplTableNameCamelLowerCase": gstr.LcFirst(tableName),
		"TplTableSharding":           sharding,
		"TplTableShardingPrefix":     tableName + "_",
	})
	content, err = view.ParseContent(context.Background(), consts.TemplateGenDaoIndexContent)
	t.AssertNil(err)
	return content
}

// Test_Gen_Dao_IndexTemplateImport asserts that the sharding handler of the rendered dao index
// file compiles: the template declares no gdb import, and GoFmt adds it because the rendered
// sharding code references gdb.
func Test_Gen_Dao_IndexTemplateImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The template itself is import-minimal: gdb is not part of it, so that the
		// rendered non-sharding content never references an import that it does not use.
		t.Assert(gstr.Contains(consts.TemplateGenDaoIndexContent, "database/gdb"), false)
		t.Assert(gstr.Contains(renderDaoIndex(t, "users", true), "database/gdb"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			dir  = t.TempDir()
			path = gfile.Join(dir, "users.go")
		)
		t.AssertNil(gfile.PutContents(path, renderDaoIndex(t, "users", true)))
		utils.GoFmt(path)

		var (
			content = gfile.GetContents(path)
			_, err  = parser.ParseFile(token.NewFileSet(), path, content, parser.ParseComments)
		)
		// The sharding handler references gdb.Model and gdb.ShardingConfig.
		t.Assert(gstr.Contains(content, "gdb.ShardingConfig"), true)
		t.Assert(gstr.Contains(content, daoIndexGdbImportLine), true)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		var (
			dir  = t.TempDir()
			path = gfile.Join(dir, "single_table.go")
		)
		t.AssertNil(gfile.PutContents(path, renderDaoIndex(t, "single_table", false)))
		utils.GoFmt(path)

		var content = gfile.GetContents(path)
		// The non-sharding content does not reference gdb, so the import is never added.
		t.Assert(gstr.Contains(content, "database/gdb"), false)
		t.Assert(gstr.Contains(content, `/testdata/gendao/internal"`), true)
	})
}
