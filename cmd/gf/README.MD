# gf

`gf` is a powerful CLI tool for building [GoFrame](https://goframe.org) application with convenience.

# Install

> You might need setting the goproxy to make through building.
> Please make sure your Golang version > v1.15.

1. Latest version
    ```
    go install github.com/gogf/gf/cmd/gf/v2
    ```

2. Check installation
   ```
   gf -v
   ```

## Commands
```html
$ gf
USAGE
    gf COMMAND [OPTION]

COMMAND
   env        show current Golang environment variables
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

# FAQ

1. Command `gf run` returns `pipe: too many open files`

   Please use `ulimit -n 65535` to enlarge your system configuration for max open files for current terminal shell session, and then `gf run`.







