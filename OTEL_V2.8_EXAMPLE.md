# OpenTelemetry V2.8 Improvements Example

This example demonstrates the configurable OpenTelemetry tracing features for SQL, HTTP requests, and HTTP responses.

## HTTP Server Configuration

### YAML Configuration
```yaml
server:
  address: ":8080"
  otelTraceRequestEnabled: true   # Enable HTTP request parameter tracing
  otelTraceResponseEnabled: true  # Enable HTTP response body tracing
```

### Programmatic Configuration
```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

func main() {
    s := g.Server()
    
    // Method 1: Using SetConfigWithMap
    s.SetConfigWithMap(g.Map{
        "otelTraceRequestEnabled":  true,
        "otelTraceResponseEnabled": true,
    })
    
    // Method 2: Using ServerConfig struct
    config := ghttp.NewConfig()
    config.Address = ":8080"
    config.OtelTraceRequestEnabled = true
    config.OtelTraceResponseEnabled = true
    s.SetConfig(config)
    
    s.BindHandler("/api/test", func(r *ghttp.Request) {
        r.Response.WriteJson(g.Map{
            "message": "Hello World",
            "input":   r.Get("input"),
        })
    })
    
    s.Run()
}
```

## Database Configuration

### YAML Configuration
```yaml
database:
  default:
    type: "mysql"
    host: "127.0.0.1"
    port: "3306"
    user: "your_user"
    pass: "your_password"
    name: "your_database"
    otelTraceSQLEnabled: true  # Enable SQL statement tracing
```

### Programmatic Configuration
```go
package main

import (
    "github.com/gogf/gf/v2/database/gdb"
)

func main() {
    // Configure database with SQL tracing enabled
    config := gdb.ConfigNode{
        Type:                "mysql",
        Host:                "127.0.0.1",
        Port:                "3306",
        User:                "your_user",
        Pass:                "your_password",
        Name:                "your_database",
        OtelTraceSQLEnabled: true,
    }
    
    db, err := gdb.New(config)
    if err != nil {
        panic(err)
    }
    
    // SQL statements will now be traced
    result, err := db.Query("SELECT * FROM users WHERE id = ?", 1)
    // ... handle result
}
```

## Trace Output Examples

### HTTP Method Tracing
All HTTP requests now include the HTTP method in traces:
- `http.method: GET`
- `http.method: POST`  
- `http.method: PUT`
- `http.method: DELETE`

### Request Parameter Tracing (when enabled)
```json
{
  "http.request.params": {
    "username": "john",
    "email": "john@example.com",
    "query_param": "value"
  }
}
```

### Response Body Tracing (when enabled)
```json
{
  "http.response.body": {
    "code": 200,
    "message": "Success",
    "data": {"id": 1, "name": "John Doe"}
  }
}
```

### SQL Tracing (when enabled)
```json
{
  "db.execution.sql": "SELECT * FROM users WHERE id = ? AND status = ?",
  "db.execution.cost": "15 ms",
  "db.execution.rows": "1"
}
```

## Benefits

1. **Configurable**: All tracing features are opt-in via configuration
2. **Performance**: Only enabled features add overhead
3. **Simple**: Flat boolean fields, no complex nested structures
4. **Comprehensive**: Covers SQL, HTTP requests, and HTTP responses
5. **Size Aware**: Respects content size limits to prevent memory issues
6. **YAML Support**: Can be configured via configuration file