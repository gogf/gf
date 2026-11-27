// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies driver-specific schema configuration handling.

package gdb

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_SetSchemaToConfigNode verifies PostgreSQL schema switching preserves its database name.
func Test_SetSchemaToConfigNode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		pgsqlNode := &ConfigNode{
			Type:      driverTypePgSQL,
			Name:      "application",
			Namespace: "public",
		}
		setSchemaToConfigNode(pgsqlNode, "tenant_a")
		t.Assert(pgsqlNode.Name, "application")
		t.Assert(pgsqlNode.Namespace, "tenant_a")

		mysqlNode := &ConfigNode{
			Type: "mysql",
			Name: "application",
		}
		setSchemaToConfigNode(mysqlNode, "tenant_a")
		t.Assert(mysqlNode.Name, "tenant_a")
	})
}
