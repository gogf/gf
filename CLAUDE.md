# # CLAUDE.md - AI Assistant Guide for GoFrame (GF)

This file provides guidance for AI assistants working with the GoFrame (GF) codebase.

## Project Overview

GoFrame (GF) is a modular, powerful, high-performance and enterprise-class application development framework for Go. It provides comprehensive functionality for web development, microservices, CLI tools, and enterprise applications.

**Current Version:** v2.9.7

GoFrame aims to make development faster, easier, and more efficient through:
- Modular architecture with loosely-coupled components
- Comprehensive standard library with 60+ packages
- High performance and production-ready design
- Rich ecosystem with tools and extensions

## Repository Structure

```
gf/
├── cmd/                           # Command line tools
│   └── gf/                        # GF CLI development toolkit
├── container/                     # Data containers
│   ├── garray/                    # Dynamic arrays with concurrency safety
│   ├── glist/                     # Concurrent-safe doubly linked list
│   ├── gmap/                      # Concurrent-safe maps with various key types
│   ├── gpool/                     # Object pools for memory reuse
│   ├── gqueue/                    # Concurrent-safe queues
│   ├── gring/                     # Concurrent-safe ring/circular buffer
│   ├── gset/                      # Concurrent-safe sets
│   ├── gtree/                     # Concurrent-safe tree containers
│   ├── gtype/                     # High-performance atomic operations
│   └── gvar/                      # Universal variable container
├── contrib/                       # Third-party integrations
│   ├── config/                    # Configuration providers (Apollo, Consul, etc.)
│   ├── drivers/                   # Database drivers (MySQL, PostgreSQL, etc.)
│   ├── metric/                    # Monitoring integrations (Prometheus, etc.)
│   ├── nosql/                     # NoSQL integrations (Redis, MongoDB, etc.)
│   ├── registry/                  # Service discovery (ETCD, Consul, etc.)
│   ├── rpc/                       # RPC frameworks (gRPC, etc.)
│   ├── sdk/                       # External service SDKs
│   └── trace/                     # Distributed tracing (Jaeger, Zipkin)
├── crypto/                        # Cryptography utilities
│   ├── gaes/                      # AES encryption/decryption
│   ├── gdes/                      # DES encryption/decryption
│   ├── gmd5/                      # MD5 hashing
│   ├── grsa/                      # RSA encryption/decryption
│   └── gsha*/                     # SHA hashing algorithms
├── database/                      # Database operations
│   ├── gdb/                       # ORM and database operations
│   └── gredis/                    # Redis operations
├── encoding/                      # Data encoding/decoding
│   ├── gbase64/                   # Base64 encoding
│   ├── gbinary/                   # Binary operations
│   ├── gcharset/                  # Charset conversion
│   ├── gcompress/                 # Compression (gzip, zlib)
│   ├── gjson/                     # JSON operations
│   ├── gini/                      # INI file operations
│   ├── gtoml/                     # TOML operations
│   ├── gxml/                      # XML operations
│   └── gyaml/                     # YAML operations
├── errors/                        # Error handling
│   ├── gcode/                     # Error codes
│   └── gerror/                    # Error management with stack trace
├── frame/                         # Framework core
│   ├── g/                         # Core entry point and global objects
│   └── gins/                      # Singleton management
├── net/                           # Network operations
│   ├── gclient/                   # HTTP client
│   ├── ghttp/                     # HTTP server and middleware
│   ├── gtcp/                      # TCP operations
│   ├── gudp/                      # UDP operations
│   ├── gipv4/                     # IPv4 utilities
│   ├── gipv6/                     # IPv6 utilities
│   ├── gtrace/                    # HTTP tracing
│   └── gsvc/                      # Service registry and discovery
├── os/                            # Operating system operations
│   ├── gcache/                    # Cache management
│   ├── gcfg/                      # Configuration management
│   ├── gcmd/                      # Command line parsing
│   ├── gcron/                     # Scheduled tasks (cron)
│   ├── gctx/                      # Context management
│   ├── genv/                      # Environment variables
│   ├── gfile/                     # File operations
│   ├── glog/                      # Logging
│   ├── gproc/                     # Process management
│   ├── gres/                      # Resource packaging
│   ├── gsession/                  # Session management
│   ├── gtime/                     # Time operations
│   └── gview/                     # Template engine
├── test/                          # Testing utilities
│   └── gtest/                     # Testing helper functions
├── text/                          # Text processing
│   ├── gregex/                    # Regular expressions
│   └── gstr/                      # String operations
├── util/                          # Utilities
│   ├── gconv/                     # Type conversion
│   ├── gmeta/                     # Metadata management
│   ├── grand/                     # Random number generation
│   ├── gtag/                      # Struct tag parsing
│   ├── guid/                      # UUID generation
│   ├── gutil/                     # Utility functions
│   └── gvalid/                    # Data validation
└── internal/                      # Internal packages
    ├── command/                   # CLI command implementations
    ├── json/                      # JSON operations (internal)
    ├── mutex/                     # Mutex utilities
    └── utils/                     # Internal utilities
```

## Build System

### Prerequisites

- **Go 1.23.0+** (required for latest features)
- **Git** (for dependency management)
- **Make** (for build automation)
- **golangci-lint** (for code quality checks)

### Common Build Commands

```bash
# Get the framework
go get -u github.com/gogf/gf/v2

# Build all packages
go build ./...

# Install dependencies
go mod download

# Tidy module dependencies
go mod tidy
# or: make tidy (executes on all modules)

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run specific test
go test -run TestFunctionName

# Code quality check
make lint
```

### Make Commands

GoFrame provides several make targets for development:

- `make tidy`: Execute `go mod tidy` on all modules
- `make lint`: Run golangci-lint for code quality checks
- `make branch to=vX.Y.Z`: Create fix branches for specific versions

## Architecture & Key Concepts

### Framework Design Principles

1. **Modular Architecture**: Each package is independent and loosely coupled
2. **High Performance**: Focus on memory efficiency and execution speed
3. **Production Ready**: Battle-tested components used in enterprise applications
4. **Developer Friendly**: Intuitive APIs with comprehensive documentation
5. **Backward Compatibility**: Careful version management with semantic versioning

### Core Concepts

- **Singleton Pattern**: Global objects accessible via `g.*` methods
- **Chaining API**: Fluent interfaces for better code readability
- **Concurrency Safety**: Thread-safe operations across container types
- **Context Propagation**: Built-in context support for request tracing
- **Configuration Management**: Unified configuration with multiple providers
- **Error Chain**: Stack trace and error code management

### Key Components

#### Core Framework (`frame/g`)

Entry point providing global access to all framework components:

```go
g.Server()    // HTTP server
g.DB()        // Database
g.Redis()     // Redis client
g.Log()       // Logger
g.Cfg()       // Configuration
g.Validator() // Data validator
```

#### Container Package (`container/*`)
High-performance data structures with concurrency safety:
- `garray`: Dynamic arrays with automatic read/write locks
- `gmap`: Various map types (HashMap, TreeMap, ListMap)
- `gset`: Set implementations (HashSet, TreeSet)
- `gtype`: Atomic operations for basic types

#### HTTP Server (`net/ghttp`)
Full-featured HTTP server with middleware support, routing, and OpenAPI integration.

#### Database ORM (`database/gdb`)
Powerful ORM with support for multiple databases, query builders, and migrations.

## Code Style & Conventions

### Go Standards

GoFrame follows standard Go conventions:
- **Package naming**: lowercase, single word when possible
- **Function naming**: CamelCase for public, camelCase for private
- **Variable naming**: camelCase with meaningful names
- **Constants**: UPPER_CASE_WITH_UNDERSCORES

### GoFrame Specific Conventions

- **Chaining APIs**: Use method chaining for fluent interfaces
- **Error handling**: Always check and handle errors explicitly
- **Context usage**: Pass context.Context as first parameter when applicable
- **Resource cleanup**: Use defer for cleanup operations
- **Nil checks**: Always check for nil before using pointers

### Code Quality Tools

The project uses `golangci-lint` with specific rules:
- **gofmt**: Automatic code formatting
- **govet**: Static analysis
- **ineffassign**: Detect ineffective assignments
- **misspell**: Spell checking in comments and strings
- **gocyclo**: Cyclomatic complexity checking

### Testing Guidelines

- **Test naming**: `TestFunctionName` for unit tests
- **Benchmark naming**: `BenchmarkFunctionName` for benchmarks
- **Example naming**: `ExampleFunctionName` for documentation examples
- **Coverage**: Aim for high test coverage on critical paths
- **Table-driven tests**: Use for testing multiple scenarios

## Development Tasks

### Adding New Packages

1. Create package directory following naming conventions
2. Implement main functionality with proper documentation
3. Add comprehensive tests with examples
4. Update go.mod if external dependencies are needed
5. Add package documentation to README if needed

### Working with HTTP Server

```go
// Basic server setup
s := g.Server()
s.BindHandler("/", func(r *ghttp.Request) {
    r.Response.Write("Hello GoFrame")
})
s.Run()
```

### Database Operations

```go
// Database query example
result, err := g.DB().Table("users").Where("age > ?", 18).All()
if err != nil {
    g.Log().Error(ctx, err)
    return
}
```

### Validation Usage

```go
// Validation example
if err := g.Validator().Data(data).Rules("name@required|length:2,20").Run(ctx); err != nil {
    return err
}
```

## Testing

### Test Structure

- **Unit tests**: `*_test.go` files in the same package
- **Example tests**: `*_example_test.go` for documentation
- **Benchmark tests**: `*_bench_test.go` for performance testing
- **Integration tests**: Separate test packages when needed

### Running Tests

```bash
# All tests
go test ./...

# Specific package
go test ./util/gvalid

# With coverage
go test -cover ./...

# Race condition detection
go test -race ./...

# Verbose output
go test -v ./...
```

## Important Files for Understanding the Codebase

- [`version.go`](version.go) - Current framework version
- [`go.mod`](go.mod) - Module dependencies and Go version
- [`Makefile`](Makefile) - Build automation scripts
- [`frame/g/g.go`](frame/g/g.go) - Core framework entry point
- [`container/garray/garray.go`](container/garray/garray.go) - Container pattern examples
- [`net/ghttp/ghttp.go`](net/ghttp/ghttp.go) - HTTP server implementation
- [`database/gdb/gdb.go`](database/gdb/gdb.go) - Database ORM core
- [`util/gvalid/gvalid.go`](util/gvalid/gvalid.go) - Validation engine

## Key Development Guidelines

### Performance Considerations

1. **Memory allocation**: Prefer sync.Pool for frequently allocated objects
2. **String operations**: Use gstr package for efficient string processing
3. **Concurrency**: Leverage gtype for atomic operations
4. **Caching**: Use gcache for application-level caching

### Security Best Practices

1. **Input validation**: Always validate user input using gvalid
2. **SQL injection**: Use parameterized queries with gdb
3. **XSS prevention**: Proper HTML escaping in templates
4. **Context timeouts**: Set appropriate timeouts for operations

### Error Handling

1. **Error wrapping**: Use gerror.Wrap() to add context
2. **Error codes**: Define error codes using gcode
3. **Logging**: Use glog for structured logging
4. **Stack traces**: gerror automatically includes stack traces

## Contributing Guidelines

### Branch Strategy

- **master**: Main development branch
- **fix/vX.Y.Z**: Bug fix branches for specific versions
- **feature/feature-name**: Feature development branches

### Commit Messages

Follow conventional commit format:
- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation updates
- `style:` formatting changes
- `refactor:` code refactoring
- `test:` test additions or modifications
- `chore:` maintenance tasks

### Pull Request Process

1. Fork the repository
2. Create feature branch from master
3. Make changes with comprehensive tests
4. Run `make lint` and `make tidy`
5. Submit pull request with detailed description
6. Address code review feedback

## Tips for AI Assistants

1. **Package relationships**: Understand the modular architecture - packages are designed to be independent
2. **Performance focus**: GoFrame prioritizes performance - consider memory allocation and concurrency in suggestions
3. **Context usage**: Many functions accept context.Context - include this in API calls when appropriate
4. **Error handling**: Follow Go error handling patterns - don't ignore errors
5. **Testing importance**: Always suggest comprehensive tests for new functionality
6. **Backward compatibility**: Be careful with API changes - GoFrame maintains strong backward compatibility
7. **Documentation**: Code should be well-documented following Go documentation standards
8. **Concurrent safety**: Many GoFrame components are concurrent-safe - mention this when relevant
9. **Configuration**: Use gcfg for configuration management rather than hardcoded values
10. **Logging**: Use glog instead of fmt.Println for production code
