# GoFrame 项目指导说明

## 项目概述

GoFrame (GF) 是一个模块化、高性能、企业级的 Golang 基础开发框架。

## 目录结构说明

### 主包

```shell
gf
├── container                 // 容器相关包
│   ├── garray               // 数组容器
│   ├── glist                // 链表容器
│   ├── gmap                 // Map容器
│   ├── gpool                // 对象池
│   ├── gqueue               // 队列
│   ├── gring                // 环形缓冲区
│   ├── gset                 // 集合
│   ├── gtree                // 树结构
│   ├── gtype                // 并发安全类型
│   └── gvar                 // 动态变量
├── crypto                    // 加密相关
│   ├── gaes                 // AES加密
│   ├── gcrc32               // CRC32校验
│   ├── gdes                 // DES加密
│   ├── gmd5                 // MD5哈希
│   └── gsha1                // SHA1哈希
├── database                  // 数据库相关
│   ├── gdb                  // 数据库ORM
│   └── gredis               // Redis客户端
├── debug                     // 调试工具
│   └── gdebug              // 调试辅助
├── encoding                  // 编码相关
│   ├── gbase64             // Base64编码
│   ├── gbinary             // 二进制编码
│   ├── gcharset            // 字符集转换
│   ├── gcompress           // 压缩解压
│   ├── ghash               // 哈希算法
│   ├── ghtml               // HTML处理
│   ├── gini                // INI解析
│   ├── gjson               // JSON处理
│   ├── gproperties         // Properties解析
│   ├── gtoml               // TOML解析
│   ├── gurl                // URL处理
│   ├── gxml                // XML处理
│   └── gyaml               // YAML处理
├── errors                    // 错误处理
│   ├── gcode               // 错误码
│   └── gerror              // 错误处理
├── frame                     // 框架核心
│   ├── g                   // 全局对象
│   └── gins                // 依赖注入
├── i18n                      // 国际化
│   └── gi18n               // 国际化支持
├── net                       // 网络相关
│   ├── gclient             // HTTP客户端
│   ├── ghttp               // HTTP服务端
│   ├── gipv4               // IPv4工具
│   ├── gipv6               // IPv6工具
│   ├── goai                // AI工具
│   ├── gsel                // 服务发现
│   ├── gsvc                // 服务治理
│   ├── gtcp                // TCP工具
│   ├── gtrace              // 链路追踪
│   └── gudp                // UDP工具
├── os                        // 系统相关
│   ├── gbuild              // 构建工具
│   ├── gcache              // 缓存管理
│   ├── gcfg                // 配置管理
│   ├── gcmd                // 命令行解析
│   ├── gcron               // 定时任务
│   ├── gctx                // 上下文管理
│   ├── genv                // 环境变量
│   ├── gfile               // 文件操作
│   ├── gfpool              // 文件池
│   ├── gfsnotify           // 文件监控
│   ├── glog                // 日志管理
│   ├── gmetric             // 指标监控
│   ├── gmlock              // 内存锁
│   ├── gmutex              // 互斥锁
│   ├── gproc               // 进程管理
│   ├── gres                // 资源管理
│   ├── grpool              // 协程池
│   ├── gsession            // 会话管理
│   ├── gspath              // 路径处理
│   ├── gstructs            // 结构体工具
│   ├── gtime               // 时间处理
│   ├── gtimer              // 定时器
│   └── gview               // 视图渲染
├── test                      // 测试工具
│   └── gtest               // 测试框架
├── text                      // 文本处理
│   ├── gregex              // 正则表达式
│   └── gstr                // 字符串工具
└── util                      // 工具类
    ├── gconv               // 类型转换
    ├── gmeta               // 元数据处理
    ├── gmode               // 运行模式
    ├── gpage               // 分页工具
    ├── grand               // 随机数
    ├── gtag                // 标签处理
    ├── guid                // UUID生成
    ├── gutil               // 通用工具
    └── gvalid              // 数据校验
```

### cmd 命令行工具

```shell
cmd
├── gf                        // GF CLI主程序
│   ├── gfcmd                // CLI命令入口
│   │   ├── build            // 项目构建 (cmd_build.go)
│   │   ├── run              // 热编译运行 (cmd_run.go) 
│   │   ├── init             // 项目脚手架 (cmd_init.go)
│   │   ├── gen              // 代码生成入口 (cmd_gen.go)
│   │   ├── docker           // 容器化操作 (cmd_docker.go)
│   │   ├── install          // 依赖管理 (cmd_install.go)
│   │   ├── fix              // 代码修复 (cmd_fix.go)
│   │   ├── update           // 框架升级 (cmd_up.go)
│   │   ├── env              // 环境变量管理 (cmd_env.go)
│   │   ├── pack             // 二进制打包 (cmd_pack.go)
│   │   └── doc              // 文档生成 (cmd_doc.go)
│   ├── internal/cmd/         // 命令实现核心
│   │   ├── cmd_build.go     // 构建命令：交叉编译支持/构建参数配置
│   │   ├── cmd_doc.go       // 文档命令：Swagger/API文档自动化生成
│   │   ├── cmd_docker.go    // Docker命令：镜像构建/推送/多阶段编译
│   │   ├── cmd_env.go       // 环境管理：变量查看/设置/环境切换  
│   │   ├── cmd_fix.go       // 代码修复：自动修复常见语法问题
│   │   ├── cmd_gen.go       // 代码生成：统一入口路由
│   │   ├── cmd_gen_ctrl.go  // MVC控制器：RESTful接口生成
│   │   ├── cmd_gen_dao.go   // DAO层：数据库表映射生成
│   │   ├── cmd_gen_enums.go // 枚举代码：自动生成枚举类型和方法
│   │   ├── cmd_gen_pb.go    // Protobuf：协议文件编译生成
│   │   ├── cmd_gen_pbentity.go // Protobuf实体：数据库表到proto转换
│   │   ├── cmd_gen_service.go // 微服务接口：GRPC服务代码生成
│   │   ├── cmd_init.go      // 项目初始化：模块化脚手架生成
│   │   ├── cmd_install.go   // 依赖管理：自动分析并安装go依赖
│   │   ├── cmd_pack.go      // 打包发布：支持二进制/Docker/zip多种格式
│   │   ├── cmd_run.go       // 运行管理：热编译/配置重载/进程监控
│   │   ├── cmd_tpl.go       // 模板管理：自定义代码模板系统
│   │   ├── cmd_up.go        // 框架升级：版本检测与自动更新
│   │   ├── cmd_version.go   // 版本管理：CLI/Golang/框架版本信息
│   │   ├── cmd_z_init_test.go // 初始化测试：脚手架生成验证
│   │   └── cmd_z_unit_*_test.go // 单元测试：各命令功能验证
│   ├── internal/cmd/gen/    // 代码生成模板
│   │   ├── tpl_field.go     // 字段级模板(列映射/类型转换)
│   │   ├── tpl_table.go     // 表级模板(CRUD操作/关系映射)
│   │   ├── tpl_test.go      // 测试用例模板
│   │   ├── tpl_ctrl.go      // 控制器模板(RESTful方法)
│   │   ├── tpl_service.go   // 服务层模板(业务逻辑)
│   │   └── tpl_pbentity.go  // Protobuf实体模板
│   ├── test/                // 测试相关
│   │   ├── cmd_z_unit_build_test.go  // 构建命令单元测试
│   │   ├── cmd_z_unit_gen_dao_test.go // DAO生成测试
│   │   └── testdata/        // 测试用例数据
│   ├── internal             // 内部实现
│   ├── test                 // 测试代码
│   ├── go.mod               // 模块文件
│   ├── go.sum               // 依赖校验
│   ├── go.work              // 工作区配置
│   ├── LICENSE              // 许可证
│   ├── main.go              // 主入口
│   ├── Makefile             // 构建配置
│   └── README.MD            // 说明文档
```

### contrib 组件库

```shell
contrib
├── config                  // 配置中心支持
│   ├── apollo             // Apollo配置中心
│   ├── consul             // Consul配置中心
│   ├── kubecm             // Kubernetes ConfigMap支持
│   ├── nacos              // Nacos配置中心
│   └── polaris            // Polaris配置中心
├── drivers                 // 数据库驱动
│   ├── clickhouse         // ClickHouse驱动
│   ├── dm                 // 达梦数据库驱动
│   ├── mssql              // SQL Server驱动
│   ├── mysql              // MySQL驱动
│   ├── oracle             // Oracle驱动
│   ├── pgsql              // PostgreSQL驱动
│   ├── sqlite             // SQLite驱动
│   └── sqlitecgo          // SQLite CGO驱动
├── metric                  // 指标监控
│   └── otelmetric         // OpenTelemetry指标支持
├── nosql                   // NoSQL支持
│   └── redis              // Redis支持
├── registry                // 服务注册发现
│   ├── consul             // Consul支持
│   ├── etcd               // Etcd支持
│   ├── file               // 文件注册中心
│   ├── nacos              // Nacos支持
│   ├── polaris            // Polaris支持
│   └── zookeeper          // Zookeeper支持
├── rpc                     // RPC支持
│   └── grpcx              // gRPC扩展支持
├── sdk                     // SDK支持
│   └── httpclient         // HTTP客户端SDK
└── trace                   // 链路追踪
    ├── otlpgrpc           // OpenTelemetry gRPC支持
    └── otlphttp           // OpenTelemetry HTTP支持
```

### examples 示例库

```shell
examples
├── balancer                  // 负载均衡示例
│   ├── http                 // HTTP负载均衡
│   └── polaris              // Polaris负载均衡
├── config                    // 配置中心示例
│   ├── apollo               // Apollo配置中心
│   ├── consul               // Consul配置中心
│   ├── kubecm               // Kubernetes ConfigMap
│   ├── nacos                // Nacos配置中心
│   └── polaris              // Polaris配置中心
├── converter                 // 类型转换示例
│   ├── alias-type-convert   // 别名类型转换
│   ├── alias-type-scan      // 别名类型扫描
│   ├── struct-convert       // 结构体转换
│   └── struct-scan          // 结构体扫描
├── database                  // 数据库示例
│   └── mysql                // MySQL数据库
├── httpserver                // HTTP服务示例
│   ├── default-value        // 默认值处理
│   ├── proxy                // 代理服务
│   ├── rate                 // 限流控制
│   ├── response-with-json   // JSON响应
│   ├── serve-file           // 文件服务
│   ├── swagger              // Swagger文档
│   ├── upload-file          // 文件上传
│   └── swagger-set-template // Swagger模板
├── metric                    // 指标监控示例
│   ├── basic                // 基础指标
│   ├── callback             // 回调指标
│   ├── dynamic-attributes   // 动态属性
│   ├── global-attributes    // 全局属性
│   ├── http-client          // HTTP客户端指标
│   ├── http-server          // HTTP服务端指标
│   ├── meter-attributes     // 计量器属性
│   └── prometheus           // Prometheus集成
├── nosql                     // NoSQL示例
│   └── redis                // Redis操作
├── os                        // 系统操作示例
│   ├── cron                 // 定时任务
│   └── log                  // 日志管理
├── pack                      // 打包示例
│   ├── hack                 // 打包工具
│   ├── manifest             // 清单文件
│   ├── packed               // 打包结果
│   └── resource             // 资源文件
├── registry                  // 服务注册发现示例
│   ├── consul               // Consul注册中心
│   ├── etcd                 // Etcd注册中心
│   ├── file                 // 文件注册中心
│   ├── nacos                // Nacos注册中心
│   └── polaris              // Polaris注册中心
├── rpc                       // RPC示例
│   └── grpcx                // gRPC扩展
├── tcp                       // TCP示例
│   └── server               // TCP服务
└── trace                     // 链路追踪示例
    ├── grpc-with-db         // gRPC+数据库
    ├── http                 // HTTP链路
    ├── http-with-db         // HTTP+数据库
    ├── inprocess            // 进程内追踪
    ├── inprocess-grpc       // 进程内gRPC
    ├── otlp                 // OpenTelemetry
    ├── processes            // 进程管理
    └── provider             // 追踪提供者
```

## 编码规范

1. 命名规范

- 包名使用小写
- 结构体、接口名使用大驼峰
- 方法名使用大驼峰
- 变量名使用小驼峰

2. 代码格式

- 使用`gofmt`标准格式化
- 遵循 Go 官方代码规范
- 每个包都应有详细的文档注释

3. 错误处理

- 使用`gerror`包进行错误处理
- 错误信息应该清晰明确

4. 测试规范

- 所有公开接口需要单元测试
- 测试文件命名为`xxx_test.go`
- 基准测试命名为`BenchmarkXxx`

5. 依赖管理

- 使用`go mod`进行依赖管理
- golang 的版本根据 go.mod 文件中的 go 版本进行管理

## 项目特定指南

1. 模块开发

- 遵循模块化设计原则
- 使用依赖注入模式
- 保持向后兼容性

2. 文档编写

- 使用英文编写代码注释
- 中英文文档同步更新
- 示例代码需要可运行

3. 性能考虑

- 注意内存分配
- 避免不必要的类型转换
- 合理使用缓存机制

## 代码生成建议

生成代码时请遵循以下原则：

- 符合 Go 语言惯用法
- 保持代码简洁清晰
- 注重性能和可维护性
- 添加必要的注释说明
