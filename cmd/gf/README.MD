# gf

`gf` is a powerful CLI tool for building [GoFrame](https://goframe.org) application with convenience.


## 1. Install

### 1). Manually Install

> You might need setting the goproxy to make through building.
> Please make sure your Golang version > v1.15.

1. Latest version
    ```
    go install -g github.com/gogf/gf/cmd/gf/v2@latest
    ```

2. Check installation
   ```
   gf -v
   ```
### 2). PreBuilt Binary

You can also install `gf` tool using pre-built binaries: https://github.com/gogf/gf/releases

After downloads, please use `gf_xxx_xxx install` command to install gf binary to system binary path.

1. `Mac`
    ```shell
    wget -O gf https://github.com/gogf/gf/releases/download/v2.0.0-rc2/gf_darwin_amd64 && chmod +x gf && ./gf install
    ```
   > If you're using `zsh`, you might need rename your alias by command `alias gf=gf` to resolve the conflicts between `gf` and `git fetch`.

2. `Linux`
    ```shell
    wget -O gf https://github.com/gogf/gf/releases/download/v2.0.0-rc2/gf_linux_amd64 && chmod +x gf && ./gf install
    ```

3. `Windows`

   Manually download, execute it and then follow the instruction.

4. Database `sqlite` and `oracle` are not support in `gf gen` command in default as it needs `cgo` and `gcc`, you can manually make some changes to the source codes and do the building.

## 2. Commands
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

## 3. FAQ

### 1). Command `gf run` returns `pipe: too many open files`

Please use `ulimit -n 65535` to enlarge your system configuration for max open files for current terminal shell session, and then `gf run`.







