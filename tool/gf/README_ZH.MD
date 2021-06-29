# GF-CLI
[English](README.MD) | 简体中文

`gf`是一款强大的[GoFrame](https://goframe.org)开发工具链，使得我们开发基于`GoFrame`框架的项目更加便捷。

## 1. 安装

1. 最新版本
    ```
    go install github.com/gogf/gf/tool/gf@latest
    ```
1. 指定版本
    ```
    go install github.com/gogf/gf/tool/gf@v1.16.0
    ```
> 注意：在`gf gen`命令中，由于`sqlite`和`oracle`数据库需要`cgo`和`gcc`支持，因此预编译的二进制中不提供对这两个数据库的支持，您需要手动修改源码，去掉对应源码文件中指定数据库类型的`import`注释后手动编译支持。

## 2. 命令
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
    update     update current gf binary to latest one (might need root/admin permission)
    install    install gf binary to system (might need root/admin permission)
    version    show current binary version info

OPTION
    -y         all yes for all command without prompt ask
    -?,-h      show this help or detail for specified command
    -v,-i      show version information

ADDITIONAL
    Use 'gf help COMMAND' or 'gf COMMAND -h' for detail about a command, which has '...'
    in the tail of their comments.
```

## 3. 文档

完善详尽的中文文档请参考`GoFrame`官网板块：[开发工具](https://itician.org/pages/viewpage.action?pageId=1114260)

## 4. FAQ

### 1). `gf run` 命令报错 `pipe: too many open files`

请执行`ulimit -n 65535`命令扩展您当前终端会话支持的最大文件打开数，随后再执行`gf run`。需要注意的是该命令仅对当前终端会话有效。







