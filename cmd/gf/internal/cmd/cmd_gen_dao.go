// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v3"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v3"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v3"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v3"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v3"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v3"

	// do not add dm in cli pre-compilation,
	// the dm driver does not support certain target platforms.
	// _ "github.com/gogf/gf/contrib/drivers/dm/v3"
	"github.com/gogf/gf/cmd/gf/v3/internal/cmd/gendao"
)

type (
	cGenDao = gendao.CGenDao
)
