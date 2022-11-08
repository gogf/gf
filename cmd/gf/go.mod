module github.com/gogf/gf/cmd/gf/v2

go 1.16

require (
	github.com/gogf/gf/contrib/drivers/mssql/v2 v2.1.0
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.1.0
	github.com/gogf/gf/contrib/drivers/pgsql/v2 v2.1.0
	github.com/gogf/gf/contrib/drivers/sqlite/v2 v2.1.0
	github.com/gogf/gf/v2 v2.1.0
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/rivo/uniseg v0.4.2 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/tools v0.1.11
)

replace (
	github.com/gogf/gf/contrib/drivers/mssql/v2 => ../../contrib/drivers/mssql/
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/drivers/pgsql/v2 => ../../contrib/drivers/pgsql/
	github.com/gogf/gf/contrib/drivers/sqlite/v2 => ../../contrib/drivers/sqlite/
	github.com/gogf/gf/v2 => ../../
)
