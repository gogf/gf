# 数据库驱动程序

用于gdb包的数据库驱动程序。

## 安装

以 `mysql` 为例。

```shell
go get -u github.com/gogf/gf/contrib/drivers/mysql/v2
# 方便复制
go get -u github.com/gogf/gf/contrib/drivers/clickhouse/v2
go get -u github.com/gogf/gf/contrib/drivers/dm/v2
go get -u github.com/gogf/gf/contrib/drivers/mssql/v2
go get -u github.com/gogf/gf/contrib/drivers/oracle/v2
go get -u github.com/gogf/gf/contrib/drivers/pgsql/v2
go get -u github.com/gogf/gf/contrib/drivers/sqlite/v2
go get -u github.com/gogf/gf/contrib/drivers/sqlitecgo/v2
```

选择并将驱动程序导入到您的项目中：

```go
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
```

通常在 `main.go` 的顶部导入：

```go
package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	// 其他导入的包。
)

func main() {
	// 主要逻辑。
}
```

## 支持的驱动程序

### MySQL/MariaDB/TiDB

```go
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
```

### SQLite

```go
import _ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
```

#### cgo 版本

32位Windows请使用cgo版本

```go
import _ "github.com/gogf/gf/contrib/drivers/sqlitecgo/v2"
```

### PostgreSQL

```go
import _ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
```

注意：

- 不支持 `Replace` 功能。

### SQL Server

```go
import _ "github.com/gogf/gf/contrib/drivers/mssql/v2"
```

注意：

- 不支持 `Replace` 功能。
- 不支持 `LastInsertId`。
- 仅支持服务器版本 >= `SQL Server2005`
- 仅支持 datetime2 和 datetimeoffset 类型来自动处理 created_at/updated_at/deleted_at 列，因为 datetime 类型在将列值作为字符串传递时不支持微秒精度。

### Oracle

```go
import _ "github.com/gogf/gf/contrib/drivers/oracle/v2"
```

注意：

- 不支持 `Replace` 功能。
- 不支持 `LastInsertId`。

### ClickHouse

```go
import _ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
```

注意：

- 不支持 `InsertIgnore/InsertGetId` 功能。
- 不支持 `Save/Replace` 功能。
- 不支持 `Transaction` 功能。
- 不支持 `RowsAffected` 功能。

### DM

```go
import _ "github.com/gogf/gf/contrib/drivers/dm/v2"
```

注意：

- 不支持 `Replace` 功能。

## 自定义驱动程序

自定义驱动程序非常快速和简单，您可以参考当前驱动程序的源代码来进行开发。
如果您有关于支持新驱动程序的PR（Pull Request），我们将非常感激地接受您的提交到当前仓库。