# GoFrame 项目指导说明

## 项目概述

GoFrame (GF) 是一个模块化、高性能、企业级的 Golang 基础开发框架。

## 项目目录结构

```shell
gf
├── cmd                        // 命令行工具
│   ├── gf                    // GF CLI工具
│   └── internal             // 内部命令实现
├── container                 // 容器类型
│   ├── garray               // 数组
│   ├── glist                // 列表
│   ├── gmap                 // 映射
│   ├── gpool               // 对象池
│   ├── gqueue              // 队列
│   ├── gset                // 集合
│   ├── gstack              // 栈
│   ├── gtree               // 树
│   └── gtype               // 并发安全基本类型
├── contrib                   // 第三方贡献组件库
│   ├── drivers              // 数据库驱动
│   │   ├── mysql           // MySQL驱动适配
│   │   ├── pgsql           // PostgreSQL驱动适配
│   │   ├── sqlite          // SQLite驱动适配
│   │   ├── mssql           // SQL Server驱动适配
│   │   └── oracle          // Oracle驱动适配
│   ├── registry            // 注册中心实现
│   │   ├── consul         // Consul服务注册
│   │   ├── etcd           // ETCD服务注册
│   │   ├── file           // 文件系统服务注册
│   │   ├── nacos          // Nacos服务注册
│   │   ├── polaris        // Polaris服务注册
│   │   └── zookeeper      // ZooKeeper服务注册
│   ├── config              // 配置中心适配
│   │   ├── kubecm         // K8s ConfigMap配置中心
│   │   └── polaris        // Polaris配置中心
│   └── trace               // 链路追踪实现
├── crypto                   // 加密解密
│   ├── gaes                // AES加密
│   ├── gdes                // DES加密
│   ├── gmd5                // MD5哈希
│   └── gsha1               // SHA1哈希
├── database                 // 数据库功能
│   ├── gdb                 // 数据库ORM
│   ├── gredis              // Redis客户端
│   └── gkafka              // Kafka客户端
├── debug                    // 调试相关功能
│   ├── gdebug              // 调试工具包
│   └── checksum            // 校验和工具
├── encoding                 // 编码解码功能
│   ├── gbase64             // Base64编解码
│   ├── gbinary             // 二进制编码
│   ├── gcharset            // 字符集转换
│   ├── gcompress           // 压缩解压
│   ├── ghash               // 哈希编码
│   ├── ghtml               // HTML编码
│   ├── gjson               // JSON编解码
│   ├── gtoml               // TOML编解码
│   ├── gurl                // URL编解码
│   ├── gxml                // XML编解码
│   └── gyaml               // YAML编解码
├── errors                   // 错误处理
│   ├── gerror              // 错误处理包
│   └── gcode               // 错误码管理
├── example                  // 示例代码
│   ├── http                // HTTP示例
│   ├── database            // 数据库示例
│   └── other               // 其他示例
├── frame                    // 框架核心组件 
│   ├── g                   // 核心包
│   └── gins                // 实例管理
├── i18n                     // 国际化支持
│   └── gi18n               // 多语言管理
├── internal                 // 内部实现
│   ├── command             // 命令行工具实现
│   ├── filepool            // 文件池
│   ├── intlog              // 内部日志
│   └── utils               // 内部工具函数
├── net                      // 网络功能
│   ├── ghttp               // HTTP客户端/服务端
│   ├── gtcp                // TCP组件
│   ├── gudp                // UDP组件
│   ├── goai                // OpenAPI接口
│   ├── gsvc                // 服务注册发现
│   └── gsel                // 负载均衡
├── os                       // 系统功能
│   ├── gcache              // 缓存管理
│   ├── gcfg                // 配置管理
│   ├── gcmd                // 命令行
│   ├── genv                // 环境变量
│   ├── gfile               // 文件操作
│   ├── glog                // 日志管理
│   ├── gmlock              // 内存锁
│   ├── gproc               // 进程管理
│   ├── gres                // 资源管理
│   ├── gsession            // 会话管理
│   ├── gtime               // 时间管理
│   ├── gtimer              // 定时器
│   └── gview               // 模板引擎
├── test                     // 测试相关
│   ├── gtest               // 测试框架
├── text                     // 文本处理
│   ├── gstr                // 字符串处理
│   └── gregex              // 正则表达式
└── util                     // 实用工具
    ├── gconv               // 类型转换
    ├── gmode               // 运行模式
    ├── gutil               // 工具函数
    └── gvalid              // 数据校验
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
