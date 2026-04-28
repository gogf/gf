# GoFrame Config Editor

A web-based visual configuration editor for GoFrame projects. It reads your `config.yaml` and provides an interactive UI to view, edit, and save configuration fields with type-aware inputs, validation, and i18n support.

## Quick Start

```bash
gf config                       # Start on port 8888, auto-detect config file
gf config -p 9000               # Use a custom port
gf config -f manifest/config/config.yaml  # Specify config file path
```

The browser opens automatically at `http://127.0.0.1:<port>`.

## Features

### Supported Modules

| Module | Config Node | Description |
|--------|-------------|-------------|
| Server | `server` | HTTP server settings (address, timeouts, TLS, sessions, logging, PProf) |
| Database | `database` | Database connections (host, port, credentials, pool, timeouts) |
| Redis | `redis` | Redis connections (address, auth, pool, sentinel, cluster) |
| Logger | `logger` | Logging configuration (level, rotation, output) |
| Viewer | `viewer` | Template engine settings (paths, delimiters, auto-encode) |

### UI Features

- **Type-aware inputs**: bool fields get toggle switches, duration fields get text input with placeholder hints, map/slice fields get JSON editors
- **Default value display**: each field shows its default value from struct tags
- **Validation**: fields with `v:"required"` tags are validated on blur
- **Modified tracking**: changed fields are marked with a blue indicator bar
- **Group collapse**: fields are organized by logical groups (Basic, Connection, Pool, etc.)
- **Search**: search fields by name, key, description, or type (supports Chinese)
- **i18n**: switch between English and Chinese field descriptions
- **Export format**: save as YAML, TOML, or JSON
- **Keyboard shortcut**: `Ctrl/Cmd + S` to save

### Config File Detection

When no `-f` flag is provided, the editor searches these paths in order:

1. `config.yaml` / `config.yml` / `config.toml` / `config.json`
2. `config/config.yaml` (and variants)
3. `manifest/config/config.yaml` (and variants)
4. `app.yaml` / `app.yml`

### Nested Config Support

GoFrame stores database and redis configs under group names:

```yaml
database:
  default:
    host: 127.0.0.1
    port: 3306
redis:
  default:
    address: 127.0.0.1:6379
```

The editor correctly reads and writes these nested structures.

## Architecture

```
cmd/gf/internal/cmd/
├── cmd_config.go                          # CLI command + REST API handlers
├── resources/
│   ├── templates/index.html               # Vue 3 + Tailwind CSS SPA
│   ├── static/vue.global.prod.js          # Vue 3 runtime
│   ├── static/tailwind.min.css            # Tailwind CSS
│   └── i18n/{en,zh-CN}.yaml              # Field descriptions
os/gcfg/
├── gcfg_schema.go                         # Schema registry (FieldSchema, ModuleSchema, SchemaRegistry)
└── gcfg_z_unit_schema_test.go             # Unit tests
```

### API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/schemas` | Returns all registered module schemas (fields, types, defaults, rules) |
| GET | `/api/config` | Returns current config values from file |
| POST | `/api/config/validate` | Validates config values against schema rules |
| POST | `/api/config/save` | Saves config to file (preserves YAML comments) |
| GET | `/api/i18n/:lang` | Returns i18n translations for the given language |

### Struct Tags

Configuration field metadata is extracted from struct tags:

| Tag | Purpose | Example |
|-----|---------|---------|
| `json` | YAML/JSON key | `json:"address"` |
| `d` | Default value | `d:":0"` |
| `v` | Validation rule (gvalid) | `v:"required"` |
| `dc` | Description + i18n key | `dc:"Server address\|i18n:config.server.address"` |

## Development

### Building

```bash
go build ./cmd/gf/...
```

### Testing

```bash
go test -count=1 -v ./os/gcfg/... -run TestSchema
```
