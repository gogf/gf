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

GaussDB is based on **PostgreSQL 9.2**, which predates several modern PostgreSQL features (like `ON CONFLICT` introduced in PostgreSQL 9.5). However, GaussDB supports the SQL standard `MERGE` statement, which we use to implement some upsert operations.

### Fully Supported UPSERT Operations

All ORM upsert operations are **FULLY SUPPORTED** using `MERGE` statement or alternative implementations:

- ✅ **Save()** - Insert or update (upsert) - Uses MERGE INTO
- ✅ **Replace()** - Replace existing record - Alias for Save()
- ✅ **InsertIgnore()** - Insert and ignore duplicate key errors
  - With primary key in data: Uses MERGE INTO for conflict detection
  - Without primary key: Uses INSERT with error catching
- ✅ **OnConflict()** - Custom conflict column specification - Works with MERGE
- ✅ **OnDuplicate()** - On duplicate key update with custom fields
  - Uses MERGE when not updating conflict keys
  - Uses UPDATE+INSERT when updating conflict keys (GaussDB MERGE limitation workaround)
- ✅ **OnDuplicateEx()** - Exclude specific fields from update - Uses MERGE
- ✅ **OnDuplicateWithCounter()** - Counter operations on duplicate - Fully supported

### Usage Examples

```go
// Basic Save (upsert)
result, err := db.Model("user").Data(data).Save()

// Save with conflict detection on specific column
result, err := db.Model("user").Data(data).OnConflict("email").Save()

// Insert Ignore (skip if exists)
result, err := db.Model("user").Data(data).InsertIgnore()

// OnDuplicate - update specific fields on conflict
result, err := db.Model("user").
    Data(data).
    OnConflict("id").
    OnDuplicate("name", "email").
    Save()

// OnDuplicateEx - update all except specified fields
result, err := db.Model("user").
    Data(data).
    OnConflict("id").
    OnDuplicateEx("created_at").
    Save()

// OnDuplicate with Counter
result, err := db.Model("user").
    Data(data).
    OnConflict("id").
    OnDuplicate(g.Map{
        "login_count": gdb.Counter{Field: "login_count", Value: 1},
    }).
    Save()

// OnDuplicate with Raw SQL
result, err := db.Model("user").
    Data(data).
    OnConflict("id").
    OnDuplicate(g.Map{
        "updated_at": gdb.Raw("CURRENT_TIMESTAMP"),
    }).
    Save()
```

### Implementation Notes

1. **MERGE Statement**: GaussDB supports the SQL standard MERGE statement, which is used for most upsert operations
2. **Conflict Key Updates**: When OnDuplicate attempts to update a conflict key (e.g., primary key), MERGE cannot be used. In this case, the driver automatically falls back to UPDATE+INSERT approach
3. **EXCLUDED Keyword**: PostgreSQL's `EXCLUDED` (used in ON CONFLICT) is automatically converted to the MERGE equivalent `T2` prefix
4. **Atomic Operations**: All operations maintain atomicity and consistency

## Supported Features

- ✅ Basic CRUD operations (Insert, Select, Update, Delete)
- ✅ Transactions
- ✅ Batch operations
- ✅ Array data types (int, float, text, etc.)
- ✅ JSON/JSONB data types
- ✅ Schema/namespace support
- ✅ Prepared statements
- ✅ Connection pooling

## Supported Features

- ✅ Basic CRUD operations (Insert, Select, Update, Delete)
- ✅ **Save/Upsert operations** (using MERGE statement)
- ✅ **InsertIgnore** (using MERGE statement)  
- ✅ **Replace** (using MERGE statement)
- ✅ Transactions
- ✅ Batch operations
- ✅ Array data types (int, float, text, etc.)
- ✅ JSON/JSONB data types
- ✅ Schema/namespace support
- ✅ Prepared statements
- ✅ Connection pooling4
# Database: postgres

Tests for unsupported features (OnConflict/OnDuplicate operations requiring ON CONFLICT syntax) will be skipped with explanatory messages. Tests for Save/InsertIgnore operations (using MERGE statement) will pass successfully.
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
