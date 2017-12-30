Go Frame - 轻量级的Go语言开发框架。

## 安装
```
go get -u gitee.com/johng/gf
````

## 使用
```go
import "gitee.com/johng/gf/g/xxx/xxx"
```

## 说明
    .
    ├── g                【框架目录】
    │   ├── container        【常用数据结构】
    │   │   ├── gbtree           B+树
    │   │   ├── glist            并发安全的双向链表
    │   │   ├── gmap             并发安全的哈希表
    │   │   └── gset             并发安全的集合
    │   │   
    │   ├── database         【数据库操作】
    │   │   └── gdb              通用关系型数据库操作封装(目前仅支持MySQL、PostgreSQL)
    │   │   
    │   ├── encoding         【数据编码】
    │   │   ├── gbase64          BASE64
    │   │   ├── gbinary          二进制操作
    │   │   ├── gcompress        数据压缩
    │   │   ├── gcrc32           CRC32
    │   │   ├── ghash            常用哈希函数
    │   │   ├── ghtml            HTML编码
    │   │   ├── gjson            JSON
    │   │   ├── gmd5             MD5
    │   │   ├── gsha1            SHA1
    │   │   └── gurl             URL
    │   │   
    │   ├── frame            【开发框架】
    │   │   ├── gcfg             通用配置管理
    │   │   ├── gins             单例对象管理
    │   │   └── gmvc             MVC模式封装基类
    │   │   
    │   ├── net              【网络通信】
    │   │   ├── ghttp            HTTP客户端及服务端
    │   │   ├── gipv4            IPv4操作
    │   │   ├── gipv6            IPv6操作
    │   │   ├── grouter          路由管理
    │   │   ├── gscanner         端口扫描
    │   │   ├── gsession         SESSION管理
    │   │   ├── gsmtp            SMTP邮件管理
    │   │   ├── gtcp             TCP服务端
    │   │   └── gudp             UDP服务端
    │   │   
    │   ├── os               【系统管理】
    │   │   ├── gcache           缓存管理
    │   │   ├── gcmd             命令行控制
    │   │   ├── genv             环境变量
    │   │   ├── gfile            文件管理
    │   │   ├── gfilepool        文件指针池
    │   │   ├── gfilespace       文件碎片管理
    │   │   ├── glog             日志管理
    │   │   ├── gmmap            MMAP文件内存映射
    │   │   ├── gtime            时间管理
    │   │   └── gview            视图管理(通用的模板管理对象，独立于MVC模块)
    │   │   
    │   └── util             【常用工具】
    │       ├── grand            随机数
    │       ├── gregx            正则表达式
    │       ├── gvalid           表单校验
    │       └── gutil            其他工具
    │   
    ├── geg              【框架示例】
    ├── vendor           【第三方包】
    └── version.go       【版本信息】
    
## 配置
1. 默认读取文件运行目录的config.json文件；
1. 支持多配置文件管理，默认配置文件目录为当前运行目录；
1. config.json配置文件说明：
    ```json
    {
        "viewpath"   : "视图文件存放路径",
        "database"   : {
            "数据库分组名称" : [
                {
                    "host"     : "地址",
                    "port"     : "端口",
                    "user"     : "账号",
                    "pass"     : "密码",
                    "name"     : "数据库名称",
                    "type"     : "数据库类型(mysql|pgsql)",
                    "role"     : "数据库角色(master|slave)",
                    "charset"  : "字符编码",
                    "priority" : "数据库优先级(负载均衡处理)"
                }
            ]
        }
    }
    ```
1. config.json示例文件：
    ```json
    {
        "viewpath"   : "/home/john/Workspace/Go/gitee.com/johng/gf/geg/frame/mvc/view",
        "database"   : {
            "default" : [
                {
                    "host"     : "127.0.0.1",
                    "port"     : "3306",
                    "user"     : "root",
                    "pass"     : "123456",
                    "name"     : "test",
                    "type"     : "mysql",
                    "role"     : "master",
                    "charset"  : "utf8",
                    "priority" : "100"
                },
                {
                    "host"     : "127.0.0.1",
                    "port"     : "3306",
                    "user"     : "root",
                    "pass"     : "123456",
                    "name"     : "test",
                    "type"     : "mysql",
                    "role"     : "slave",
                    "charset"  : "utf8",
                    "priority" : "100"
                }
            ],
            "ucenter" : [
                {
                    "host"     : "127.0.0.1",
                    "port"     : "3306",
                    "user"     : "root",
                    "pass"     : "123456",
                    "name"     : "ucenter",
                    "type"     : "mysql",
                    "role"     : "master",
                    "charset"  : "utf8",
                    "priority" : "100"
                }
            ]
        }
    }
    ```

