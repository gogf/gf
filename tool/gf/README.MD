# GoFrame CLI Tool

Powerful CLI tool for building [GoFrame](https://goframe.org) application with convenience.

## 1. Install

> You might need setting the goproxy to make through building.

1. Latest version
   > Go version >= v1.16
    ```
    go install github.com/gogf/gf/tool/gf@latest
    ```
   > Go version < v1.16
    ```
    go install github.com/gogf/gf/tool/gf
    ```
1. Specified version
   > Go version >= v1.16
    ```
    go install github.com/gogf/gf/tool/gf@v1.16.0
    ```
1. Database `sqlite` and `oracle` are not support in `gf gen` command in default as it needs `cgo` and `gcc`, you can manually make some changes to the source codes and do the building.

## 2. Commands
```html
$ gf
USAGE
    gf COMMAND [ARGUMENT] [OPTION]

COMMAND
    env        show current Golang environment variables
    get        install or update GF to system in default...
    gen        automatically generate go files for ORM models...
    mod        extra features for go modules...
    run        running go codes with hot-compiled-like feature...
    init       initialize an empty GF project at current working directory...
    help       show more information about a specified command
    pack       packing any file/directory to a resource file, or a go file...
    build      cross-building go project for lots of platforms...
    docker     create a docker image for current GF project...
    swagger    swagger feature for current project...
    update     update current gf binary to latest one
    version    show current binary version info

OPTION
    -y         all yes for all command without prompt ask
    -?,-h      show this help or detail for specified command
    -v,-i      show version information

ADDITIONAL
    Use 'gf help COMMAND' or 'gf COMMAND -h' for detail about a command, which has '...'
    in the tail of their comments.
```

## 3. FAQ

### 1). Command `gf run` returns `pipe: too many open files`

Please use `ulimit -n 65535` to enlarge your system configuration for max open files for current terminal shell session, and then `gf run`.







