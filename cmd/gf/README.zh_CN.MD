[English](README.MD) | 简体中文

# gf

`gf` 是一个强大的 CLI 工具，用于便捷地构建 [GoFrame](https://goframe.org) 应用程序。

## 1. 安装

## 1) 预编译二进制文件

您也可以使用预构建的二进制文件安装 `gf` 工具：<https://github.com/gogf/gf/releases>

1. `Mac` & `Linux`

   ```shell
    wget -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(go env GOOS)_$(go env GOARCH) && chmod +x gf && ./gf install -y && rm ./gf
   ```

   > 如果您使用 `zsh`，您可能需要通过命令 `alias gf=gf` 重命名别名以解决 `gf` 和 `git fetch` 之间的冲突。

2. `Windows`
   手动下载，在命令行中执行，然后按照说明操作。

3. 数据库支持

   |     数据库     | 内置支持 |                                                                             说明                                                                             |
   | :--------: | :-------------: | :-------------------------------------------------------------------------------------------------------------------------------------------------------------: |
   |   mysql    |       是       |                                                                                -                                                                                |
   |  mariadb   |       是       |                                                                                -                                                                                |
   |    tidb    |       是       |                                                                                -                                                                                |
   |   mssql    |       是       |                                                                                -                                                                                |
   |   oracle   |       是       |                                                                                -                                                                                |
   |   pgsql    |       是       |                                                                                -                                                                                |
   |   sqlite   |       是       |                                                                                -                                                                                |
   | sqlitecgo  |       否        | 要在 32 位架构系统上支持 sqlite 数据库，请手动向[源代码](./internal/cmd/cmd_gen_dao.go)添加包导入并进行构建。 |
   | clickhouse |       是        |                                                                                -                                                                                |
   |     dm     |       否        |                              手动向[源代码](./internal/cmd/cmd_gen_dao.go)添加包导入并进行构建。                              |

## 2) 手动安装

```shell
go install github.com/gogf/gf/cmd/gf/v2@latest # 最新版本
go install github.com/gogf/gf/cmd/gf/v2@v2.5.5 # 特定版本(应该 >= v2.5.5)
```

## 2. 命令

```shell
$ gf -h
用法
    gf 命令 [选项]

命令
    up         升级项目中的 GoFrame 版本/工具到最新版本
    env        显示当前 Golang 环境变量
    fix        升级到新 GoFrame 版本后自动修复代码
    run        运行 go 代码，具有热编译功能
    gen        自动生成 dao/do/entity/pb/pbentity 的 go 文件
    tpl        模板解析和构建命令
    init       创建并初始化一个空的 GoFrame 项目
    pack       将任何文件/目录打包到资源文件或 go 文件
    build      为多个平台交叉编译 go 项目
    docker     为当前 GoFrame 项目构建 docker 镜像
    install    将 gf 二进制文件安装到系统(可能需要 root/admin 权限)
    version    显示当前二进制文件的版本信息
    doc        下载 https://pages.goframe.org/ 本地运行

选项
    -y, --yes       对所有命令都使用 yes，不再提示
    -v, --version   显示当前二进制文件的版本信息
    -d, --debug     显示内部详细的调试信息
    -h, --help      显示此命令的更多信息

附加信息
    使用 "gf 命令 -h" 获取有关命令的详细信息。
```

## 3. 常见问题

### 1). 命令 `gf run` 返回 `pipe: too many open files`

请使用 `ulimit -n 65535` 扩大系统配置以增加当前终端 shell 会话的最大打开文件数，然后再运行 `gf run`。
