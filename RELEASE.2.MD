# `v1.14.2` (2020-10-27)

# GoFrame

`GF(Go Frame)`是一款模块化、高性能、生产级的Go基础开发框架。实现了比较完善的基础设施建设以及开发工具链，提供了常用的基础开发模块，如：缓存、日志、队列、数组、集合、容器、定时器、命令行、内存锁、对象池、配置管理、资源管理、数据校验、数据编码、定时任务、数据库ORM、TCP/UDP组件、进程管理/通信等等。并提供了Web服务开发的系列核心组件，如：`Router`、`Cookie`、`Session`、`Middleware`、服务注册、模板引擎等等，支持热重启、热更新、域名绑定、`TLS/HTTPS`、`Rewrite`等特性。

## 特点

* 模块化、松耦合设计；
* 模块丰富、开箱即用；
* 简便易用、易于维护；
* 高代码质量、高单元测试覆盖率；
* 社区活跃，大牛谦逊低调脾气好；
* 详尽的开发文档及示例；
* 完善的本地中文化支持；
* 设计为团队及企业使用；

## 支持我们

OSC最佳开源项目评选开始了，如果您喜欢`GoFrame`，欢迎为`GoFrame`投上您宝贵的一票🙏 https://www.oschina.net/p/goframe


# Change Log

由于`GoFrame`是模块化设计，因此每个版本的更新记录都会以模块的形式进行介绍。

重要更新：
1. 将框架内所有的`json`操作从标准库替换为`json-iterator/go`，提高操作效率。
1. 缓存模块重构底层设计，增加适配器设计模式，并增加内存及`Redis`适配器支持。其中内存适配器默认核心模块提供，`Redis`适配器由社区模块提供：https://goframe.org/os/gcache/adapter
1. 增加可自定义的校验规则注册特性：https://goframe.org/util/gvalid/customrule
1. `Web Server`增加所有配置项示例：https://goframe.org/net/ghttp/config/example
1. `ORM`新增基于`Redis`的`SQL`缓存适配器：https://goframe.org/database/gdb/model/cache
1. `ORM`新增模型关联实验特性：https://goframe.org/database/gdb/model/association
1. `ORM`改进时间自动更新特性增加自定义时间字段：https://goframe.org/database/gdb/model/auto-time
1. 错误处理模块新增`Current`及`Next`方法：https://goframe.org/errors/gerror/index



## `net`
1. `ghttp`
    - `Client`
        - 增加`GetVar/PutVar/PostVar`等`*Var`请求方法，用于发起`HTTP`请求获取内容之后直接返回泛型对象，方便类型转换，特别是针对于返回`XML/JSON`的结果处理将会更加简便：https://goframe.org/net/ghttp/client/demo/index
        - 增加`SetProxy/Proxy`方法，用于设置客户端代理，支持`HTTP/Socket5`代理类型：https://goframe.org/net/ghttp/client/demo/proxy
        - 增加`SetRedirectLimit/RedirectLimit`方法，用于设置页面跳转数量限制。
    - `Request`
        - 增加`ParseQuery`, `ParseForm`方法，用于解析指定类型的参数，并绑定到给定的对象。
        - 增加`GetHeader`方法，用于获取指定`Header`参数。
        - 增加`GetRemoteIp`方法，用于获取请求客户端IP。在IP白名单限制时应当使用`GetRemoteIp`而不是`GetClientIp`进行判断，后者可以通过`Header`伪造。
        - 增加`ReloadParam`方法，往往用在中间件处理中，当中间件修改了请求参数，需要通过调用该方法重新解析一下请求参数。
        - 增加`GetRouterMap`方法，用于获得所有的路由参数返回为`map`。
    - `Response`
        - 将`Output`方法名称改为`Flush`，用于将缓冲区的数据写入到客户端数据流中。
    - `Server`
        - `Server`增加所有配置项示例：https://goframe.org/net/ghttp/config/example
        - 增加`SessionCookieOutput`配置，用于控制是否输出`SessionId`到`Cookie`中，默认为开启。
        - 改进路由解析，增加对`URI`带有重复的`/`符号的支持。
        - `Pprof`功能路由支持`Domain`绑定。
        - 其他一些细节改进。
    - `Cookie`
        - 增加`SetHttpCookie`方法，用于根据标准库`http.Cookie`对象设置`Cookie`。
        - 其他一些功能改进
    

## `database`
1. `gdb`
    - 新增模型关联实验特性：https://goframe.org/database/gdb/model/association
    - 改进时间自动更新特性增加自定义时间字段：https://goframe.org/database/gdb/model/auto-time
    - 新增基于`Redis`的`SQL`缓存适配器：https://goframe.org/database/gdb/model/cache
    - 新增对输入参数的键名-字段名自动识别映射特性：https://goframe.org/database/gdb/senior
    - 新增`DB.HasTable`方法，用于判断是否当前数据库存在指定数据表。
    - 新增`Model.HasField`方法，用于判断是否当前数据表存在指定字段。
    - 新增`Model.ScanList`方法，用于智能地将当前`struct`/`slice`绑定到指定的`list`对应属性上。
    - 新增`Result.MapKeyValue`方法，用于将当前`Result`转换为`map[string]Value`类型。
    - 新增`Result.IsEmpty/Len/Size/ScanList`方法。
    - 增加`ListItemValues`及`ListItemValuesUnique`方法，用于自动获取`list`中指定名称的键值或属性值，构成`slice`返回。
    - `SQL`日志内容增加分组名称打印。
    - 改进`DataToMapDeep`方法。
    - 其他一些细节改进工作。
    
1. `gredis`
    - 新增`TLS`特性支持，并支持配置文件配置，https://goframe.org/database/gredis/config

## `container`
1. `gvar`
    - 增加`Scan`及`ScanDeep`方法，用于`struct`/`slice`自动识别转换。
    - 增加`ListItemValues`及`ListItemValuesUnique`方法，用于自动获取`list`中指定名称的键值或属性值，构成`slice`返回。
    - 增加`MapStrAny`接口实现方法。

## `os`
1. `gcache`
    - 增加`GetVar`方法，用于获取缓存数据并返回为泛型对象。
    - 增加`Update`方法，用于仅修改缓存数值，不修改缓存过期时间。
    - 增加`UpdateExpire`方法，用于仅修改缓存过期时间，不修改缓存数值。
    - 重构底层设计，增加适配器设计模式，并增加内存及`Redis`适配器支持。其中内存适配器默认核心模块提供，`Redis`适配器由社区模块提供：https://goframe.org/os/gcache/adapter
    - 注意，本次模块的修改会有部分方法不兼容，部分方法增加了`error`参数返回，升级时请注意查看。编译时将不会通过。
    - 其他一些功能改进。
1. `gfile`
    - 增加`ScanDirFileFunc`方法，用于自定义函数处理的递归目录文件遍历。
    - 改进`Scan*`方法，增加递归层级限制，默认层级限制为`100000`.

1. `gfsnotify`
    - 去掉模块初始化时的`Watcher`对象创建，调整为运行时按需创建，并且增加了并发安全控制。

1. `grpool`
    - 增加`AddWithRecover`方法，用于添加异步任务时给定一个`recover`处理方法，当任务`panic`时交由该`recover`方法处理，防止异步任务`panic`引起整个进程崩溃。
    > 这里解决的痛点是`recover`只能捕获到当前`goroutine`的`panic`，因此只能在创建异步任务的时候指定`recover`处理方法。

1. `gtime`
    - 增加`ParseDuration`方法，增加了对时间单位`d`的支持，表示天。
    - 改进`New`方法，支持通过字符串、时间戳、`time.Time`对象创建`gtime.Time`对象，https://goframe.org/os/gtime/time
    - 改进`Add/AddStr/ToLocation/ToZone/UTCLocal/AddDate/Truncate/Round`方法，这些方法调用时，不再修改当前对象本身，而是创建并返回一个新的`gtine.Time`对象，以便保证和标准库`time.Time`的逻辑一致，防止混淆。
    - 其他一些细节改进。
1. `gtimer`
    - 增加`Reset`方法，用于重置定时任务的计时。
1. `gfcache`
    - 去掉了该模块，该模块的功能作用不是特别大。

## `debug`
1. `gdebug`
    - 新增`GoroutineId`方法，用于获取当前执行的`goroutine id`，仅作调试使用。

## `encoding`
1. `gjson`
    - 新增`GetScan/GetScanDeep`方法。
    - 新增`ToScan/ToScanDeep`方法。
    - 新增`LoadContentType`方法，用于根据指定类型的内容创建`Json`操作对象。
    - 新增`IsValidDataType`方法，用于判断给定的数据类型是否支持解析。
    - 其他一些改进。
    - 单元测试完善。

1. `gcompress`
    - 新增`GzipFile/UnGzipFile`基于`gzip`压缩算法的文件压缩/解压。

## `i18n`
1. `gi18n`
    - 新增`TranslateFormat/TranslateFormatLang/Tf/Tfl`方法: https://goframe.org/i18n/gi18n/index

## `text`
1. `gstr`
    - 增加`SnakeFirstUpperCase`方法，用于在字母大写前增加连接符，并不会处理数字，例如：`SnakeFirstUpperCase("RGBCodeMd5")`将会返回`rgb_code_md5`。

## `util`
1. `gconv`
    - 增加对指针基本类型的转换支持。
    - 增加`Scan/ScanDeep`方法，用于自动识别转换`Struct/[]Struct`。
    - 改进`MapDeep`方法的层级递归处理。
    - 其他一些细节改进，性能改进。

1. `gutil`
    - 增加`ListItemValues`及`ListItemValuesUnique`方法，用于自动获取`list`中指定名称的键值或属性值，构成`slice`返回。
    - 增加`ItemValue`方法，用于获取指定`map/*map/struct/*struct`类型的键值/属性值。
    - 增加`MapOmitEmpty`方法，用于过滤`map`中的空值。
    - 增加`SliceDelete`方法，用于数组项删除。
    - 增加`Try`方法，通过闭包执行给定的方法，如果方法产生`panic`则该方法返回`error`，否则返回`nil`。
    - 改进`TryCatch(try func(), catch ...func(exception interface{}))`为`TryCatch(try func(), catch ...func(exception error))`

1. `gvalid`
    - 增加自定义规则特性，开发者可注册自定义的校验规则：https://goframe.org/util/gvalid/customrule
    - 其他一些功能改进。

## `error`
1. `gerror`
    - 新增`Current`方法，用于获取当前错误层级的`error`接口对象。
    - 新增`Next`方法，用于获取层级错误的下一级错误`error`接口对象。当下一层级不存在时，返回`nil`。
    - 文档更新：https://goframe.org/errors/gerror/index

## Bug Fix
1. 修复`garray`模块的`Unique`方法问题。
1. 修复`glog`中定时器懒初始化时的`goroutine`泄露问题。
1. 修复`gstr`中名称`Case`转换相关方法在名称中带有数字+特殊字符时的名称转换问题。
1. 修复`ghttp`模块中的`CORS`跨域设置`Header`细节问题。
1. 其他BUG修复：https://github.com/gogf/gf/issues?q=is%3Aissue+label%3Abug+is%3Aclosed


# `v1.13.1` (2020-06-10)

# GoFrame

`GF(Go Frame)`是一款模块化、高性能、生产级的Go基础开发框架。实现了比较完善的基础设施建设以及开发工具链，提供了常用的基础开发模块，如：缓存、日志、队列、数组、集合、容器、定时器、命令行、内存锁、对象池、配置管理、资源管理、数据校验、数据编码、定时任务、数据库ORM、TCP/UDP组件、进程管理/通信等等。并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、Middleware、服务注册、模板引擎等等，支持热重启、热更新、域名绑定、TLS/HTTPS、Rewrite等特性。

## 特点

* 模块化、松耦合设计；
* 模块丰富、开箱即用；
* 简便易用、易于维护；
* 高代码质量、高单元测试覆盖率；
* 社区活跃，大牛谦逊低调脾气好；
* 详尽的开发文档及示例；
* 完善的本地中文化支持；
* 设计为团队及企业使用；

## 发展

`GoFrame`开始得比较早，`2011`年始于北京一个智能物联网平台项目，那时还没有这么多物联网的现行标准，`Go`的标准库以及生态也未如此丰富。`2017`年的时候`GF`才开始发布测试版，`2018`年`1024`程序员节日的时候才发布`v1.0`正式版，为`Go`生态发展添砖加瓦。开源以来快速迭代、发展成长，广受开发者和企业的青睐，也有许多的开发者加入了贡献行列。`GF`原本是为开发团队设计的，因此她的开发效率和可维护性做得非常好，有着很高的代码质量以及丰富的单元测试和示例，并且`GF`是目前中文化文档做的最好的`Golang`开发框架。

# Change Log

1. 应多数开发者的要求，框架要求的最低`Golang`运行版本降级为了`v1.11`。
1. 新增`GoFrame`视频教程地址：
    - bilibili：https://www.bilibili.com/video/av94410029
    - 西瓜视频： https://www.ixigua.com/pseries/6809291194665796100/
1. 将不常用的`guuid`模块迁移到 github.com/gogf/guuid 作为社区模块维护，保持`gf`主仓库的轻量级。
1. 新增`guid`模块，用于高效轻量级的唯一字符串生成：https://goframe.org/util/guid/index


## `tool chain`

1. 工具链更新：https://goframe.org/toolchain/cli
1. 新增`gf env`命令，更优雅地查看当前`Golang`环境变量信息。
1. 新增`gf mod path`命令，用于将当前`go modules`包拷贝到`GOPATH`中，以便使用原始的`GOPATH`方式开发项目。
1. 对现有`cli`命令进行了一些改进，提高使用体验；预编译二进制版本在部分平台下提供了`upx`压缩，使得下载的文件更小。

## `container`
1. `garray`
    - https://goframe.org/container/garray/index
    - 简化数组使用方式，支持类似于`var garray.Array`的变量定义使用方式；
    - 增加`Walk`方法，用于自定义的数组元素处理方法；
    - 增加`ContainsI`方法，用于大小写忽略匹配的数组元素项存在性查找；
    - 完善单元测试，代码覆盖率`94%`；
    - 代码改进，提高性能；
    - 修复一些问题；
1. `gchan`
    - 由于该封装包实际意义不是很大，因此从主框架中删除；
1. `glist`
    - https://goframe.org/container/glist/index
    - 简化链表使用方式，支持类似于`var glist.List`的变量定义使用方式；
    - 完善单元测试，代码覆盖率`99%`；
1. `gmap`
    - https://goframe.org/container/gmap/index
    - 简化`Map`使用方式，支持类似于`var gmap.Map`的变量定义使用方式；
    - 完善单元测试，代码覆盖率`81%`；
    - 代码改进，提高性能；
1. `gset`
    - https://goframe.org/container/gset/index
    - 简化集合使用方式，支持类似于`var gset.Set`的变量定义使用方式；
    - 增加`Walk`方法，用于自定义的集合元素处理方法；
    - 完善单元测试，代码覆盖率`90%`；
    - 代码改进，提高性能；
1. `gtree`
    - https://goframe.org/container/gtree/index
    - 简化树型使用方式，支持类似于`var gtree.BTree`的变量定义使用方式；
    - 完善单元测试，代码覆盖率`90%`；
1. `gvar`
    - https://goframe.org/container/gvar/index
    - 完善单元测试，代码覆盖率`69%`；
    - 代码组织结构调整，提高维护性；
    - 代码改进，提高性能；

## `database`

1. `gdb`
    - 增加`Transaction(f func(tx *TX) error) (err error)`接口方法，用于通过闭包实现事务封装处理：https://goframe.org/database/gdb/transaction
    - 去掉不常用的`From`接口方法，改进`Table`及`Model`方法的参数为不定参数，并支持通过不定参数传递表别名：https://goframe.org/database/gdb/chaining/select
    - 增加`DryRun`特性，支持空跑时只执行查询不执行写入/更新/删除操作：https://goframe.org/database/gdb/senior
    - 增加`create_at`, `update_at`写入时间、更新时间字段自动填充特性：https://goframe.org/database/gdb/chaining/auto-time
    - 增加`delete_at`软删除特性：https://goframe.org/database/gdb/chaining/auto-time
    - 增加`Having`链式操作方法，用于`having`条件查询：https://goframe.org/database/gdb/chaining/select
    - `Result`结果对象增加`Chunk`方法，用于自定义的数据分批处理：https://goframe.org/database/gdb/result
    - 改进`Schema`数据库运行时切换特性；
    - 改进对`pgsql`, `mssql`, `sqlite`, `oracle`数据库字段类型的支持；
    - 进一步完善单元测试；
    - 代码组织结构调整，提高维护性；
    - 代码改进，提高性能；
1. `gredis`
    - 增加`MaxActive`连接池参数默认配置为`100`，限制默认的连接数量；
    - 改进`Conn`连接对象的`Do`方法，支持对`map/slice/struct`类型进行自动的`json.Marshal`处理，注意获取数据时使用`DoVar`方法获取：https://goframe.org/database/gredis/usage
    - 完善单元测试，代码覆盖率`72%`；

## `net`
1. `ghttp`
    - 增加`Prefix`及`Retry`客户端链式操作方法；
    - 增加客户端原始请求打印特性：https://goframe.org/net/ghttp/client/demo/dump
    - 增加`ClientMaxBodySize`的服务端配置，用于限制客户端提交的`Body`大小，默认为`8MB`；在涉及到上传的Server中需要增加该配置的大小，在配置文件中指定对应的大小即可，如`ClientMaxBodySize="100MB"`：https://goframe.org/net/ghttp/config
    - 改进`SessionId`生成的随机性，提高`Session`安全性：https://goframe.org/os/gsession/index
    - 改进`ghttp.Server`实现了标准库的`http.Handler`接口，便于与其他第三方的服务如`Prometheus`进行代码集成；
    - 其他大量的代码细节改进工作，提高性能及持久维护性；
    - 完善单元测试，代码覆盖率`61%`；

1. `gipv4`
    - 增加`GetIpArray`方法，用于获取当前主机的所有IPv4地址；
    - 增加`GetMacArray`及`GetMac`方法，用于获取当前主机的`MAC`地址信息；
    - 修改`IntranetIP`方法名称为`GetIntranetIp`，修改`IntranetIPArray`方法名称为`GetIntranetIpArray`；

## `encoding`
1. `gjson`
    - 新增`GetMaps`获取`JSON`内部节点变量方法；
    - 改进`NewWithTag`方法对`map/struct`的处理；
    - 完善单元测试，代码覆盖率`77%`；
1. `gyaml`
    - 升级依赖的第三方`yaml`解析包，解决了`map[interface{}]interface{}`转换问题；

## `error`
1. `gerror`
    - 新增`NewfSkip`方法，用于创建`skip`指定堆栈的错误对象；
    - 放开框架所有的堆栈链路打印，展示错误时真实的链路调用详情；

## `os`
1. `gcache`
    - 增加`GetVar`方法，用于获得可以便捷转换为其他数据类型的"泛型"变量；
    - 标记`Removes`方法废弃，改进`Remove`方法参数为不定参数，统一使用`Remove`方法删除单个/多个键值对；
    - 完善单元测试，代码覆盖率`96%`；

1. `genv`
    - 增加`GetVar`方法，用于获得可以便捷转换为其他数据类型的"泛型"变量；

1. `gfile`
    - 改进`CopyDir/CopyFile`复制目录/文件方法；
    - 新增`ScanDirFunc`方法，用于支持自定义处理回调的目录检索；
    - 完善单元测试，代码覆盖率`64%`；

1. `glog`
    - 增加支持`Context`上下文变量的日志打印特性：https://goframe.org/os/glog/context

1. `gres`
    - 改进打包特性，增强生成二进制文件及Go文件的压缩比，比旧版本增加`20%`压缩率，使得编译生成的二进制文件体积更小；
    - 代码结构改进，提高执行效率及可持久维护性；

1. `gsession`
    - 改进`SessionId`默认生成方法，采用`guid.S`方法生成；
    - 增加`SetId`及`SetIdFunc`方法，用于自定义`SessionId`及自定义的`SessionId`生成方法；

## `frame`
1. `g`
    - 新增`g.Table`方法，用于快速创建数据库模型操作对象；

## `i18n`
1. `gi18n`
    - 新增`GetContent`方法，用于获取指定`i18n`关键字为转译内容；
    - 改进代码细节，提高性能和持久可维护性；
    - 完善单元测试，代码覆盖率`74%`；

## `test`
1. `gtest`
    - 增加`AssertNQ`断言方法，用于强类型的不相等判断；

## `text`
1. `gstr`
    - 增加`SubStrRune`方法，用于支持`unicode`的字符串截取；
    - 增加`StrLimitRune`方法，用于支持`unicode`的字符串截断隐藏；
    - 增加`LenRune`方法，用于替换`RuneLen`方法，统一方法命名风格；
    - 增加`PosRune/PosIRune/PosRRune/PosRIRune`方法，用于支持`unicode`的字符串左右位置查找；
    - 增加`CompareVersionGo`方法，用于`Golang`风格的版本号大小比较；
    - 完善单元测试，代码覆盖率`75%`；

## `util`
1. `gconv`
    - 改进`Convert`转换方法，支持常见`map`类型的转换；
    - 改进类型转换过程中异常错误的捕获，通过`error`返回；
    - 其他一些细节改进；
    - 完善单元测试，代码覆盖率`63%`；

1. `grand`
    - 增加`B`方法，用于获得随机的二进制数据；
    - 改进代码底层实现，部分接口性能提高`50%`；
    - 完善单元测试，代码覆盖率`74%`；

1. `guid`
    - 新增`guid`模块，用于高效轻量级的唯一字符串生成：https://goframe.org/util/guid/index

1. `gutil`
    - 增加`MapContains`方法，用于判断map中是否包含指定键名；
    - 增加`MapDelete`方法，用于删除map中指定的键名，可以为多个键名；
    - 增加`MapMerge`方法，用于合并两个map；
    - 增加`MapMergeCopy`方法，用于拷贝多个map；
    - 增加`MapContainsPossibleKey`方法，用于查找指定键名，忽略大小写及字符`'-'/'_'/'.'/' '`；

1. `gvalid`
    - 所有默认的错误提示改为了英文；
    - 错误提示的配置改为了通过`i18n`来配置实现，以便支持国际化：https://goframe.org/util/gvalid/message
    - 身份证号规则名称从`id-number`改为了`resident-id `；
    - 银行卡号规则名称从`luhn`改为了`bank-card`；
    - 完善单元测试，代码覆盖率`96%`；

## Bug Fix
1. 修复`gcompress`的多文件`zip`压缩问题；
1. 修复`ghttp.Client`获取返回的过期`Cookie`的问题；
1. 修复`gres.File`对于`http.File`接口的实现细节；
1. 修复`garray.Pop*`方法的边界问题；
1. 修复`gres`中`Readdir`方法参数为`0`时报错的问题；
1. 其他一些修复：https://github.com/gogf/gf/issues?q=is%3Aissue+label%3Abug




# `v1.12.1` (2020-03-31)

大家好啊！久等啦！

由于自从上次版本的发布以来，越来越多小伙伴加入了`GF`的大家庭，并提供了许多不错的建议和反馈，这次版本对其中大部分反馈进行了处理，包括大部分的改进建议和部分新特性，因此这次的版本发布时隔了两个多月。`GF`非常注重代码质量以及可持续维护性，这次版本也进一步对框架大部分模块的示例、注释和单元测试用例进行了完善，目前单元测试用例数量约为`1991`例，代码覆盖率为`71%`，覆盖了所有模块的绝大部分主要功能。

`GF`框架提供了比较常用、高质量的基础开发模块，轻量级、模块化、高性能，推荐大家阅读框架源码了解更多细节。本次发布有个别的不兼容升级，往往批量替换即可，以下`change log`比较完善，建议升级前仔细阅读。

本次发布即意味下一版本开发计划的开启，欢迎更多小伙伴参与开源贡献：https://github.com/gogf/gf/projects/8

感谢大家支持！Enjoy your `GF`！


# GoFrame

`GF(Go Frame)`是一款模块化、高性能、生产级的Go基础开发框架。实现了比较完善的基础设施建设以及开发工具链，提供了常用的基础开发模块，如：缓存、日志、队列、数组、集合、容器、定时器、命令行、内存锁、对象池、配置管理、资源管理、数据校验、数据编码、定时任务、数据库ORM、TCP/UDP组件、进程管理/通信等等。并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、Middleware、服务注册、模板引擎等等，支持热重启、热更新、域名绑定、TLS/HTTPS、Rewrite等特性。

## 特点
* 模块化、松耦合设计；
* 模块丰富，开箱即用；
* 简便易用，易于维护；
* 社区活跃，大牛谦逊低调脾气好；
* 高代码质量、高单元测试覆盖率；
* 详尽的开发文档及示例；
* 完善的本地中文化支持；
* 更适合企业及团队使用；

## 地址
- 官网：https://goframe.org
- 主库：https://github.com/gogf/gf 
- 码云：https://gitee.com/johng/gf 

# Change Log

从`GF v1.12`版本开始，框架要求的最低`Golang`运行版本为`v1.13`，由于`Golang`新版本都是向后兼容的，因此推荐大家更新使用`Golang`新版本：https://golang.google.cn/dl/

> 本次版本也新增了`Swagger`的工具及插件支持，另行独立发布。

## `tool chain`
1. `gen model`命令新增对`pgsql/mssql/sqlite/oracle`的模型生成支持。
1. `gen model`命令生成模型新增公开包变量`Columns`用于获得表的字段名称。

## `net`
1. `ghttp`
    - 注意：从该版本开始，`Server`默认关闭了平滑重启特性。开发者可以通过相应的配置选项打开。
    - 改进`Client.Get`方法，增加可选的请求参数。
    - 新增`Client`链式操作方法：`Header`, `HeaderRaw`, `Cookie`, `ContentType`, `ContentJson`, `ContentXml`, `Timeout`, `BasicAuth`, `Ctx`：https://goframe.org/net/ghttp/client/chain
    - 新增`Request.GetCtx/GetCtxVar/SetCtxVar/Context`上下文变量管理方法，用于请求内部的上下文变量特性：
        - 自定义变量：https://goframe.org/net/ghttp/request/custom
        - 上下文变量：https://goframe.org/net/ghttp/request/context
    - 新增`Request.GetUploadFile/GetUploadFiles`方法，以及`UploadFile`类型，极大简化文件上传处理逻辑：https://goframe.org/net/ghttp/client/demo/upload
    - 新增`Request.GetPage`方法，用于便捷地获得分页对象：
        - 基本介绍：https://goframe.org/util/gpage/index
        - 动态分页：https://goframe.org/util/gpage/dynamic
        - 静态分页：https://goframe.org/util/gpage/static
        - Ajax分页：https://goframe.org/util/gpage/ajax
        - URL模板：https://goframe.org/util/gpage/template
        - 自定义分页：https://goframe.org/util/gpage/custom
    - 改进`Response.Redirect*`方法，增加自定义的跳转HTTP状态码参数。
    - 改进`CORS`特性，完善跨域功能处理，并完全遵守`W3C`关于`OPTIONS`请求方法的规范约定：https://goframe.org/net/ghttp/cors
    - 模板视图对象增加`Request`内置变量，用于模板获得请求参数。由于`GF`框架的模板引擎采用两级缓存设计，减少`IO`开销的同时提升了执行效率，即使增加了大量的内置变量以及内置方法，经过大规模的性能测试，性能比其他`WebServer`库相同逻辑下高出`50% - 200%`的效率。
    - 新增`Server`实验性的`Plugin`特性。
    - 改进`Server`底层路由存储、检索逻辑，优化代码，提升效率。
    - 改进分组路由注册的源码位置记录功能，当路由注册冲突时能够精准定位及提示重复路由源码位置。
    - 改进`Server`日志处理。
    - 完善单元测试。


## `database`
1. `gdb`
    - 代码重构改进，增加`Driver`驱动接口设计，方便开发者自定义驱动实现。新增`Register`包方法，用于开发者注册自定义的数据库类型驱动：https://goframe.org/database/gdb/driver
    - 新增`GetArray`接口及实现，用于获取指定字段列的数据，构造成数组返回：https://goframe.org/database/gdb/chaining/select
    - 新增`InsertIgnore`接口及实现，用于写入时忽略写入冲突，仅对`mysql`数据库类型有效：https://goframe.org/database/gdb/chaining/insert-save
    - 新增`Schema`接口及实现，用于动态切换并获取指定名称的数据库对象：https://goframe.org/database/gdb/chaining/schema
    - 新增`FieldsStr/FieldsExStr`模型方法，用于获取表字段，并构造成字符串返回：hhttps://goframe.org/database/gdb/chaining/fields-retrieve
    - 新增`LockUpdate/LockShared`模型链式操作方法，用于实现悲观锁操作：https://goframe.org/database/gdb/chaining/lock
    - 改进`Where/Data`方法对更新参数输入方式的支持，提高灵活性：
        - https://goframe.org/database/gdb/chaining/select
        - https://goframe.org/database/gdb/chaining/update-delete
    - 查询结果对象`Result`新增`Array`方法，用于获得指定字段的数值，构造成数组返回：https://goframe.org/database/gdb/result
    - 改进`OmitEmpty`方法对`Data`输入参数的过滤，当给定的`Data`参数为空时间对象时，将会被过滤。
    - 增加默认的数据库连接池参数：`MaxIdleConns=10`。
    - 其他一些改进。
    - 完善单元测试。

1. `gredis`
    - 增加/修改默认的数据库连接池参数：`MaxIdle=10`, `IdleTimeout=10s`, `MaxConnLifetime=30s`, `Wait=true`。
    - 完善单元测试。

## `container`
1. 所有容器对象新增`UnmarshalValue(interface{}) error`接口方法实现，用于`gconv`转换时根据任意类型变量创建/设置对象内容，`GF`框架的所有容器对象均实现了该接口。

1. `garray`
    - 新增`RemoveValue`方法，用于根据数值检索并删除元素项。
    - 新增`FilterNil`方法，用于遍历并删除数组中的`nil`元素项。
    - 新增`FilterEmpty`方法，用于遍历并删除数组中的空值元素项，空值包括：`0, nil, false, "", len(slice/map/chan) == 0`。
    - 改进`Set/Fill/InsertBefore/InsertAfter`方法严谨性，将原有的链式操作返回值对象修改为`error`返回值。
    - 改进`Get/Remove/PopRand/PopLeft/PopRight/Rand`方法严谨性，增加`found`的`bool`返回值，标识当前方法是否有数据返回，当空数组或者操作越界时`found`返回值为`false`。
    - 改变`Rands`方法原有逻辑，保证按照给定大小返回随机数组。
    - 完善单元测试。

1. `gpool`
    - 调整缓存池过期时间参数类型，旧版本为`int`类型表示毫秒，新版本为`time.Duration`类型使得时间控制更灵活。
    - 内部代码优化，性能改进。
    - 完善单元测试。
    - 完善注释。

## `os`
1. `glog`
    - 增加`Rotation`日志滚动切分特性，新增按照文件大小或过期时间进行日志切分，并支持切分文件数量限制、对日志文件进行自动压缩、可自定义压缩级别（`1-9`）、支持对切分文件过期时间清理：https://goframe.org/os/glog/rotate
    - 新增`LevelPrefixes`特性，支持对日志级别的前缀名称进行自定义：https://goframe.org/os/glog/level
    - 新增`SetLevelStr`方法，并增加按照字符串进行日志级别配置的特性：
        - https://goframe.org/os/glog/level
        - https://goframe.org/os/glog/config
    - 新增`SetDefaultLogger`包方法，用于设置全局默认的`Logger`对象。

1. `gres`
    1. 改进资源内容编码设计，采用新的压缩算法，减少资源文件大小，例如原本`15MB`的网站静态资源文件（`css/js/html`等），资源文件打包后约为`4MB`左右：https://goframe.org/os/gres/index
    1. 注意：该改进与旧版本无法兼容，需要重新打包原有的资源文件。
    1. 完善单元测试。

1. `gcfg`
    - 去掉配置对象属性的原子并发安全控制，改为普通变量类型。由于配置管理是最常用模块之一，因此确保高效的设计及方法实现。
    - 单例对象做较大调整：为方便多文件场景下的配置文件调用，简便使用并提高开发效率，因此当给定的单例名称对应的`toml`配置文件在配置目录中存在时，将自动设置该单例对象的默认配置文件为该文件。例如：`g.Cfg("redis")`获取到的单例对象将会默认去检索并设置默认的配置文件为`redis.toml`，当该文件不存在时，则使用默认的配置文件（`config.toml`）：https://goframe.org/net/ghttp/config
    - 完善单元测试。

1. `gview`
    - 新增`concat`内置字符串拼接方法：https://goframe.org/os/gview/function/buildin
    - 完善单元测试。

1. `gfile`
    - 改进`SelfPath`获取当前执行文件路径方法，提高执行效率。
    - 改进`Join`文件路径连接方法，防止多余的路径连接符号。
    - 改建`GetContents`文件内容获取方法执行效率，降低内存占用。
    - 新增`StrToSize`方法，用于将大小字符串转换为字节数字，大小字符串例如`10m`, `5KB`, `1.25Gib`等。
    - 新增`ReadLines/ReadByteLines`方法，用于按行读取指定文件内容，并给定读取回调函数。
    - 完善单元测试。

1. `gtime`
    - 改进`gtime.Time`对象实现，统一字符串打印时间格式为`Y-m-d H:i:s`，如：`2020-01-01 12:00:00`。

1. `gcmd`
    - 命令行解析方法增加`strict`参数，用于设置当前解析是否严格解析，严格解析下如果给定了非法的选项，将会停止解析并返回错误。默认情况下为非严格解析。

1. `genv`
    - 改进`Remove`删除环境变量键值对方法，增加对多个键值对环境变量的删除。

1. `gfpool`
    - 改进代码实现，提高效率。
    - 完善单元测试。

1. `gfsnotify`
    - 改进包初始化方法，当系统原因引起默认`Watcher`对象创建失败时，直接`panic`。

1. `gproc`
    - 改进`SearchBinaryPath`方法。
    - 改进`Process.Kill`方法在`windows`及`*niux`平台下不同表现的处理。

## `encoding`
1. `gjson`
    - 代码改进、完善注释、新增大量代码示例。
    - 文档更新：
        - 基本介绍：https://goframe.org/encoding/gjson/index
        - 对象创建：https://goframe.org/encoding/gjson/object
        - 层级访问：https://goframe.org/encoding/gjson/pattern
        - Struct转换：https://goframe.org/encoding/gjson/conversion-struct
        - 动态创建修改：https://goframe.org/encoding/gjson/dataset
        - 数据格式转换：https://goframe.org/encoding/gjson/conversion-format

## `frame`
1. `g`
    - 新增`IsNil`方法，用于判断当前给定的变量是否为`nil`，该方法有可能会使用到反射来实现判断。
    - 新增`IsEmpty`方法，用于判断当前给定的变量是否为`空`，该方法优先使用断言判断但有可能会使用到反射来实现判断。空值包括：`0, nil, false, "", len(slice/map/chan) == 0`。
    - 标记废弃`SetServerGraceful`方法，转而统一交给`Server`的配置来管理。
1. `gins`
    - 改进代码结构，方便维护。
    - 改进、完善单元测试。
1. `gmvc`
    - 新增`M`类型，为`*gdb.Model`的别名简称，用于工具链自动生成模型中的`M`属性。

## `text`
1. `gstr`
    - 新增`HasPrefix/HasSuffix`方法。
    - 新增`OctStr`方法，用于将八进制字符串转换为其对应的`unicode`字符串，例如转换为中文。常用于`gRPC`底层通信编码中。
    - 完善单元测试。

## `debug`
1. `gdebug`
    - 改进代码结构，方便维护。
    - 新增`TestDataPath`方法，用于当前包单元测试中获得当前包中`testdata`目录的绝对路径。

## `util`
1. `gconv`
    - 改进`String`字符串转换方法，增加对`time.Time/*time.Time/gtime.Time`类型的内置支持。
    - 改进`Map/Struct`转换方法，增加对一些特殊场景的细节处理。经过大规模的使用，大量的反馈改进，不断完善了细节。
    - 改进`Struct`转换方法，增加对`UnmarshalValue(interface{}) error`接口的支持。
    - 完善单元测试。

1. `grand`
    - 注意：不兼容调整，原有的`Str`方法更名为`S`方法，用以获取指定长度的随机字符串、数字，并增加`symbol`参数，指定是否可以随机返回特殊可见字符。
    - 新增`Str`方法，用于从指定的字符串字符中随机获取指定长度的字符串。该方法同时支持`unicode`字符串，例如：中文：https://goframe.org/util/grand/index
    - 新增`Symbols`方法，用于随机返回指定场孤独的特殊可见字符：https://goframe.org/util/grand/index
    - 完善单元测试。

1. `gvalid`
    - 长度校验规则增加对`unicode`字符串的支持，例如：中文。

# Bug Fix
1. 修复`Server`的视图对象配置失效问题。
1. 修复`Server`中间件在中间件`panic`情况下，忽略`Middleware.Next`方法控制，导致鉴权中间件失效的问题。
1. 修复`gudp.Server`在请求包大小超过`64bytes`时的断包问题，并调整默认缓冲区大小为`1024byte`，开发者可自定义缓冲区大小。
1. 修复`gfile.MTimeMillisecond`方法返回错误的文件修改毫秒时间戳。
1. 修复`gconv.Int64`对负数转换的支持。
1. 其他一些修复。
1. 详见：https://github.com/gogf/gf/issues?q=label%3Abug



# `v1.11.2` (2020-01-14)

`GF(Go Frame)` https://goframe.org 是一款模块化、高性能、生产级的Go基础开发框架。实现了比较完善的基础设施建设，包括常用的核心开发组件， 如：缓存、日志、文件、时间、队列、数组、集合、字符串、定时器、命令行、文件锁、内存锁、对象池、连接池、资源管理、数据校验、数据编码、文件监控、 定时任务、数据库ORM、TCP/UDP组件、进程管理/通信、并发安全容器等等。 并提供了Web服务开发的系列核心组件，如：Router、Cookie、Session、Middleware、服务注册、配置管理、模板引擎等等， 支持热重启、热更新、多域名、多端口、多服务、HTTPS、Rewrite等特性。

`GF`有着丰富的基础模块、完善的工具链、详尽的开发文档。开源近两年以来，`GF`得到越来越多小伙伴的肯定和支持，从寂寂无名到现在被广泛应用于微服务、物联网、区块链、电商系统、银行系统等企业级的生产项目中，经历了百万级、千万级项目的考验，2019年度被码云`gitee`评选为`GVP`最有价值开源项目。`GF`正在快速地成长中，目前保持着1-2个月迭代版本的发布规律，社区活跃，欢迎加入`GF`大家庭。

最后，祝大家2020新年快乐，鼠年大吉！



## 新特性

1. 新年新气象，官网文档大量更新：https://goframe.org/index
1. `GF`工具链更新：https://goframe.org/toolchain/cli
    - 新增`gf run`热编译运行命令；
    - 新增`gf docker` Docker镜像编译命令；
    - 新增`gf gen model` 强大的模型自动生成命令；
    - `gf build`命令增加对配置文件配置支持；
    - 大量命令行工具改进工作；
    - 新增自动代理设置特性；
1. 数据库`ORM`新特性：
    - 增加`prefix`数据表前缀支持：https://goframe.org/database/gdb/config
    - 新增`Schema`数据库对象并改进数据库切换特性：https://goframe.org/database/gdb/chaining/schema
    - 新增`WherePri`方法，用于自动识别主键的条件方法：https://goframe.org/database/gdb/chaining/select
    - 文档及示例大量更新，覆盖95%以上的功能特性；



## 功能改进

### `container`
1. `garray`
    - 新增`New*ArrayRange`方法，用于初始化创建指定数值范围的数组。
    - 新增`Iterator*`方法，用于数组项元素回调遍历。
    - 完善单元测试。
1. `gvar`
    - 改进`MapStrStr`、`MapStrStrDeep`方法实现。

### `net`
1. `ghttp`
    - 改进HTTP客户端，增加对提交参数的自动`Content-Type`识别功能。
    - `Request`对象增加`Parse`方法，用于快捷的对象转换即参数校验。
    - `Request.GetPost*`方法全部标记为`deprecated`，统一客户端参数提交方式为`QueryString`, `Form`, `Body`。
    - 去掉`Response`模板解析时的`Get`/`Post`内置变量，新增`Query`, `Form`, `Request`内置变量：https://goframe.org/net/ghttp/response/template
    - 改进`Response.WriteJson*`及`Response.WriteXml*`方法，增加对`string`, `[]byte`类型参数的支持。
    - `Server`新增`GetRouterArray`方法，用于向应用层暴露并获取`Server`的路由列表。
    - `Server`新增`Use`方法，该方法为`BindMiddlewareDefault`的别名，用以全局中间件的注册。
    - `Server`新增`RouteOverWrite`配置项，用于控制是否在注册路由冲突时自动覆盖，默认关闭并提示。
    - `Server`新增`Graceful`配置项，用于在单服务场景下控制平滑重启特性的开启/关闭，默认开启。
    - 完善单元测试。
1. `gtcp`
    - 改进简单协议下的数据包发送接收功能。
    - 将连接池默认的缓存过期时间`30`秒修改为`10`秒。
    - 完善单元测试。

### `database`
1. `gdb`
    - 新增`As`数据表别名方法。
    - 改进数据表、字段的安全字符自动识别添加功能。
    - 新增`DB`数据库对象切换方法。
    - 新增`TX`链式操作事务支持方法。
    - 完善单元测试。
### `os`
1. `gcfg`
    - 新增`GetMapStrStr`方法。
1. `gcmd`
    - 增加参数解析的`strict`严格参数，默认严格解析，不存在指定参数/选项名称时则报错返回。
1. `genv`
    - 改进`Remove`方法支持多个环境变量的删除。
1. `gfile`
    - 改进`TempDir`临时目录获取方法，在`*nix`系统下默认为`/tmp`目录。
    - 新增`ReadLines`, `ReadByteLines`方法，用以按行回调读取文件内容。
    - 新增`Copy*`方法，用以文件/目录的拷贝，支持递归。
    - 新增`Replace*`方法，用以目录下的文件内容替换，支持递归。
    - 改进`Scan*`方法，用以检索并返回指定目录下的所有文件/目录，支持文件模式指定，支持递归。
    - 完善单元测试。
1. `gproc`
    - 改进命令行运行方法。
    - 改进`Shell`命令文件检索逻辑。
    - 改进实验性的进程间通信设计。
1. `gtime`
    - 将包方法以及`Time`对象的时间戳方法`Second`, `Millisecond`, `Microsecond`, `Nanosecond`标记为废除，
    并新增`Timestamp`, `TimestampMilli`, `TimestampMicro`, `TimestampNano`替换。
    - 需要注意的是以上修改可能和老版本存在兼容性问题。

1. `gview`
    - 解析功能、缓存设计改进。
    - 新增`encode`, `decode`HTML编码/解码模板函数。
    - 新增`concat`字符串拼接模板函数。
    - 新增`dump`模板函数，功能类似于`g.Dump`方法。
    - 新增`AutoEncode`配置项，用于自动转码输出的`HTML`内容，常用于防止`XSS`，默认关闭。需要注意的是该特性并不会影响`include`内置函数: https://goframe.org/os/gview/xss
    - 单元测试完善。

### `crypto`
1. `gmd5`
    - 增加`MustEncrypt`, `MustEncryptBytes`, `MustEncryptString`, `MustEncryptFile`方法。
1. `gsha1`
    - 增加`MustEncryptFile`方法

### `encoding`
1. `gbase64`
    - 新增`MustEncodeFile`, `MustEncodeFileToString`, `MustDecode`, `MustDecodeToString`方法。
1. `gjson`/`gparser`
    - 新增`GetMapStrStr`方法。
    - 新增`Must*`方法，用于指定数据格式的转换失败时产生`panic`错误，而不会返回`error`参数。

### `util`
1. `gconv`
    - 改进`Convert`方法增加对`[]int32`, `[]int64`, `[]uint`, `[]uint32`, `[]uint64`, `[]float32`, `[]float64`数据类型的转换支持。
    - 改进`String`字符串转换方法对指针参数的支持。
    - 改进`Map*` Map转换方法的代码结构及性能。
    - 新增`Floats`, `Float32s`, `Float64s`对`[]float32`, `[]float64`类型转换方法。
    - 新增`Ints`, `Int32s`, `Int64s`对`[]int`, `[]int32`, `[]int64`类型转换方法。
    - 新增`Uints`, `Uint32s`, `Uint64s`对`[]uint`, `[]uint32`, `[]uint64`类型转换方法。
    - 完善单元测试。


### `frame`
1. `gins`
    - 所有的单例对象在获取失败时产生`panic`错误。

## Bug Fix
1. 增加对常见错误路由格式例如`/user//index`的兼容支持。
1. 修复`gtcp`/`gudp`在数据接收时的间隔时间单位问题。
1. 修复`gfile`/`gspath`/`gfsnotify`包对文件的存在性判断不严谨问题。
1. 修复`gproc.Kill`方法在`windows`系统下的运行阻塞问题。
1. 修复`gstr.TrimLeftStr`/`gstr.TrimRightStr`在被替换字符串长度小于替换字符串长度时的数组溢出问题。



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