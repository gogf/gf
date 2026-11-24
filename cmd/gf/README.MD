English | [简体中文](README.zh_CN.MD)

# gf

`gf` is a powerful CLI tool for building [GoFrame](https://goframe.org) application with convenience.

## 1. Install

## 1) PreCompiled Binary

You can also install `gf` tool using pre-built binaries: <https://github.com/gogf/gf/releases>

1. `Mac` & `Linux`

   ```shell
    wget -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH) && chmod +x gf && ./gf install -y && rm ./gf
   ```

   > If you're using `zsh`, you might need rename your alias by command `alias gf=gf` to resolve the conflicts between `gf` and `git fetch`.

2. `Windows`
   Manually download, execute in command line it and then follow the instruction.

3. Database support

   | DB | builtin support | remarks |
   | --- | --- | --- |
   | `mysql` | yes | - |
   | `mariadb` | yes | - |
   | `tidb` | yes | - |
   | `mssql` | yes | - |
   | `oracle` | yes | - |
   | `pgsql` | yes | - |
   | `sqlite` | yes | - |
   | `sqlitecgo` | no | 32-bit `sqlite`: import in `cmd_gen_dao.go` and rebuild. |
   | `clickhouse` | yes | - |
   | `dm` | no | Import in `cmd_gen_dao.go` and rebuild. |

## 2) Manually Install

```shell
go install github.com/gogf/gf/cmd/gf/v2@latest # latest version
go install github.com/gogf/gf/cmd/gf/v2@v2.5.5 # certain version(should be >= v2.5.5)
```

## 2. Commands

The `gf` commands are grouped into three categories: project and runtime, code generation, and environment and helpers.

### 2.1 Project and runtime

| Command | Description |
| --- | --- |
| `gf init` | Create and initialize an empty `GoFrame` project |
| `gf run` | Run Go code with hot-compiled-like behavior, rebuilding and restarting on change |
| `gf build` | Cross-build the project for multiple platforms and architectures |
| `gf pack` | Pack any file or directory into a resource file or Go source file |
| `gf docker` | Build a Docker image for the current project |
| `gf up` | Upgrade the `GoFrame` version or tool version in the current project |

### 2.2 Code generation

| Command | Description |
| --- | --- |
| `gf gen dao` | Generate `dao`, `do` and `entity` code from database tables |
| `gf gen ctrl` | Parse `api` definitions and generate controller and SDK code |
| `gf gen service` | Parse structs and methods in packages and generate `service` interface files |
| `gf gen pb` | Parse `proto` files and generate protobuf code |
| `gf gen pbentity` | Generate `entity` message files in protobuf3 format |
| `gf gen enums` | Parse Go files in the project and generate enum definition files |

### 2.3 Environment and helpers

| Command | Description |
| --- | --- |
| `gf env` | Show current Golang environment variables |
| `gf fix` | Auto-fix code after upgrading to a new `GoFrame` version |
| `gf tpl` | Template parsing and building |
| `gf doc` | Download `https://pages.goframe.org/` to run locally |
| `gf install` | Install the `gf` binary to the system (may require admin permission) |
| `gf version` | Show version information of the current binary |
| `gf agents` | Link local AI coding-agent paths to `.agents` and `AGENTS.md` |

### 2.4 Full help output

```shell
$ gf -h
USAGE
    gf COMMAND [OPTION]

COMMAND
    up         upgrade GoFrame version/tool to latest one in current project
    env        show current Golang environment variables
    fix        auto fixing codes after upgrading to new GoFrame version
    run        running go codes with hot-compiled-like feature
    gen        automatically generate go files for dao/do/entity/pb/pbentity
    tpl        template parsing and building commands
    init       create and initialize an empty GoFrame project
    pack       packing any file/directory to a resource file, or a go file
    build      cross-building go project for lots of platforms
    docker     build docker image for current GoFrame project
    install    install gf binary to system (might need root/admin permission)
    version    show version information of current binary
    doc        download https://pages.goframe.org/ to run locally
    agents     link local AI coding-agent paths to .agents and AGENTS.md

OPTION
    -y, --yes       all yes for all command without prompt ask
    -v, --version   show version information of current binary
    -d, --debug     show internal detailed debugging information
    -h, --help      more information about this command

ADDITIONAL
    Use "gf COMMAND -h" for details about a command.
```

## 3. FAQ

### 1). Command `gf run` returns `pipe: too many open files`

Please use `ulimit -n 65535` to enlarge your system configuration for max open files for current terminal shell session, and then `gf run`.
