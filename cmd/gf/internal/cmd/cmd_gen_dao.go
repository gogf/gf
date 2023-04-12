package cmd

import (
	// DO NOT import clickhouse in default, as it will fail building cli binary in some platforms.
	// You can change the imports here and build the cli binary manually if clickhouse is necessary for you.
	// _ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"

	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/gendao"
)

type (
	cGenDao = gendao.CGenDao
)
