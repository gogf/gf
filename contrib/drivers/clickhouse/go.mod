module github.com/gogf/gf/contrib/drivers/clickhouse/v2

go 1.15

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.0.14
	github.com/gogf/gf/v2 v2.0.0
	github.com/google/uuid v1.3.0 // indirect
)

replace (
	github.com/ClickHouse/clickhouse-go/v2 => github.com/gogf/clickhouse-go/v2 v2.0.14-compatible
	github.com/gogf/gf/v2 => ../../../
)
