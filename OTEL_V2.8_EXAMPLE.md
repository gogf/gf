# OpenTelemetry V2.8 Improvements Example

This example demonstrates the new configurable OpenTelemetry tracing features for SQL, HTTP requests, and HTTP responses.

**Updated to OpenTelemetry v1.38.0 with Independent OTEL Parameters**

## HTTP Server Configuration

### New Independent OTEL Configuration (Recommended)
```yaml
server:
  address: ":8080"
  otel:
    traceRequestEnabled: true   # Enable HTTP request parameter tracing
    traceResponseEnabled: true  # Enable HTTP response body tracing
```

### Legacy Configuration (Still Supported)
```yaml
server:
  address: ":8080"
  otelTraceRequestEnabled: true   # Enable HTTP request parameter tracing  
  otelTraceResponseEnabled: true  # Enable HTTP response body tracing
```

## Database Configuration

### New Independent OTEL Configuration (Recommended)
```yaml
database:
  default:
    type: "mysql"
    host: "127.0.0.1"
    port: "3306"
    user: "your_user"
    pass: "your_password"
    name: "your_database"
    otel:
      traceSQLEnabled: true  # Enable SQL statement tracing
```

### Legacy Configuration (Still Supported)
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

## Programmatic Configuration

### HTTP Server - New Independent OTEL Configuration
```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/internal/otel"
    "github.com/gogf/gf/v2/net/ghttp"
)

func main() {
    s := g.Server()
    
    // Configure using new independent OTEL configuration
    config := ghttp.NewConfig()
    config.Address = ":8080"
    config.Otel = otel.Config{
        TraceRequestEnabled:  true,
        TraceResponseEnabled: true,
    }
    s.SetConfig(config)
    
    s.BindHandler("/api/test", func(r *ghttp.Request) {
        // This handler will have its request parameters and response traced
        r.Response.WriteJson(g.Map{
            "message": "Hello World",
            "input":   r.Get("input"),
        })
    })
    
    s.Run()
}
```

### HTTP Server - Legacy Configuration (Still Supported)
```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

func main() {
    s := g.Server()
    
    // Enable tracing via configuration map (legacy approach)
    s.SetConfigWithMap(g.Map{
        "OtelTraceRequestEnabled":  true,
        "OtelTraceResponseEnabled": true,
    })
    
    s.BindHandler("/api/test", func(r *ghttp.Request) {
        // This handler will have its request parameters and response traced
        r.Response.WriteJson(g.Map{
            "message": "Hello World",  
            "input":   r.Get("input"),
        })
    })
    
    s.Run()
}
```

### Database - New Independent OTEL Configuration
```go
package main

import (
    "github.com/gogf/gf/v2/database/gdb" 
    "github.com/gogf/gf/v2/internal/otel"
)

func main() {
    // Configure database with new independent OTEL configuration
    config := gdb.ConfigNode{
        Type: "mysql",
        Host: "127.0.0.1",
        Port: "3306",
        User: "your_user",
        Pass: "your_password",
        Name: "your_database",
        Otel: otel.Config{
            TraceSQLEnabled: true,
        },
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

### Database - Legacy Configuration (Still Supported)  
```go
package main

import (
    "github.com/gogf/gf/v2/database/gdb"
)

func main() {
    // Configure database with legacy OTEL configuration
    config := gdb.ConfigNode{
        Type: "mysql",
        Host: "127.0.0.1", 
        Port: "3306",
        User: "your_user",
        Pass: "your_password",
        Name: "your_database",
        OtelTraceSQLEnabled: true,  // Legacy field
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

1. **OpenTelemetry v1.38.0**: Updated to the latest OpenTelemetry version with improved performance and features
2. **Independent Configuration**: New modular OTEL configuration structure for better organization
3. **Configurable**: All new tracing features are opt-in via configuration
4. **Performance**: Only enabled features add overhead 
5. **Backward Compatible**: Legacy configuration fields still work alongside new structure  
6. **Comprehensive**: Covers SQL, HTTP requests, and HTTP responses
7. **Size Aware**: Respects content size limits to prevent memory issues

## Migration Guide

### From Legacy to New Configuration

#### HTTP Server
```go
// Legacy (still works)
s.SetConfigWithMap(g.Map{
    "OtelTraceRequestEnabled": true,
})

// New (recommended)
config := ghttp.NewConfig()
config.Otel.TraceRequestEnabled = true
s.SetConfig(config)
```

#### Database
```go
// Legacy (still works)  
config := gdb.ConfigNode{
    OtelTraceSQLEnabled: true,
}

// New (recommended)
config := gdb.ConfigNode{
    Otel: otel.Config{
        TraceSQLEnabled: true,
    },
}
```

The new configuration provides better organization and allows for future OTEL features to be grouped logically while maintaining full backward compatibility.