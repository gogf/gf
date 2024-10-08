// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// DoCommit commits current sql and arguments to underlying sql driver.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	out, err = d.Core.DoCommit(ctx, in)
	if err != nil {
		return
	}
	if len(out.Records) > 0 {
		// remove auto added field.
		for i, record := range out.Records {
			delete(record, rowNumberAliasForSelect)
			out.Records[i] = record
		}
	}
	return
}
