# `v1.10.0` (2019-12-05)

各位`gfer`久等了，较上一次发布时间过去已有两个多月了，这段时间`GF`也在不断地迭代改进，细节比较多，拟了个大概，以下是`release log`。

另外，`GoFrame`也参加了2019最受欢迎中国开源软件评选投票，明天就结束了，欢迎为`GF`投票啊：https://www.oschina.net/project/top_cn_2019 网页可以投一票，微信也可以投一票。

## 新特性

1. `Web Server`新特性：
    - 改进中间件及分组路由实现：https://goframe.org/net/ghttp/router/middleware
    - 增加文件配置管理特性：https://goframe.org/net/ghttp/config
    - 改进参数获取：https://goframe.org/net/ghttp/request
    - 改进文件上传：https://goframe.org/net/ghttp/client/demo/upload
1. `Session`增加内置的多种`Storage`实现：
    - 基本介绍：https://goframe.org/os/gsession/index
    - 文件存储：https://goframe.org/os/gsession/file
    - 内存存储：https://goframe.org/os/gsession/memory
    - `Redis`存储：https://goframe.org/os/gsession/redis
1. 增加日志组件单例对象，并优化配置管理：
    - https://goframe.org/frame/g/index
    - https://goframe.org/os/glog/config
1. 常用的`container`容器增加`JSON`数据格式的`Marshal`/`UnMarshal`接口实现：
    - https://goframe.org/container/gmap/index
    - https://goframe.org/container/garray/index
    - https://goframe.org/container/gset/index
    - https://goframe.org/container/gvar/index
    - https://goframe.org/container/gtype/index
    - https://goframe.org/container/glist/index
    - https://goframe.org/container/gvar/index
1. 新增`guuid`模块，用于通用的`UUID`生成：https://goframe.org/util/guuid/index

## 功能改进

### `net`
1. `ghttp`
    - 改进请求流程处理性能；
    - `Server`增加对`Logger`日志对象的配置；
    - `Server`开放了`GetRouterMap`方法，用于获得当前服务的路由列表信息，使得开发者可以更方便地实现自定义权限管理；
    - `Server`配置管理优化；
    - `Client`客户端对象进行了大量的改进工作；
    - `Client`客户端对象增加多文件上传功能；
    - `Request`对象增加`GetError`方法，用于获取当前处理错误；
    - `Request`对象增加独立的视图对象及视图变量绑定功能，使得每个请求可以独立视图管理，也可以通过中间件切换请求对象的视图对象。默认情况下该功能关闭，视图解析时使用的是`Server`对象的视图对象；
    - 改建`Response`对象的`CORS`功能；
    - 增加`Response.WriteTplDefault`方法，用于解析并返回默认的模板内容；
    - 增加更多的单元测试用例；
    - 其他改进；
1. `gipv4`/`gipv6`
    - 一些改进工作；
1. `gtcp`/`gudp`
    - 一些改进工作；

### `database`
1. `gdb`
    - 大量细节改进工作；
    - 去掉查询数据为空时的`sql.ErrNoRows`错误返回，保留`Struct`/`Structs`/`Scan`方法在操作数据为空的该错误返回；
    - 调试模式开启时，输出的SQL语句改进为完整的带参数的SQL，仅作参考；
    - `Where`方法增加对`gmap`数据类型支持，包括顺序性的`ListMap`/`TreeMap`等等；
    - 查询缓存方法`Cache`的缓存时间参数类型修改为`time.Duration`；
    - 修改`Record`/`Result`的数据类型转换方法名称，原有的转换方法标记为`deprecated`；
    - `Record`/`Result`查询结果类型增加`IsEmpty`方法，用于判断结果集是否为空；
    - `Record`类型增加`GMap`方法，用于将查询记录转换为`gmap`类型；
    - 增加`Option`/`OptionOmitEmpty`方法，用于输入参数过滤，包括`Data`参数及`Where`参数：https://goframe.org/database/gdb/empty
    - 增加字段排除方法`FieldsEx`：https://goframe.org/database/gdb/senior
    - 增加日志功能特性：https://goframe.org/database/gdb/senior
    - 改进数据库配置管理：https://goframe.org/database/gdb/config
    - 增加大量单元测试；
1. `gredis`
    - 返回数据类型转换改进：https://github.com/gogf/gf/issues/415
    - 完善单元测试；
    - 其他改进；

### `os`
1. `gcache`
    - 需要注意了：缓存的有效时间参数从`interface{}`类型调整为了`time.Duration`类型，因此不再兼容之前的`int`参数类型，以保证更好的性能；
1. `gfcache`
    - 由于`gcache`组件的缓存时间参数类型的变更，因此该组件的时间参数也变更为了`time.Duration`类型；
1. `gcfg`
    - 增加`Available`方法，用以判断配置是否有效；
1. `gfile`
    - 增加`Chdir`方法，用于工作目录切换；
1. `gtime`
    - 增加`JSON`数据格式的`Marshal`/`UnMarshal`接口实现；

### `container`
1. `gmap`
    - 增加`MapStrAny`方法，用于常见`map`类型的转换；
    - 增加`MapCopy`方法，用于底层`map`数据复制；
    - 增加`FilterEmpty`方法，用于`map`空值过滤；
    - 增加`Pop`/`Pops`方法，用于随机返回`map`中的数据项（并删除）；
    - 增加`Replace`方法，用于给定的`map`数据覆盖底层`map`数据项；
    - 完善单元测试；
    - 其他改进；
1. `garray`
    - 增加`Interfaces`转换方法，返回`[]interface{}`类型；
    - 对排序数组增加`SetComparator`方法用户自定义修改比较器；
    - 完善单元测试；
    - 其他改进；
1. `glist`
    - 增加`NewFrom`方法，基于给定的`[]interface{}`变量创建链表；
    - 增加`Join`方法，用于将链表项使用给定字符串连接为字符串返回；
    - 完善单元测试；
    - 其他改进；
1. `gset`
    - 增加`AddIfNotExistFunc`/`AddIfNotExistFuncLock`方法；
    - 完善单元测试；
    - 其他改进；
1. `gtree`
    - 增加`Replace`方法，用于更新现有树的数据项；
    - 其他改进；
1. `gtype`
    - 一些细节改进工作，不一一列出；
    - 完善基准测试、单元测试；
1. `gvar`
    - 增加`Ints`/`Uints`类型转换方法；
    - 其他改进；

### `crypto`
1. `gmd5`
    - 小细节改进；
1. `gsha1`
    - 小细节改进；

### `text`
1. `gstr`
    - 改进`SplitAndTrim`方法，将`SplitAndTrimSpace`标记为`deprecated`；
    - 增加`TrimStr`方法；
    - 完善单元测试；
    - 其他改进；

### `debug`

1. `gdebug`
    - 增加`CallerFileLineShort`/`FuncPath`/`FuncName`方法；
    - 其他改进；

### `encoding`

1. `gbase64`
    - 增加`EncodeToString`/`EncodeFile`/`EncodeFileToString`/`DecodeToString`方法；
    - 完善单元测试；
1. `gjson`
    - 完善单元测试；

### `frame`

1. `g`/`gins`
    - https://goframe.org/frame/g/index
    - 增加`CreateVar`方法；
    - 完善单元测试；
    - 其他改进；

### `util`

1. `gconv`
    - 改进优化部分类型转换方法性能；
    - 增加`Uints`/`SliceUint`类型转换方法；
    - 增加`UnsafeStrToBytes`/`UnsafeBytesToStr`高性能的类型转换方法；
    - 增加对`MapStrAny`接口方法的支持，用于常见`map`类型的转换；
    - 其他改进；
1. `gvalid`
    - 改进对中国身份证号的识别校验功能；
    - 增加`luhn`银行卡号的校验功能；
1. `grand`
    - 一些性能改进工作；


## Bug Fix
1. 解决`WebSocket`关闭时的`hijacked`报错问题：https://github.com/gogf/gf/issues/381
1. 解决静态文件服务时大文件的内存占用问题；
1. 修复前置`Nginx`后默认情况下的`Cookie`域名设置问题；
1. 修复`gconv.Struct`在属性为`[]struct`并且输入属性参数为空时的转换失败问题：https://github.com/gogf/gf/issues/405
1. 其他一些修复；



# `v1.9.3` (2019-09-24)

该版本实际为`v2.0`的大版本发布，为避免`go module`机制严格要求`v2`版本以上需要修改`import`并加上`v2`后缀，因此使用了`v1.9`版本进行发布。

## 新特性

1. 新增`gf`命令行开发辅助工具：https://goframe.org/toolchain/cli
1. 新增`gres`资源管理器模块：https://goframe.org/os/gres/index
1. 重构`Session`功能，新增`gsession`模块，`WebServer`默认使用文件存储`Session`：https://goframe.org/net/ghttp/session
1. `WebServer`新增中间件特性，并保留原有的HOOK设计，两者都可实现请求拦截、预处理等等特性：https://goframe.org/net/ghttp/router/middleware
1. 新增`gi18n`国际化管理模块：https://goframe.org/i18n/gi18n/index
1. 新增`gini`模块：https://goframe.org/encoding/gini/index
1. `WebServer`新增更便捷的层级路由注册方式：https://goframe.org/net/ghttp/group/level
1. `gcmd`命令行参数解析模块重构，增加`Parser`解析对象：https://goframe.org/os/gcmd/index
1. 新增`gdebug`模块，用于堆栈信息获取/打印：https://goframe.org/debug/gdebug/index


## 重大调整
1. 去掉`1.x`版本中已经被标记为`deprecated`的方法；
1. 调整`container`分类的容器模块，将默认并发安全参数调整为默认非并发安全；
1. 目录调整:
    - 去掉`third`目录，统一使用`go module`管理包依赖；
    - 将原有`g`目录中的模块移出到框架主目录，原有的`g`模块移动到`frame/g`目录；
    - 将原有`geg`示例代码目录名称修改为`.example`；



## 功能改进

1. `ghttp`
    - 改进`Request`参数解析方式：https://goframe.org/net/ghttp/request
    - 改进跨域请求功能，新增`Origin`设置及校验功能：https://goframe.org/net/ghttp/cors
    - `Cookie`及`Session`的`TTL`配置数据类型修改为`time.Duration`;
    - 新增允许同时通过`Header/Cookie`传递`SessionId`；
    - 新增`ConfigFromMap/SetConfigWithMap`方法，支持通过`map`参数设置WebServer；
    - 改进默认的`CORS`配置，增加对常见`Header`参数的默认支持；
    - 新增`IsExitError`方法，用于开发者自定义处理`recover`错误处理时，过滤框架本身自定义的非异常错误；
    - 新增`SetSessionStorage`配置方法，用于开发者自定义`Session`存储；
    - `ghttp.Request`新增更多的参数获取方法；
1. `gdb`
    - 增加对SQL中部分字段的自动转义(`Quote`)功能；
    - 增加对方法操作以及链式操作中的`slice`参数的支持；
    - 增加`SetLogger`方法用于开发者自定义数据库的日志打印；
    - 增加`Master/Slave`方法，开发者可自主选择数据库操作执行的主从节点；
    - 增加对`mssql/pgsql/oracle`的单元测试；
    - 在`debug`模式支持完整带参数整合的SQL语句调试打印；
    - 增加了更多的功能方法；
1. `glog`
    - 新增`Default`方法用于获取默认的`Logger`对象；
    - 新增`StackWithFilter`方法用于自定义堆栈打印过滤；
    - 增加了更多的功能方法；
1. `gfile`
    - 部分方法名称调整：`Get/PutBinContents`修改为`Get/PutBytes`；
    - 增加`ScanDirFile`方法，用于仅检索文件目录，支持递归检索；
    - 增加了更多的功能方法；
1. `gview`
    - 新增`SetI18n`方法用于设置视图对象自定义的`gi18n`国际化对象；
    - 新增对`gres`资源管理器的内置支持；
1. `gcompress`
    - 增加`zip`算法的文件/目录的压缩/解压方法；
    - 文件/目录压缩参数支持多路径；
1. `gconv`
    - 改进对`[]byte`数据类型参数的支持；
    - 新增`Unsafe`转换方法，开发者可在特定场景下使用，提高转换效率；
    - 新增`MapDeep/StructDeep/StructsDeep`方法，支持递归`struct`转换；
1. `gjson/gparser`
    - 改进类型自动识别功能；
    - 新增`LoadJson/LoadXml/LoadToml/LoadYaml/LoadIni`方法用于自定义的数据类型内容加载；
    - 增加了更多的功能方法；
1. `gerror`
    - 改进错误堆栈获取逻辑；
    - 增加了更多的功能方法；
1. `gmap/garray/gset/glist/gvar`
    - 改进并发安全基准测试脚本；
    - 修改`garray.StringArray`为`garray.StrArray`；
    - 增加了更多的功能方法；
1. `gdes`
    - 规范化修改方法名称；
1. `gstr`
    - 增加`Camel/Snake`相关命名转换方法；
    - 增加了更多的功能方法；
1. `genv`
    - 增加了更多的功能方法；


## Bug Fix
1. 修复`gvalid`校验`struct`时的`tag`自定义错误失效的问题；
1. 修复`gcfg`配置管理模块在特定情况下的内容类型自动识别失败问题；
1. 修复`gqueue`在用户主动关闭队列时的并发安全问题；
1. 修复`session`在开发者设置的`TTL`过大时的整型变量溢出问题；