# TPL - Go Template CLI

基于 GoFrame 的项目模板生成工具，通过 `go get` 下载远程 Go 模块作为模板生成新项目。

## 特性

- 支持任意 Go 模块作为项目模板
- 支持嵌套 Go 模块路径（如 `github.com/gogf/gf/cmd/gf/v2`）
- 支持 Git 仓库子目录作为模板（通过 sparse checkout）
- 自动替换 import 路径（基于 AST 解析）
- 交互式模板选择和版本选择
- 自动清理 `.git`、`go.work` 等文件
- 支持自定义 go.mod 模块路径
- 支持依赖升级到最新版本

## 安装

```bash
go build -o tpl.exe .
```

## 快速开始

```bash
# 交互式模式（推荐）
tpl init

# 快速创建单体项目
tpl init github.com/gogf/template-single my-project

# 快速创建大仓项目
tpl init github.com/gogf/template-mono my-mono
```

## 使用方式

### 基本用法

```bash
tpl init <repo> [project-name] [options]
```

### 交互式模式

不带任何参数运行 `tpl init`，将进入交互式模式：

```bash
请选择项目模板:
--------------------------------------------------
  [1] template-single - 单体项目模板
  [2] template-mono - 大仓项目模板
--------------------------------------------------
选择模板 [1-2]: 1
已选择: github.com/gogf/template-single

请输入项目名称: my-project
请输入 Go 模块路径 (留空则使用项目名 "my-project"): github.com/myorg/my-project
是否升级依赖到最新版本 (go get -u)? [y/N]: n
```

### 命令行示例

```bash
# 使用官方模板
tpl init github.com/gogf/template-single my-project       # 单体项目模板
tpl init github.com/gogf/template-mono my-mono            # 大仓项目模板

# 嵌套 Go 模块（自动识别）
tpl init github.com/gogf/gf/cmd/gf/v2 my-gf               # 嵌套模块路径

# 版本控制
tpl init github.com/gogf/gf/cmd/gf/v2 my-gf -s            # 交互式选择版本
tpl init github.com/gogf/gf/cmd/gf/v2@v2.8.0 my-gf        # 指定版本

# 自定义模块路径
tpl init github.com/gogf/template-single my-project -m github.com/myorg/myproject

# Git 子目录（通过 sparse checkout）
tpl init github.com/gogf/examples/httpserver/jwt my-jwt

# 升级依赖到最新
tpl init github.com/gogf/template-single my-project -u
```

## 参数说明

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `repo` | - | 远程仓库地址（位置参数） | 必填 |
| `name` | - | 项目名称（位置参数） | 仓库名 |
| `--select` | `-s` | 启用交互式版本选择 | false |
| `--mod` | `-m` | 自定义 go.mod 模块路径 | 项目名 |
| `--upgrade` | `-u` | 升级依赖到最新版本 | false |

## 支持的模板来源

### 1. 标准 Go 模块

任何可通过 `go get` 下载的 Go 模块都可以作为模板：

```bash
tpl init github.com/gogf/template-single my-project
tpl init github.com/your-org/your-template my-project
```

### 2. 嵌套 Go 模块

支持带有 `/vN` 版本后缀或嵌套在 `cmd/`、`contrib/` 等目录下的模块：

```bash
tpl init github.com/gogf/gf/cmd/gf/v2 my-gf
tpl init github.com/gogf/gf/contrib/drivers/mysql/v2 my-mysql
```

### 3. Git 子目录

对于不是独立 Go 模块的子目录，工具会自动使用 Git sparse checkout：

```bash
tpl init github.com/gogf/examples/httpserver/jwt my-jwt
tpl init github.com/gogf/examples/websocket/chat my-chat
```

## 工作流程

```bash
┌─────────────────┐
│  1. 环境检查     │  验证 Go/Git 安装和配置
└────────┬────────┘
         ▼
┌─────────────────┐
│  2. 版本获取     │  获取可用版本列表（支持交互式选择）
└────────┬────────┘
         ▼
┌─────────────────┐
│  3. 下载模板     │  go get 或 git sparse checkout
└────────┬────────┘
         ▼
┌─────────────────┐
│  4. 生成项目     │  复制文件、清理、替换 import 路径
└────────┬────────┘
         ▼
┌─────────────────┐
│  5. 依赖处理     │  go mod tidy（可选 go get -u）
└─────────────────┘
```

### 详细步骤

1. **环境检查** - 验证 Go 环境（版本、GOPATH、GOPROXY 等）
2. **版本获取** - 通过 `go list -m -versions` 获取可用版本
3. **下载模板** - 使用 `go get` 下载到本地模块缓存
4. **生成项目**
   - 复制模板文件到目标目录
   - 清理 `.git` 目录
   - 清理 `go.work` 和 `go.work.sum`
   - 更新 `go.mod` 中的模块名
   - 使用 AST 替换所有 Go 文件中的 import 路径
5. **依赖处理**
   - 默认执行 `go mod tidy` 整理依赖
   - 使用 `-u` 参数时执行 `go get -u ./...` 升级依赖

## 预置模板

| 模板 | 地址 | 说明 |
|------|------|------|
| template-single | `github.com/gogf/template-single` | GoFrame 单体项目模板 |
| template-mono | `github.com/gogf/template-mono` | GoFrame 大仓项目模板 |

## 环境要求

- Go 1.18+
- Git（仅子目录模板需要）
- GoFrame v2（工具本身依赖）

## 常见问题

### Q: 如何使用私有仓库作为模板？

确保配置了正确的 GOPRIVATE 和认证信息：

```bash
go env -w GOPRIVATE=github.com/your-org
```

### Q: 下载速度慢怎么办？

配置 Go 代理：

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### Q: 如何查看可用版本？

使用 `-s` 参数启用交互式版本选择：

```bash
tpl init github.com/gogf/template-single my-project -s
```

## License

MIT
