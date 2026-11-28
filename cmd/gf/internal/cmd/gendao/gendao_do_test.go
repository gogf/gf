// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// Test_Gen_Dao_DoTypeMappingImport is the regression test for
// https://github.com/gogf/gf/issues/5050: gen dao dropped the typeMapping
// import in DO files. A mapped pointer type survives the any-replacement in
// generateDo, so the generated DO file references a package it never imports
// and no longer compiles. The import must come from the appendImports
// returned by generateStructDefinition, the same way generateEntity does it.
func Test_Gen_Dao_DoTypeMappingImport(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx = context.Background()
			in  = CGenDaoInternalInput{
				CGenDaoInput: CGenDaoInput{
					DoPath:     "internal/model/do",
					EntityPath: "internal/model/entity",
					TypeMapping: map[DBFieldTypeName]CustomAttributeType{
						"numeric": {
							Type:   "*decimal.Decimal",
							Import: "github.com/shopspring/decimal",
						},
					},
				},
			}
			fieldMap = map[string]*gdb.TableField{
				"amount": {
					Index:   0,
					Name:    "amount",
					Type:    "numeric",
					Comment: "Amount",
				},
			}
		)
		structDefinitionDo, appendImports := generateStructDefinition(ctx, generateStructDefinitionInput{
			CGenDaoInternalInput: in,
			TableName:            "xxx",
			StructName:           "Xxx",
			FieldMap:             fieldMap,
			IsDo:                 true,
		})
		// The helper returns the import for the mapped type.
		t.Assert(len(appendImports), 1)
		t.Assert(appendImports[0], "github.com/shopspring/decimal")
		t.Assert(gstr.Contains(structDefinitionDo, "*decimal.Decimal"), true)

		// The DO content must carry the import once appendImports are passed.
		doContent := generateDoContent(ctx, in, "xxx", "Xxx", structDefinitionDo, appendImports)
		t.Assert(gstr.Contains(doContent, "*decimal.Decimal"), true)
		t.Assert(gstr.Contains(doContent, `"github.com/shopspring/decimal"`), true)
		// Without appendImports the import is dropped: it comes from the
		// typeMapping, not from scanning the struct text.
		doContentWithoutImports := generateDoContent(ctx, in, "xxx", "Xxx", structDefinitionDo, nil)
		t.Assert(gstr.Contains(doContentWithoutImports, `"github.com/shopspring/decimal"`), false)

		// Control: the entity content already carried the import.
		structDefinitionEntity, appendImportsEntity := generateStructDefinition(ctx, generateStructDefinitionInput{
			CGenDaoInternalInput: in,
			TableName:            "xxx",
			StructName:           "Xxx",
			FieldMap:             fieldMap,
			IsDo:                 false,
		})
		entityContent := generateEntityContent(ctx, in, "xxx", "Xxx", structDefinitionEntity, appendImportsEntity)
		t.Assert(gstr.Contains(entityContent, "*decimal.Decimal"), true)
		t.Assert(gstr.Contains(entityContent, `"github.com/shopspring/decimal"`), true)
	})
}
