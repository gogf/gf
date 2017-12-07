Go Frame - 轻量级的Go语言开发框架。

# 安装
```
go get -u gitee.com/johng/gf
````

# 使用
```go
import "gitee.com/johng/gf/g/xxx"
```

# 说明
    .
    ├── g                【框架目录】
    │   ├── container        【常用数据结构容器】
    │   │   ├── gbtree           B+树
    │   │   ├── glist            并发安全的双向链表
    │   │   ├── gmap             并发安全的哈希表
    │   │   └── gset             并发安全的集合
    │   │   
    │   ├── database         【数据库操作】
    │   │   └── gdb              MySQL、PostgreSQL的操作封装
    │   │   
    │   ├── encoding         【数据编码】
    │   │   ├── gbase64          BASE64
    │   │   ├── gbinary          二进制操作
    │   │   ├── gcompress        数据压缩
    │   │   ├── gcrc32           CRC32
    │   │   ├── ghash            常用哈希函数
    │   │   ├── ghtml            HTML相关编码
    │   │   ├── gjson            JSON
    │   │   ├── gmd5             MD5
    │   │   ├── gsha1            SHA1
    │   │   └── gurl             URL
    │   │   
    │   ├── net              【网络通信】
    │   │   ├── ghttp            HTTP客户端及服务端
    │   │   ├── gip              IP操作
    │   │   ├── gscanner         端口扫描
    │   │   ├── gsmtp            SMTP邮件管理
    │   │   ├── gtcp             TCP客户端及服务端
    │   │   └── gudp             UDP客户端及服务端
    │   │   
    │   ├── os               【系统管理】
    │   │   ├── gcache           本地缓存管理
    │   │   ├── gconsole         命令行控制
    │   │   ├── gfile            文件管理
    │   │   ├── gfilepool        文件指针池
    │   │   ├── gfilespace       文件碎片管理
    │   │   ├── glog             日志管理
    │   │   ├── gmmap            MMAP文件内存映射
    │   │   └── gtime            时间管理
    │   │   
    │   └── util             【常用工具】
    │       ├── gpage            分页
    │       ├── grand            随机数
    │       ├── gregx            正则表达式
    │       └── gutil            其他工具
    │   
    ├── geg              【框架示例】
    ├── vendor           【第三方包】
    └── version.go       【版本信息】