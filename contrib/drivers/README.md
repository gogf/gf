# drivers
Database drivers for package gdb.

# Installation
Let's take `pgsql` for example.
```
go get -u github.com/gogf/gf/contrib/drivers/pgsql/v2
```

Choose and import the driver to your project:
```
import _ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
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
Note:
- It does not support `Save/Replace` features.

## PostgreSQL
```
import _ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
```
Note:
- It does not support `Save/Replace` features.
- It does not support `LastInsertId`.

## SQL Server
```
import _ "github.com/gogf/gf/contrib/drivers/mssql/v2"
```
Note:
- It does not support `Save/Replace` features.
- It does not support `LastInsertId`.
- It supports server version >= `SQL Server2005`

## Oracle
```
import _ "github.com/gogf/gf/contrib/drivers/oracle/v2"
```
Note:
- It does not support `Save/Replace` features.
- It does not support `LastInsertId`.

## ClickHouse
```
import _ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
```
Note:
- It does not support `InsertIgnore/InsertGetId` features.
- It does not support `Ignore/Replace` features.
- It does not support `Transaction` feature.
- It does not support `RowsAffected` feature.


# Custom Drivers

It's quick and easy, please refer to current driver source. 
It's quite appreciated if any PR for new drivers support into current repo.
