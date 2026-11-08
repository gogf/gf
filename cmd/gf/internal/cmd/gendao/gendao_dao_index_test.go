// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file provides tests for the generated dao index files, which are formatted by
// GoFmt after being rendered from the dao index template.

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

// daoIndexGdbImportLine is the rendered gdb import line of the dao index file.
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

// Test_Gen_Dao_IndexTemplateImport asserts that the dao index template does not hard-code the
// gdb import, and that the import is ensured by the generator for the sharding tables only.
func Test_Gen_Dao_IndexTemplateImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// The template itself is import-minimal: gdb is not part of it, so that the
		// non-sharding files never carry an unused import.
		t.Assert(gstr.Contains(consts.TemplateGenDaoIndexContent, "database/gdb"), false)
		t.Assert(gstr.Contains(renderDaoIndex(t, "users", true), "database/gdb"), false)
	})

	gtest.C(t, func(t *gtest.T) {
		var dir = t.TempDir()

		// The sharding handler references gdb.Model, so the generator ensures the import.
		var shardingPath = gfile.Join(dir, "users.go")
		t.AssertNil(gfile.PutContents(shardingPath, renderDaoIndex(t, "users", true)))
		_, err := utils.EnsureImports(shardingPath, utils.ImportSpec{Path: daoIndexGdbImportPath})
		t.AssertNil(err)
		utils.GoFmt(shardingPath)
		var shardingContent = gfile.GetContents(shardingPath)
		t.Assert(gstr.Contains(shardingContent, daoIndexGdbImportLine), true)
		// The formatted file must be valid go code.
		_, err = parser.ParseFile(token.NewFileSet(), shardingPath, shardingContent, parser.ParseComments)
		t.AssertNil(err)

		// The non-sharding file does not reference gdb, so no import is ensured and the
		// formatted file stays free of it.
		var normalPath = gfile.Join(dir, "users_normal.go")
		t.AssertNil(gfile.PutContents(normalPath, renderDaoIndex(t, "users", false)))
		utils.GoFmt(normalPath)
		var normalContent = gfile.GetContents(normalPath)
		t.Assert(gstr.Contains(normalContent, "database/gdb"), false)
		_, err = parser.ParseFile(token.NewFileSet(), normalPath, normalContent, parser.ParseComments)
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		// EnsureImports is idempotent for the dao index file, so that re-generating an
		// existing sharding dao file does not duplicate the import.
		var (
			path    = gfile.Join(t.TempDir(), "users.go")
			content = renderDaoIndex(t, "users", true)
		)
		t.AssertNil(gfile.PutContents(path, content))
		changed, err := utils.EnsureImports(path, utils.ImportSpec{Path: daoIndexGdbImportPath})
		t.AssertNil(err)
		t.Assert(changed, true)
		changed, err = utils.EnsureImports(path, utils.ImportSpec{Path: daoIndexGdbImportPath})
		t.AssertNil(err)
		t.Assert(changed, false)
		t.Assert(gstr.Count(gfile.GetContents(path), daoIndexGdbImportLine), 1)
	})
}
