# GaussDB Driver for GoFrame

This package provides a GaussDB database driver for the GoFrame framework.

## Overview

GaussDB is Huawei's enterprise-level database that is compatible with PostgreSQL protocols. This driver adapts the PostgreSQL driver implementation to work with GaussDB.

## Installation

```bash
go get -u github.com/gogf/gf/contrib/drivers/gaussdb/v2
```

## Usage

```go
import (
    _ "github.com/gogf/gf/contrib/drivers/gaussdb/v2"
    "github.com/gogf/gf/v2/database/gdb"
)

// Configuration
gdb.AddConfigNode(gdb.DefaultGroupName, gdb.ConfigNode{
    Link: "gaussdb:username:password@tcp(127.0.0.1:9950)/database_name",
})

// Get database instance
db, err := gdb.New()
```

## Connection String Format

```
gaussdb:username:password@tcp(host:port)/database?param1=value1&param2=value2
```

Example:
```
gaussdb:gaussdb:UTpass@1234@tcp(127.0.0.1:9950)/postgres
```

## Schema/Namespace Handling

GaussDB follows PostgreSQL's schema model:
- **Database (Catalog)**: The database name in the connection string (e.g., `postgres`)
- **Schema (Namespace)**: A namespace within the database (e.g., `public`, `test`)

To use a specific schema:

```go
// Create schema if not exists
db.Exec(ctx, "CREATE SCHEMA IF NOT EXISTS my_schema")

// Set search_path to use the schema
db.Exec(ctx, "SET search_path TO my_schema")
```

## Limitations

GaussDB is based on **PostgreSQL 9.2**, which predates several modern PostgreSQL features. The following features are **NOT SUPPORTED**:

### 1. ON CONFLICT Operations (PostgreSQL 9.5+)

The following ORM methods rely on `ON CONFLICT` syntax and are not available:

- **InsertIgnore()** - Insert and ignore duplicate key errors
- **Save()** - Insert or update (upsert)
- **Replace()** - Replace existing record
- **OnConflict()** - Custom conflict handling
- **OnDuplicate()** - On duplicate key update
- **OnDuplicateEx()** - Extended on duplicate key update

**Workaround**: Use separate INSERT and UPDATE operations with proper error handling:

```go
// Instead of InsertIgnore
result, err := db.Model("user").Insert(data)
if err != nil {
    // Check if error is duplicate key error
    if strings.Contains(err.Error(), "duplicate key") {
        // Handle duplicate - either ignore or update separately
    }
}

// Instead of Save (upsert)
// First try to update
result, err := db.Model("user").Where("id", id).Update(data)
if err != nil {
    return err
}
affected, _ := result.RowsAffected()
if affected == 0 {
    // No rows updated, insert new record
    _, err = db.Model("user").Insert(data)
}
```

## Supported Features

- ✅ Basic CRUD operations (Insert, Select, Update, Delete)
- ✅ Transactions
- ✅ Batch operations
- ✅ Array data types (int, float, text, etc.)
- ✅ JSON/JSONB data types
- ✅ Schema/namespace support
- ✅ Prepared statements
- ✅ Connection pooling

## Testing

To run the test suite, ensure you have a GaussDB instance running:

```bash
# Default test connection
# Host: 127.0.0.1
# Port: 9950
# User: gaussdb
# Password: UTpass@1234
# Database: postgres

go test -v
```

Tests for unsupported features (ON CONFLICT operations) will be skipped with explanatory messages.

## Database Compatibility

- **GaussDB Version**: Based on PostgreSQL 9.2
- **Protocol Compatibility**: PostgreSQL wire protocol
- **Driver**: Uses `gitee.com/opengauss/openGauss-connector-go-pq`

## Notes

1. **Schema Usage**: Unlike MySQL where "schema" and "database" are synonymous, in PostgreSQL/GaussDB:
   - Database (catalog) is the top-level container
   - Schema is a namespace within a database
   - Tables belong to schemas within databases

2. **Connection Database**: Always connect to an existing database (like `postgres`), then create and use schemas within it.

3. **Performance**: For optimal performance, set `search_path` at the session level rather than qualifying every table name with the schema.

4. **Version Checking**: The driver does not enforce GaussDB version checking, but features relying on PostgreSQL 9.5+ functionality will fail.

## Contributing

When contributing to this driver, please note:

1. Test changes against an actual GaussDB instance
2. Ensure compatibility with PostgreSQL 9.2 features only
3. Document any additional limitations discovered
4. Update tests to skip unsupported features appropriately

## License

This driver is distributed under the same license as the GoFrame framework (MIT License).
