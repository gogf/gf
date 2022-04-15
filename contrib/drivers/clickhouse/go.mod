module github.com/gogf/gf/contrib/drivers/clickhouse/v2

go 1.15

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.0.12
	github.com/gogf/gf/v2 v2.0.0-rc
)

replace (
	github.com/ClickHouse/clickhouse-go/v2 => github.com/DGuang21/clickhouse-go/v2 v2.0.14
	github.com/gogf/gf/v2 => ../../../
)
