module github.com/gogf/gf/contrib/drivers/clickhouse/v2

go 1.15

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.0.15
	github.com/gogf/gf/v2 v2.4.0
	github.com/google/uuid v1.3.0
	github.com/shopspring/decimal v1.3.1
)

replace (
	github.com/ClickHouse/clickhouse-go/v2 => github.com/gogf/clickhouse-go/v2 v2.0.15-compatible
	github.com/gogf/gf/v2 => ../../../
)
