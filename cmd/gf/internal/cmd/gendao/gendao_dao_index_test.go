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

// Test_Gen_Dao_IndexTemplateImport asserts that the dao index template provides the gdb
// import for the sharding handler, and that the unused import is dropped for the files
// which do not reference it.
func Test_Gen_Dao_IndexTemplateImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dir      = t.TempDir()
			sharding = renderDaoIndex(t, "users", true)
			normal   = renderDaoIndex(t, "users", false)
		)

		// The sharding handler references gdb.Model and gdb.ShardingConfig.
		var shardingPath = gfile.Join(dir, "users.go")
		t.AssertNil(gfile.PutContents(shardingPath, sharding))
		utils.GoFmt(shardingPath)
		var shardingContent = gfile.GetContents(shardingPath)
		t.Assert(gstr.Contains(shardingContent, `"github.com/gogf/gf/v2/database/gdb"`), true)
		// The formatted file must be valid go code.
		_, err := parser.ParseFile(token.NewFileSet(), shardingPath, shardingContent, parser.ParseComments)
		t.AssertNil(err)

		// The non-sharding file does not reference gdb, so the import must be dropped.
		var normalPath = gfile.Join(dir, "users_normal.go")
		t.AssertNil(gfile.PutContents(normalPath, normal))
		utils.GoFmt(normalPath)
		var normalContent = gfile.GetContents(normalPath)
		t.Assert(gstr.Contains(normalContent, "database/gdb"), false)
		_, err = parser.ParseFile(token.NewFileSet(), normalPath, normalContent, parser.ParseComments)
		t.AssertNil(err)
	})
}
