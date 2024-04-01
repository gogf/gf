# drivers
Powerful database drivers for package gdb.

# Installation
Let's take `mysql` for example.
```
go get -u github.com/gogf/gf/contrib/drivers/mysql/v2
```

Choose and import the driver to your project:
```
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
```

Commonly imported at top of `main.go`:
```go
package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	// Other imported packages.
)

func main() {
	// Main logics.
}
```

# Supported Drivers

## MySQL/MariaDB/TiDB

```
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
```

## SQLite
```
import _ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
```

## PostgreSQL
```
import _ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
```
Note:
- It does not support `Replace` features.

## SQL Server
```
import _ "github.com/gogf/gf/contrib/drivers/mssql/v2"
```
Note:
- It does not support `Replace` features.
- It does not support `LastInsertId`.
- It supports server version >= `SQL Server2005`
- It ONLY supports datetime2 and datetimeoffset types for auto handling created_at/updated_at/deleted_at columns, because datetime type does not support microseconds precision when column value is passed as string. 

## Oracle
```
import _ "github.com/gogf/gf/contrib/drivers/oracle/v2"
```
Note:
- It does not support `Replace` features.
- It does not support `LastInsertId`.

## ClickHouse
```
import _ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
```
Note:
- It does not support `InsertIgnore/InsertGetId` features.
- It does not support `Save/Replace` features.
- It does not support `Transaction` feature.
- It does not support `RowsAffected` feature.

## DM
```
import _ "github.com/gogf/gf/contrib/drivers/dm/v2"
```
Note:
- It does not support `Replace` features.



# Custom Drivers

It's quick and easy, please refer to current driver source. 
It's quite appreciated if any PR for new drivers support into current repo.
