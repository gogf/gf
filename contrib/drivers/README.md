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

## MySQL

BuiltIn supported, nothing todo.

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

## Clickhouse
```
import _ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
```
Note:
- It does not support `Replace/Ignore` features.
- It does not support `LastInsertId`.
- It does not support `Transaction`.
- It does not support `RowsAffected`.

# Custom Drivers

It's quick and easy, please refer to current driver source. 
It's quite appreciated if any PR for new drivers support into current repo.
