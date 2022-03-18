module github.com/gogf/gf/cmd/gf/v2

go 1.15

require (
	github.com/gogf/gf/contrib/drivers/mssql/v2 v2.0.0-rc2
	github.com/gogf/gf/contrib/drivers/pgsql/v2 v2.0.0-rc2
	github.com/gogf/gf/contrib/drivers/sqlite/v2 v2.0.0-rc2
	github.com/gogf/gf/v2 v2.0.0
	github.com/olekukonko/tablewriter v0.0.5
)

replace (
	github.com/gogf/gf/contrib/drivers/mssql/v2 => ../../contrib/drivers/mssql/
	github.com/gogf/gf/contrib/drivers/pgsql/v2 => ../../contrib/drivers/pgsql/
	github.com/gogf/gf/contrib/drivers/sqlite/v2 => ../../contrib/drivers/sqlite/
	github.com/gogf/gf/v2 => ../../
)
