// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/stretchr/testify/assert"
)

// fix https://github.com/gogf/gf/issues/1531
func Test_Func_PG_ConvertDataForRecord(t *testing.T) {
	d := new(gdb.DriverPgsql)
	data := d.ConvertDataForRecord(ctx, g.Map{
		"code":  "test",
		"path":  []string{"c700a87b-e4d8-4aa1-aa18-38ebe107d0ae", "330cba76-8a69-4321-b783-199c53df64ae"},
		"path2": []string{},
	})
	assert.Equal(t, gdb.Raw("ARRAY['c700a87b-e4d8-4aa1-aa18-38ebe107d0ae','330cba76-8a69-4321-b783-199c53df64ae']"), data["path"])
	assert.Equal(t, gdb.Raw("ARRAY[]"), data["path2"])
}
