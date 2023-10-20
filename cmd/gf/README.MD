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

   |     DB     | builtin support |                                                                             remarks                                                                              |
   |:----------:|:---------------:|:----------------------------------------------------------------------------------------------------------------------------------------------------------------:|
   |   mysql    |       yes       |                                                                                -                                                                                 |
   |  mariadb   |       yes       |                                                                                -                                                                                 |
   |    tidb    |       yes       |                                                                                -                                                                                 |
   |   mssql    |       yes       |                                                                                -                                                                                 |
   |   oracle   |       yes       |                                                                                -                                                                                 |
   |   pgsql    |       yes       |                                                                                -                                                                                 |
   |   sqlite   |       yes       |                                                                                -                                                                                 |
   | sqlitecgo  |       no        | to support sqlite database on 32bit architecture systems, manually add package import to the [source codes](./internal/cmd/cmd_gen_dao.go) and do the building.  |
   | clickhouse |       no        |                              manually add package import to the [source codes](./internal/cmd/cmd_gen_dao.go) and do the building.                               |
   |     dm     |       no        |                              manually add package import to the [source codes](./internal/cmd/cmd_gen_dao.go) and do the building.                               |

## 2) Manually Install

```shell
go install github.com/gogf/gf/cmd/gf/v2@latest # latest version
go install github.com/gogf/gf/cmd/gf/v2@v2.5.5 # certain version(should be >= v2.5.5)
```

## 2. Commands

```html
$ gf
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
