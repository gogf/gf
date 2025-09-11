# OpenTelemetry V2.8 Improvements Example

This example demonstrates the new configurable OpenTelemetry tracing features for SQL, HTTP requests, and HTTP responses.

## HTTP Server Configuration

### Enable Request Parameter Tracing
```yaml
server:
  address: ":8080"
  otelTraceRequestEnabled: true  # Enable HTTP request parameter tracing
```

### Enable Response Body Tracing
```yaml
server:
  address: ":8080"
  otelTraceResponseEnabled: true  # Enable HTTP response body tracing
```

### Enable Both Request and Response Tracing
```yaml
server:
  address: ":8080" 
  otelTraceRequestEnabled: true   # Enable HTTP request parameter tracing
  otelTraceResponseEnabled: true  # Enable HTTP response body tracing
```

## Database Configuration

### Enable SQL Statement Tracing
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

### HTTP Server
```go
package main

import (
    "github.com/gogf/gf/v2/frame/g"
    "github.com/gogf/gf/v2/net/ghttp"
)

func main() {
    s := g.Server()
    
    // Enable tracing via configuration map
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

### Database
```go
package main

import (
    "context"
    "github.com/gogf/gf/v2/database/gdb"
    "github.com/gogf/gf/v2/frame/g"
)

func main() {
    // Configure database with SQL tracing enabled
    config := gdb.ConfigNode{
        Type: "mysql",
        Host: "127.0.0.1",
        Port: "3306",
        User: "your_user",
        Pass: "your_password",
        Name: "your_database",
        OtelTraceSQLEnabled: true,  // Enable SQL tracing
    }
    
    db, err := gdb.New(config)
    if err != nil {
        g.Log().Fatal(context.TODO(), err)
    }
    defer db.Close(context.TODO())
    
    // SQL queries will now be traced
    result, err := db.GetOne(context.TODO(), "SELECT * FROM users WHERE id = ?", 1)
    if err != nil {
        g.Log().Error(context.TODO(), err)
        return
    }
    
    g.Log().Info(context.TODO(), "User:", result)
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

1. **Configurable**: All new tracing features are opt-in via configuration
2. **Performance**: Only enabled features add overhead 
3. **Backward Compatible**: Existing tracing continues to work unchanged
4. **Comprehensive**: Covers SQL, HTTP requests, and HTTP responses
5. **Size Aware**: Respects content size limits to prevent memory issues