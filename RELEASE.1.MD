# `v1.8.0` (2019-07-15)

## 新功能改进
1. 框架目前 `69` 个开发模块（不包括内部模块），原生代码 `65302` 行（不包含第三包依赖包），单元测试覆盖率达到`77%`；
1. 新增`gerror`错误处理模块：https://goframe.org/errors/gerror/index
1. 改进`gcharset`字符编码转换模块，支持更多的字符集：https://goframe.org/encoding/gcharset/index
1. 新增`gmutex`模块，基于`channel`实现的高级互斥锁模块，支持更丰富的互斥锁特性：https://goframe.org/os/gmutex/index
1. 改进`glog`日志模块：
    - 新增日志异步输出特性：https://goframe.org/os/glog/async
    - 新增`Flags`额外功能特性：https://goframe.org/os/glog/flags
    - 新增`Json`数据格式输出：https://goframe.org/os/glog/json
    - 新增自定义`Writer`接口特性：https://goframe.org/os/glog/writer
    - **修改`Backtrace`名称为`Stack`，并改进调用堆栈输出格式；**
    - 新增`Expose`方法暴露内部默认`Logger`对象；
1. 改进`gdb`数据库ORM模块：
    - **改进错误处理，当数据库操作没有查询到数据时，`error`返回`sql.ErrNoRows`**：https://goframe.org/database/gdb/error
    - 改进`Update`/`Delete`方法支持`Order BY`及`LIMIT`特性；
    - 数据库链式操作及方法操作中，预处理变量参数支持`slice`参数：https://goframe.org/database/gdb/chaining/model
    - **修改`Priority`权重配置名称为`Weight`；**
    - 新增`Debug`配置，可配置开启/关闭调试特性：https://goframe.org/database/gdb/config
    - 新增`Offset`方法，该方法为可选链式操作方法，`pgsql`数据库可直接通过`Limit`方法第二个参数自动识别为`Offset`语法；
    - 改进数据库动态切换特性，支持不同数据库类型的当前操作数据库切换；
    - 改进简化配置文件结构：https://goframe.org/database/gdb/config
1. 改进`gconv`数据转换模块：
    - 对结构体对象转换时支持更多的标签：`gconv/c/json`；
    - 支持`*struct/[]struct/[]*struct`自动初始化创建对象/数组：https://goframe.org/util/gconv/struct
    - 新增`Strusts/StrctsDeep`方法，用于结构体数组的递归转换；
    - 新增`StructDeep`方法，用于对结构体对象的递归转换；
    - 新增`MapDeep`方法，用于对结构体属性的递归转换；
1. 改进`ghttp`模块：
    - 改进`ghttp`模块的分组路由功能，完善逻辑处理细节，程序更加稳健；
    - 改进`ghttp.Request.Get*ToStruct`方法，支持`params/param/p`标签，支持结构体递归转换，并且支持`**struct`参数的对象自动初始化；
    - 改进`ghttp.CORSDefault`的跨域设置参数，`AllowOrigin`参数调整为`*`；
1. 改进`gvalid`数据校验模块：
    - 增加对校验标签`gvalid/valid/v`的支持；
    - 改进`CheckStruct`支持对结构体对象的递归校验：https://goframe.org/util/gvalid/checkstruct
1. 改进`gtcp`TCP通信模块：
    - 改进通信包协议设计，更加轻量级高效：https://goframe.org/net/gtcp/conn/pkg
    - 改进`TCP Server`增加对`TLS`的支持：https://goframe.org/net/gtcp/tls
    - 增加`Server.Cloce`服务端关闭方法；
1. 改进`gproc`模块的通信数据结构，并使用`gtcp`的轻量级包协议重构消息发送逻辑；
1. 改进`gqueue`模块增加数据同步缓冲机制，解决大数据量下的内存占用及延迟问题；
1. 改进`gmlock`模块，使用`gmutex`模块替换内部的互斥锁，增加更多的操作方法；
1. 改进`gaes`加密模块，增加`CBC`模式的加密/解密方法：
1. 改进`garray.Range/SubSlice`方法，改进设计，提高性能；
1. 改进`gjson`/`gparser`模块实现`MarshalJSON`接口以实现自定义的`JSON`数据格式转换；
1. **改进`crypto`分类下模块的方法返回值，增加`error`错误变量返回，以保证更严谨的接口设计风格；**
1. **改进`gbase64`模块，输入输出类型发生改变，并增加多个相关方法；**
1. 改进`gflock`修改方法名`UnLock`为`Unlock`，新增`IsRLocked`方法；
1. 新增`gfile.CopyFile/CopyDir`方法，用于文件及目录的复制；
1. 改进`gjson/gparser/gvar/gcfg`模块增加更多的类型转换方法；
1. 改进`gcache`模块，过期时间参数支持`time.Duration`类型；
1. 新增`internal/structs`包，强大且便捷的结构体解析，并改进框架中所有涉及到结构体反射处理的模块；
1. 改进`gbinary`增加封装方法对`BigEndian`的支持；

## Bug Fix
1. 修复`garray.Search`返回值问题；
1. 修复`garray.Contains`, `garray.New*ArrayFromCopy`方法逻辑问题；
1. 修复`gjson.Remove`删除`slice`参数问题；
1. 修复`gtree.AVLTree.Remove`方法返回值问题；
1. 修复`gqueue.Size`不准确的大小问题；
1. 修复`queue.Close`问题；
1. 修复`gcache.GetOrSetLockFunc`当回调函数返回`nil`结果时的死锁问题；
1. 修复`gfsnotify.Add`方法默认递归监控添加失效问题；
1. 修复`gdb.Model.Scan`在某些参数类型下的失效问题；

## 注意事项

请注意以上粗体文字部分，如有使用，在您升级时可能会出现不兼容性。



# `v1.7.0` (2019-06-10)
## 新功能/改进
1. 重构改进`glog`模块：
    - 去掉日志模块所有的锁机制，改为无锁设计，执行性能更加高效
    - 增加日志内容的异步输出特性：https://goframe.org/os/glog/async
    - 增加日志输出内容的`Json`格式支持：https://goframe.org/os/glog/json
    - 增加`Flags`额外特性支持，包括文件行号打印、自定义时间格式、异步输出等特性控制：https://goframe.org/os/glog/flags
    - 增加`Writer`接口支持，便于开发者进行自定义的日志功能扩展，或者与第三方服务/模块对接集成：https://goframe.org/os/glog/writer
    - 修改`SetStdPrint`方法名为`SetStdoutPrint`
    - 修改链式方法`StdPrint`方法名为`Stdout`
    - 标记淘汰`*fln`日志输出方法，`*f`方法支持自动的换行输出
    - 新增更多的链式方法支持：https://goframe.org/os/glog/chain
1. 重构改进`gmap`模块：
    - 增加更多数据格式支持：`HashMap`/`ListMap`/`TreeMap`
    - 简化类型名称，如`gmap.StringInterfaceMap`简化为`gmap.StrAnyMap`
    - 改进`Map/Keys/Values`方法以提高性能
    - 修改`BatchSet`/`BatchRemove`方法名为`Sets`/`Removes`
    - 新增更多功能方法支持：https://goframe.org/container/gmap/index
1. 改进`gtime`时间模块：
    - 增加并完善更多的类`PHP`时间格式支持
    - 新增更多功能方法，如`FormatTo`/`LayoutTo`等等
    - 详见开发文档：https://goframe.org/os/gtime/index
1. 改进`gdb`数据库模块：
    - 增加对继承结构体的数据转换支持：https://goframe.org/database/gdb/senior
    - 新增`GetLastSql`方法，用以在调试模式下获取最近一条执行的SQL语句
    - 其他的细节处理改进
1. 改进`gtcp`通信模块：
    - 完善处理细节，提高通信性能；
    - 增加`TLS`服务端/客户端通信支持：https://goframe.org/net/gtcp/tls
    - 增加简单协议支持，便于开发者封包/解包，并解决粘包/半包问题：https://goframe.org/net/gtcp/conn/pkg
    - TCP服务端增加`Close`方法
    - 更多细节查看开发文档：https://goframe.org/net/gtcp/index
1. 改进`gconv`类型转换模块
    - 修改`gconv.TimeDuration`转换方法名称为`gconv.Duration`
    - 新增`gconv.StructDeep`及`gconv.MapDeep`方法，支持递归转换
    - 详见开发文档：https://goframe.org/util/gconv/struct
1. 改进`ghttp`模块：
    - 日志输出增加`http/https`字段：https://goframe.org/net/ghttp/logs
    - 新增`ghttp.Server.SetKeepAlive`设置方法，用以开启/关闭`KeepAlive`特性
    - 增加`ghttp.Request.GetUrl`方法，用以获取当前完整的URL请求地址
    - `ghttp.Client`客户端支持开发者自定义`Transport`属性，`ghttp.Client.Post`方法支持`浏览器模式`：https://goframe.org/net/ghttp/client
1. 新增`gtree`树形数据结构容器支持：https://goframe.org/container/gtree/index
1. 改进`gudp`通信模块，具体请参考开发文档：https://goframe.org/net/gudp/index
1. 改进`gcfg`配置管理模块，所有`Get*`方法增加默认值支持：https://goframe.org/os/gcfg/index
1. `gredis`模块新增`DoVar`/`ReceiveVar`方法以便于开发者对执行结果进行灵活的数据格式转换：https://goframe.org/database/gredis/index
1. `gcache`模块`BatchSet`/`BatchRemove`方法名修改为`Sets`/`Removes`
1. 改进`gjson`/`gparser`模块，增加更多方法：https://goframe.org/encoding/gjson/index
1. 改进`gfile.MainPkgPath`方法，以支持不同平台的开发环境；
1. 改进`grpool`协程池模块，提高执行性能：https://goframe.org/os/grpool/index
1. 改进`TryCatch`方法，当开发者不传递`Catch`参数时，默认抑制并忽略错误的处理
1. 改进`gmlock`模块，增加`TryLockFunc`/`TryRLockFunc`方法，并且为`gmlock.Mutex`高级互斥锁对象增加`TryLockFunc`/`TryRLockFunc`方法
1. 去除`gvar.VarRead`接口类型支持

## Bug Fix
1. 解决`gdb`模块与其他第三方`ORM`模块同时使用的冲突；
1. 修复`gcron.AddOnce`方法的细节逻辑问题；
1. 修复内部`empty`模块的`IsEmpty`方法对结构体属性的空校验错误；
1. 修复`gview`模板引擎的并发安全问题；
1. 修复`ghttp.Server`的SESSION初始化过期时间问题；

# `v1.6.0` (2019-04-09)

## 新功能/改进
1. `gcron`定时任务模块增加运行日志记录功能：https://goframe.org/os/gcron/index
1. `gredis`增加全局分组配置功能，并增加更多的配置选项`maxIdle/maxActive/idleTimeout/maxConnLifetime`：https://goframe.org/database/gredis/index
1. `gcfg`模块增加更多的默认配置文件检索路径，并且增加全局分组配置特性，增加`Instance`单例方法：https://goframe.org/os/gcfg/index
1. `gview`模块增加更多的默认配置文件检索路径，并且增加`Instance`单例方法：https://goframe.org/os/gview/index
1. `ghttp`模块新功能及改进：
    - 新增`CORS`HTTP(S)跨域请求特性: https://goframe.org/net/ghttp/cors
    - 增加`TLSConfig`配置功能；
    - 去掉路由注册方法的`error`返回值，当产生注册错误时直接终端打印错误/输出到日志文件；
    - 增加在`HTTP Code 302`跳转时的`Set-Cookie`支持；
    - 增加对`SESSION ID`的安全性检查；
    - 增加对基于`HTTPS`的`WebSocket`支持（`WSS`）：https://goframe.org/net/ghttp/websocket/index
    - `Request`对象增加`Error`方法，用于输出自定义错误信息到`WebServer`错误日志中；
    - 其他一些改进；
1. `gdb`模块新功能及改进：
    - 新增`Instance`单例管理方法；
    - 新增`Structs/Scan`链式操作方法，`gdb.DB/TX`新增`GetStructs/GetScan`方法，用于结果集`struct`/`slice`映射转换：https://goframe.org/database/gdb/chaining
    - 新增`Safe`链式操作方法（默认非并发安全），用于链式安全控制：https://goframe.org/database/gdb/chaining
    - `Where`链式操作方法改进：
        - 方法支持任意的`string/map/slice/struct/*struct`类型；
        - 逻辑调整，当链式操作中存在多个`Where`方法调用时，自动转换为`And`条件；
        - 支持`slice`条件参数，常用在`SELECT IN`查询中，例如：`Where("uid IN(?)", g.Slice{1,2,3})`；
        - 支持在`map`类型条件参数的`key`中传递条件，例如：`Where(g.Map{"uid>?", uid})`；
1. `gconv`及`gvalid`模块改进并去掉对私有`struct`方法属性的转换/校验；
1. `gconv.Map`转换方法新增对`json tag`: `-`, `omitempty`的支持: https://goframe.org/util/gconv/map
1. `gstr`模块新增 `ReplaceI/ReplaceIByArray/ReplaceIByMap`大小写非敏感替换方法；
1. `gutil`模块增加`IsEmpty`方法，用于判断给定变量是否为空（整型0, 布尔false, slice/map长度为0, 其他为nil的情况，判断为空），并增加快捷方法`g.IsEmpty`；
1. `gutil`模块增加`Export`方法，用于导出返回格式化打印的变量内容字符串，并增加快捷方法`g.Export`；
1. `gspath`增加缓存及非缓存检索检索方法`Search`/`SearchWithCache`；
1. `gjson`模块增加默认的`UseNumber`功能支持；
1. `gmap`增加`SetIfNotExistFunc/SetIfNotExistFuncLock`方法；
1. 迁移`greuseport`模块到新的仓库：https://github.com/gogf/greuseport
1. 大量的单元测试完善；

## Bug Fix
1. 修复`gqueue`模块的资源竞争问题；
1. 修复`gconv.GTime`转换失败问题；
1. 修复`gconv.String`在转换`int`参数时字节溢出问题；
1. 修复`ghttp.Request`的`HTTP Basic Auth`校验问题；
1. 修复`gxml`针对于非`UTF-8`编码内容转换的并发安全问题；
1. 修复`gtime`部分`Format`（`G`&`j`）格式失效问题；
1. 修复`gudp.Conn`对象的`RemoteAddr`获取客户端连接地址方法问题；
1. 修复`gmap/gcache`模块的`GetOrSetFuncLock`方法，增加对回调方法返回值的`nil`判断，只有非nil返回值才会被保存；



# `v1.5.8` (2019-02-28)

## 新特性
1. 主库从`gitee`迁移到了`github`( https://github.com/gogf/gf )，`gitee`作为镜像站，用于国内的代码贡献及ISSUE提交，迁移说明详见：https://goframe.org/upgradeto150
1. 对常用的`container`数组模块: `garray`做了大量改进/完善工作，新增大量常用方法，并完善单元测试用例及方法注释，详见API文档：https://godoc.org/github.com/gogf/gf/container/garray
1. 对常用的`container`集合模块: `gset`做了大量改进/完善工作，新增大量常用方法，并完善单元测试用例及方法注释，详见API文档：https://godoc.org/github.com/gogf/gf/container/gset
1. 对常用的`container`MAP模块: `gmap`做了大量改进/完善工作，新增大量常用方法，并完善单元测试用例及方法注释，详见API文档：https://godoc.org/github.com/gogf/gf/container/gmap
1. 对常用的字符串模块: `gstr`做了大量改进/完善工作，新增大量常用方法，并完善单元测试用例及方法注释，详见API文档：https://godoc.org/github.com/gogf/gf/text/gstr
1. 改进`gform`中对`struct`/`*struct`参数的支持，`*Insert/*Save/*Replace/*Update/Where/Data`方法的参数调整为`interface{}`类型，并支持任意类型的: `string/map/slice/struct/*struct`参数传递，具体请参考：https://goframe.org/database/orm/chaining
1. 新增/完善若干模块的单元测试用例, 包括：`gvalid`/`gregex`/`garray`/`gset`/`gmap`/`gstr`/`gconv`/`ghttp`/`gdb`；
1. 由于`gkafka`模块比较重，且不是框架核心模块，因此将该模块迁移到新的仓库中独立管理，并去掉相关依赖包：https://github.com/gogf/gkafka
1. 新增`greuseport`模块，用以实现TCP的`REUSEPORT`特性：https://godoc.org/github.com/gogf/gf/net/greuseport

## 新功能/改进
1. 去掉模板引擎内置变量中自动初始化`session`对象带来的内存占用问题；
1. `ghttp.Client`改进，增加若干方法，详见：https://goframe.org/net/ghttp/client 
1. `ghttp`分组路由增加`COMMON`方法，用以注册常用的`HTTP METHOD`(`GET/PUT/POST/DELETE`)路由；
1. 更新框架依赖的`golang.org/x/sys`模块；
1. 改进`gform`的批量操作(`Batch*`操作)返回结果对象，可以通过该结果对象获得批量操作准确的受影响记录行数；
1. 将`gstr`/`gregex`模块从`util`分类迁移到了`text`分类目录下；
1. 将`gtest`模块从`util`分类迁移到了`test`分类目录下；
1. 完善`glog`方法注释；

## Bug Fix
1. 修复带点的邮件格式,用`gvalid.Check`的"`email`"规则不能匹配成功;
1. 修复`gvalid.Check`在`regex`规则下的检查失败问题；
1. 修复`gcron`模块定时规则中天和周不允许`?`符号的问题；
1. 修复`ghttp.Server`在部分异常情况下仍然返回`200`状态码的问题；
1. 修复`gfpool`模块中由于原子操作问题造成的高并发"内存泄露"问题；
1. 修复分组路由注册对象/控制时，方法`Index`的路由仅能通过`/xxx/index`访问的问题；
1. 修复模板引擎使用中，当不存在`config.toml`(即使没使用)配置文件时的报错问题；
1. 其他一些修复；



# `v1.4.6` (2019-01-24)

## 新特性
1. 新增并发安全的高性能任务定时器模块`gtimer`, 类似于Java的`Timer`，但是比较于Java的`Timer`更加强大，内部实现采用灵活高效的`分层时间轮`设计，被设计为可管理维护百万级别以上数量的定时任务。`gtimer`为`GF`框架的核心模块之一，单元测试覆盖率达到`93.6%`：[https://goframe.org/os/gtimer/index](https://goframe.org/os/gtimer/index)
1. 采用任务定时器`gtimer`重构`gcron`定时任务模块，去掉第三方`github.com/robfig/cron`包的使用。`gcron`增加单例模式的定时任务：[https://goframe.org/os/gcron/index#](https://goframe.org/os/gcron/index#)；
1. `gconv`类型转换模块支持对`struct`结构体中的**指针属性**转换：[https://goframe.org/util/gconv/struct](https://goframe.org/util/gconv/struct)；
1. `gform`增加对数据库类型的自动识别特性，这一特性在需要将查询结果`json`编码返回时非常有用: [https://goframe.org/database/orm/index](https://goframe.org/database/orm/index)
1. `Travis CI`增加对`386`架构的自动化测试支持(目前已支持`386`和`amd64`)；

## 新功能
1. `ghttp`模块新增`Exit`、`ExitAll`、`ExitHook`方法，用于HTTP请求处理流程控制: [https://goframe.org/net/ghttp/service/object](https://goframe.org/net/ghttp/service/object)；
1. `grand`模块增加`Meet/MeetProb`方法，用于给定概率的随机满足判断，增加别名方法`N/Str/Digits/Letters`；
1. `gvalid`数据/表单校验模块增加`16X`及`19X`手机号的校验支持；

## 功能改进
1. `gform`设置默认的数据库连接池`CONN_MAX_LIFE`参数值为`30`秒；
1. 改进`glist`模块，提高约`20%`左右性能，并增加若干链表操作方法；
1. 改进`gqueue`模块，提高约`50`左右性能，并增加模块对`select`语法的支持(使用`Queue.C`): [https://goframe.org/container/gqueue/index](https://goframe.org/container/gqueue/index)；
1. 改进`gmlock`内存锁模块，并完善单元测试用例：[https://goframe.org/os/gmlock/index](https://goframe.org/os/gmlock/index)；
1. 改进并发安全容器所有的模块，调整并发安全控制非必需参数`safe...bool`为`unsafe...bool`；
1. 改进`gpool`对象复用模块，支持并发安全；
1. 更新`gkafka`模块的第三方依赖包；
1. 完善`ghttp`模块的单元测试用例；


## Bug Fix
1. 修复`gmd5`模块操作文件时的文件指针未关闭问题；
1. 修复`gcache`缓存项过期删除失效问题；
1. 其他修复；

# `v1.3.8` (2018-12-26)

## 新特性
1. 对`gform`完成重构，以提高扩展性，并修复部分细节问题、完善单元测试用例([https://goframe.org/database/orm/index](https://goframe.org/database/orm/index))；
1. `WebServer`路由注册新增分组路由特性([https://goframe.org/net/ghttp/group](https://goframe.org/net/ghttp/group));
1. `WebServer`新增`Rewrite`路由重写特性([https://goframe.org/net/ghttp/static](https://goframe.org/net/ghttp/static));
1. 增加框架运行时对开发环境的自动识别；
1. 增加了`Travis CI`自动化构建/测试；

## 新功能
1. 改进`WebServer`静态文件服务功能，增加`SetStaticPath`/`AddStaticPath`方法([https://goframe.org/net/ghttp/static](https://goframe.org/net/ghttp/static))；
1. `gform`新增`Filter`链式操作方法，用于过滤参数中的非表字段键值对（[https://goframe.org/database/orm/linkop](https://goframe.org/database/orm/linkop)）；
1. `gcache`新增`Data`方法，用以获取所有的缓存数据项；
1. `gredis`增加`GetConn`方法获取原生redis连接对象；

## 功能改进
1. 改进`gform`的`Where`方法，支持`slice`类型的参数，并更方便地支持`in`操作查询（[https://goframe.org/database/orm/linkop](https://goframe.org/database/orm/linkop)）；
1. 改进`gproc`进程间通信数据结构，将`pid`字段从`16bit`扩展为`24bit`；
1. 改进`gconv`/`gmap`/`garray`，增加若干操作方法；
1. 改进`gview`模板引擎中的`date`内置函数，当给定的时间戳为空时打印当前的系统时间；
1. 改进`gview`模板引擎中，当打印的变量不存在时，显示为空（标准库默认显示为`<no value>`）；
1. 改进`WebServer`，去掉`HANGUP`的信号监听，避免程序通过`nohup`运行时产生异常退出问题；
1. 改进`gcache`性能，并完善基准测试；

## Bug Fix
1. 修复`gcache`在非LRU特性开启时的缓存关闭资源竞争问题，并修复`doSetWithLockCheck`内部方法的返回值问题；
1. 修复`grand.intn`内部方法在`x86`架构下的随机数位溢出问题；
1. 修复`gbinary`中`Int`方法针对`[]byte`参数长度自动匹配造成的字节长度溢出问题；
1. 修复`gjson`由于官方标准库`json`不支持`map[interface{}]*`类型造成的Go变量编码问题；
1. 修复`garray`中部分方法的数据竞争问题，修复二分查找排序问题；
1. 修复`ghttp.Request.GetVar`方法获取参数问题；
1. 修复`gform`的数据库连接池不起作用的问题；





# `v1.2.11` (2018-11-26)
## 新特性
1. `ORM`新增对`SQLServer`及`Oracle`的支持([https://goframe.org/database/orm/database](https://goframe.org/database/orm/database))；
1. 完成`gvalid`模块校验结果的顺序特性([https://goframe.org/util/gvalid/checkmap](https://goframe.org/util/gvalid/checkmap));
1. 改进`ghttp.Request.Exit`，使得调用该方法时立即退出业务执行，开发者无需调用`Exit`方法时再使用`return`返回([https://goframe.org/net/ghttp/service/object](https://goframe.org/net/ghttp/service/object))；
1. 模板引擎新增若干内置函数:`text/html/htmldecode/url/urldecode/date/compare/substr/strlimit/hidestr/highlight/toupper/tolower/nl2br` ([https://goframe.org/os/gview/funcs](https://goframe.org/os/gview/funcs));
1. 模板引擎新增内置变量`Config` ([https://goframe.org/os/gview/vars](https://goframe.org/os/gview/vars));
1. 改进`gconv.Struct`转换默认规则，支持不区分大小写的键名与属性名称匹配；
1. `gform`配置文件支持`linkinfo`自定义数据库连接字段([https://goframe.org/database/orm/config](https://goframe.org/database/orm/config))；
1. `gfsnotify`模块增加对特定回调的取消注册功能（[https://goframe.org/os/gfsnotify/index](https://goframe.org/os/gfsnotify/index)）；



## 新功能
1. 改进`ghttp.Request`，增加`SetParam/GetParam`请求流程自定义变量设置/获取方法，用于在请求流程中的回调函数共享变量（[https://goframe.org/net/ghttp/request](https://goframe.org/net/ghttp/request)）;
1. 改进`ghttp.Response`，增加`ServeFileDownload`方法，用于WebServer引导客户端下载文件([https://goframe.org/net/ghttp/response](https://goframe.org/net/ghttp/response));
1. `gvar`模块新增`gvar.VarRead`只读接口，用于控制对外只暴露数据读取功能；
1. 增加`g.Throw`抛异常方法，`g.TryCatch`异常捕获方法封装;
1. 改进`gcron`模块，增加自定义的Cron管理对象，增加`New/Start/Stop`方法;


## 功能改进
1. WebServer添加`RouterCacheExpire`配置参数，用于设置路由检索缓存过期时间；
1. WebServer允许同一`HOOK`事件被多次绑定注册，先注册的回调函数优先级更高([https://goframe.org/net/ghttp/service/hook](https://goframe.org/net/ghttp/service/hook));
1. 当前工作目录为系统临时目录时，`gcfg`/`gview`/`ghttp`模块默认不添加工作目录到搜索路径;
1. 改进`WebSocket`默认支持跨域请求([https://goframe.org/net/ghttp/websocket](https://goframe.org/net/ghttp/websocket));
1. 改进`gtime.Format`支持中文;
1. 改进`gfsnotify`，支持编辑器对文件非执行标准编辑时(RENAME+CHMOD)的热更新问题;
1. 改进`gtype.Set`方法，增加Set原子操作返回旧的变量值;
1. `gfile.ScanDir`增加支持`pattern`多个文件模式匹配，使用'`,`'符号分隔多个匹配模式;
1. `gcfg`模块增加获取配置变量为`*gvar.Var`;
1. `gstr`模块增加对中文截取方法;
1. 改进`gtime.StrToTime`对常用时间格式匹配模式，新增`gtime.ParseTimeFromContent`方法;
1. 修改配置管理、模板引擎、调试模式的环境变量名称为大写下划线标准格式;
1. 改进`grand`模块随机数生成设计，底层使用`crypto/rand`+缓冲区实现高速的随机数生成([https://goframe.org/util/grand/index](https://goframe.org/util/grand/index));

## 问题修复
1. 修复`gspath`模块在`windows`下搜索失效问题；
1. 修复`gspath`模块Search时带有indexFiles的检索问题;
1. bug fix INZS1([https://github.com/gogf/gf/issues/INZS1](https://github.com/gogf/gf/issues/INZS1));
1. 修复`gproc.ShellRun`在windows下的执行问题;






# `v1.0.898 stable` (2018-10-24)

## 新特性
1. `gf-orm`增加`sqlite`数据库类型支持([http://gf.johng.cn/database/orm/database](http://gf.johng.cn/database/orm/database))；
1. 增加`gkafka`模块，对kafka的客户端程序封装，支持分组消费及指定起始位置等特性，并提供简便易用的API接口([http://gf.johng.cn/database/gkafka/index](http://gf.johng.cn/database/gkafka/index))；
1. 增加go语言最新版本的`go modules`特性支持；
1. 增加`gcron`定时任务模块([http://gf.johng.cn/os/gcron/index](http://gf.johng.cn/os/gcron/index))；
1. `Web Server`增加路由注册项获取/打印特性，所有的路由注册/回调注册一览无余；
1. 模板引擎增加全局变量管理，并增加多个常用的内置函数及内置变量（[http://gf.johng.cn/os/gview/funcs](http://gf.johng.cn/os/gview/funcs)）；
1. `gredis`改进为单例操作方式(基于基层连接池特性)，每次操作`redis`服务器时开发者无需显示调用`Close`方法执行关闭（[http://gf.johng.cn/database/gredis/index](http://gf.johng.cn/database/gredis/index)）；
1. `gf-orm`增加数据库操作自动`Close`特性(基于底层链接池特性)，开发者无需再`defer db.Close()`，并增加`g.DB`数据库对象单例别名（[http://gf.johng.cn/database/orm/linkop](http://gf.johng.cn/database/orm/linkop)）；
1. 增加`gvar`通用动态变量模块（[http://gf.johng.cn/container/gvar/index](http://gf.johng.cn/container/gvar/index)）；
1. 数据结构容器增加`并发安全特性开启/关闭功能`，当关闭后和普通的数据结构无异，且在非并发安全模式下性能会得到提高;
1. 新增`gmlock`内存锁模块（[http://gf.johng.cn/os/gmlock/index](http://gf.johng.cn/os/gmlock/index)）；
1. 增加`gaes`算法模块（[http://gf.johng.cn/crypto/gaes/index](http://gf.johng.cn/crypto/gaes/index)）；
1. `gproc`模块增加执行`shell`命令方法（[http://gf.johng.cn/os/gproc/index](http://gf.johng.cn/os/gproc/index)）；
1. 新增`gfcache`模块，用于带自动缓存更新的文件内容操作（文档待完善）；

## 新功能
1. `glog`增加链式操作方法，增加日志级别管理控制、分类管理、调试管理功能；
1. `g.View`增加分组名称设置，支持通过`g.*`对象管理器获取多个命名的单例模板引擎对象；
1. `glog`增加对文件名称格式的自定义设置，支持`gtime日期格式`；
1. `gconv`增加`Ints/Uints/Floats/Interfaces`转换方法;
1. `gjson`增加`Append`方法;
1. `gparser`增加`NewUnsafe/Append`方法;
1. `gcache`增加`GetOrSet/GetOrSetFunc/GetOrSetFuncLock`方法；
1. `gset`增加`LockFunc/RLockFunc`方法;
1. `ghttp.Response`方法完善，增加`ParseTpl/ParseTplContent/TplContent`方法，`Template`修改为`Tpl`方法;
1. `ghttp.Request`增加获取用户真实IP判断;
1. `Session`增加`Contains`方法;
1. 完善`ghtml`模块，增加多个方法;
1. `gcache`新增`Contains/SetIfNotExist`方法；
1. `gvalid`增加`Error`对象，用以管理校验错误信息；
1. `gvalid`模块增加`struct tag`的校验规则、自定义错误提示信息绑定的支持特性（[http://gf.johng.cn/util/gvalid/index](http://gf.johng.cn/util/gvalid/index)）；
1. `ghttp`增加输入参数与`struct`的`绑定机制`，并增加对应`params`标签支持（[http://gf.johng.cn/net/ghttp/service/handler](http://gf.johng.cn/net/ghttp/service/handler)）；
1. `ghttp.Request`增加服务端`BasicAuth`功能(文档待完善)；
1. `gvalid`增加字段校验别名，用于自定义返回结果字段，并更新WebServer中相关使用的模块；
1. `gf-orm`链式操作增加`ForPage`方法，调整`Chunks`方法;
1. `ghttp`对象路由注册增加`Init&Shut`自动回调方法，增加重复路由注册检测功能；
1. `gfsnotify`增加默认递归`Add/Remove`特性；
1. `ghttp.Response`增加`ServiceFile`方法；
1. 其他一些新功能；

## 功能改进
1. 改进`ghttp.Server`配置管理；
1. 改进`gcache`底层对象继承关系，改进部分设计细节，提高性能；
1. 改进`gfpool`文件指针池，修复部分错误，提升性能，并增加基准测试代码；
1. 改进`gmap`系列并发安全map数据结构，增加多个易用性的方法；
1. 改进`gconv.Struct`对象转换功能（[http://gf.johng.cn/util/gconv/index](http://gf.johng.cn/util/gconv/index)）；
1. 改进`grand`随机数生成规则，提供了极高的随机数生成性能，并保证每一次调用随机方法时生成的都是不同的随机数值（[http://gf.johng.cn/util/grand/index](http://gf.johng.cn/util/grand/index)）；
1. 改进`gfile`文件内容操作方法，增加若干常用的文件内容读取方法；
1. 改进`gtime`模块，并增加时区转换方法；
1. 改进`COOKIE`，去掉锁机制；
1. 改进`SESSION`获取方法，新增多个类型获取方法；
1. 改进`g.DB/g.Config`单例缓存键名；
1. 改进`gtcp/gudp`超时错误判断机制;
1. 改进`gtype`底层统一修改为原子操作;
1. 改进`gvalid`对`struct`的`string`属性的默认值非必需校验;
1. 改进`gvalid`在关联规则下的非必需校验；
1. 改进`gf-orm`在调试模式下日志自动输出功能；
1. `ghttp.Server/gspath`模块静态文件检索改进;
1. 优化`ghttp.ServerConfig`配置，增加`struct/method``名称到uri`的转换规则，通过`SetNameToUri`方法进行灵活配置（[http://gf.johng.cn/net/ghttp/service/object](http://gf.johng.cn/net/ghttp/service/object)）；
1. 改进`*any/:name`路由匹配规则，支持不带名字的`*/:`路由规则;
1. 修改默认配置文件名称 `config.yml` -> `config.toml`（[http://gf.johng.cn/os/gcfg/index](http://gf.johng.cn/os/gcfg/index)）；
1. 调整服务注册的`BindControllerMethod`及`BindObjectMethod`逻辑为绑定路由到指定的方法执行；
1. 改进`garray`二分查找方法，增加安全操作处理；
1. 改进`gdb.Result/Recorde` `ToXml`方法，增加可选的`rootTag`参数；
1. 其他一些改进；

## 问题修复
1. 修复`ghttp.Server`在`windows`下的重启失效问题；
1. 修复`ghttp.Server`服务注册与回调注册路由重复判断问题；
1. 修复`garray`排序数组`Add`变参时的死锁问题；
1. 修复`gfsnotify`默认递归监控整个`gspath.Add`添加的目录的问题;
1. 修复`ghttp.BindParams`对`@file`文件上传标识符的转义问题；
1. 修复`ghttp.Server`日志路径丢失问题;
1. 修复`多WebServer`下的状态检测问题;
1. 修复`gvalid`模块`min/max`校验问题；
1. 修复控制器和执行对象服务注册时绑定'/'路由的问题；
1. 修复`gvalid.CheckStruct`自定义错误提示失效问题；
1. `ghttp.Server`修复`hook`与`serve`方法的路由影响，并新增跳转方法；
1. 其他一些修复；

## 其他改动
1. 去掉`gfile.IsExecutable`方法；
1. 目录调整，将`加密/解密`相关的包从`encoding`目录迁移到`crypto`目录下；
1. 增加`gfsnotify/gfcache`调试信息；
1. `gf-orm`允许写入的键值为`nil`时往数据库中写入`null`;
1. 统一使用`gview.Params`类型作为模板变量类型；
1. `gconv.MapToStruct`方法名称修改为`gconv.Struct`；
1. `ghttp.Server`完善重启及停止的终端提示信息；
1. 完善`gring`模块，增加`约瑟夫问题`代码作为`gring`示例程序；
1. 其他一些改动；



# `v0.99.682 beta` (2018-08-07)
## 新特性
    1、新增gdes包，用于DES加密/加密算法处理；
    2、新增gkafka包，kafka的golang客户端；
    3、新增gpool对象复用池，比较于标准库的sync.Pool更加灵活强大，可自定义对象的缓存时间、创建方法、销毁方法(http://gf.johng.cn/686654)；
    4、完成网络通信gtcp/gudp包的重构，并进行了大量的改进工作，新增了详尽的开发文档及示例代码(http://gf.johng.cn/494382)；
    5、增加gring并发安全环，标准库container/ring包的并发安全版本，并做了易用性的封装(http://gf.johng.cn/686655)；
    6、gtime包新增了自定义日期格式话的支持，格式化语法类似PHP的date语法(http://gf.johng.cn/494387)；
    7、gdb增加调试模式特性，使用SetDebug方法实现，在调试模式下可以获取详细的SQL执行记录，增加了详细的开发文档及示例代码(http://gf.johng.cn/702801)；
    8、gdb增加查询缓存特性，使用Cache方法实现，增加了详细的开发文档及示例代码(http://gf.johng.cn/702801)；
    9、ghttp.Server路由功能增加字段匹配规则特性，支持如：/order/list/{page}.html 动态路由规则特性(http://gf.johng.cn/702766)；
    10、gpage分页包增加分页URL规则生成模板特性，内部可使用{.page}变量指定页码位置(http://gf.johng.cn/716438)；
    11、增加gmap.Map对象，这是gmap.InterfaceInterfaceMap的别名；

## 新功能
    1、gdb增加MaxIdleConnCount/MaxOpenConnCount/MaxConnLifetime三项配置，并增加SetMaxConnLifetime方法；
    2、ghttp.Client增加HTTP账号密码设置功能(SetBasicAuth)；
    3、glog新增对系统换行符号的自适配调整(\n|\r\n)；
    4、增加glog控制台调试模式打印开关(SetDebug)；
    5、gcfg增加SetFileName方法设置默认读取的配置文件名称；
    6、gcfg/gjson/gparser包新增Int8/16/32/64,Uint8/16/32/64方法；
    7、增加gzip方法的封装(Zip/Unzip)；
    8、gview增加模板变量分隔符设置方法SetDelimiters；
    9、ghttp.Response增加Writef、Writefln方法；

## 功能改进
    1、改进gfilepool文件指针池设计；改进gfile文本内容写入，增加指针池使用
    2、gdb包增加调试模式特性，并支持在调试模式下获得已执行的SQL列表结果
    3、改进gproc进程间通信机制，增加进程消息分组特性，并限定队列大小
    4、gdb结果方法处理增加ToXml/ToJson方法
    5、gregx包名修改为gregex
    6、改进gtime.StrToTime方法，新增对常见标准时间日期的自动转换，以及对时区的自动识别支持，并调整gconv,gvalid对该包的引用
    7、增加对字符集转换的封装，gxml包中使用新增的字符集转换包来做处理
    8、ghttp.Server.EnableAdmin页面Restart接口支持GET参数newExeFilePath支持
    9、ghttp.Server平滑重启机制增加可自定义重启可执行文件路径，特别是针对windows系统特别有用(因为windows下不支持可执行文件覆盖更新)
    10、改进ghttp.Server静态文件检索设计，增加开发环境时的main包源码目录查找机制；改进gcfg/gview的main包源码目录查找机制
    11、优化gcache设计，LRU特性非默认开启；优化gtype/gcache基准测试脚本；新增gregx基准测试脚本，改进设计，提升性能
    12、gfile包增加GoRootOfBuild方法，用于获取编译时的GOROOT数值；并改进glog包中backtrace的GOROOT路径过滤处理；
    13、改进grpool代码质量，并改进对池化goroutine数量的限制设计
    14、改进gdb.Map/List及g.Map/List的类型定义，改用别名特性以便支持原生类型输入(map/slice)，并修复gdb.Model.Update方法参数处理问题
    15、调整ghttp包示例代码目录结构，增加ghttp.Client自定义Header方法，ghttp.Cookie增加Map方法用于获得客户端提交的所有cookie值，构造成map返回
    16、删除gcharset中的getcharset方法
    17、去掉gmap中常用的基本数据类型转换获取方法
    18、改进gconv.String方法，当无法使用基本类型进行字符串转换时，使用json.Marshal进行转换
    19、gvalid.CheckObject方法名称修改为gvalid.CheckStruct


## 问题修复
    1、修正gstr.IsNumeric错误
    2、修复当xml中encoding字符集为非UTF-8字符集时报错的问题
    3、修正gconv包float32->float64精度问题
    4、修复gpage包分页计数问题
    5、修复gdb批量数据Save错误
    6、去掉gpool中math.MAXINT64常量的使用，以修复int64到int类型的转换错误，兼容32位系统
    7、修正ghttp包没有使用Server仍然初始化相关异步goroutine的问题




# `v0.98.503 beta` (2018-05-21)
## 新特性
    1、平滑重启特性( http://gf.johng.cn/625833 )；
    2、gflock文件锁模块( http://gf.johng.cn/626062 )；
    3、gproc进程管理及通信模块( http://gf.johng.cn/626063 )；
    4、gpage分页管理模块，强大的动态分页及静态分页功能，并为开发者自定义分页样式提供了极高的灵活度( http://gf.johng.cn/597431 )；
    5、ghttp.Server增加多端口监听特性，并支持HTTP/HTTPS( http://gf.johng.cn/494366 , http://gf.johng.cn/598802 )；
    6、增加gspath目录检索包管理工具，支持对多目录下的文件检索特性；
    7、ghttp包控制器及执行对象注册增加更灵活的动态路由特性，路由规则增加{method}变量支持；

## 新功能
    1、gutil包增加MapToStruct方法，支持将map数据类型映射为struct对象；
    2、gconv
        1)、gconv包增加按照类型名称字符串进行类型转换；
        2)、gconv包新增Time/TimeDuration类型转换方法；
    3、ghttp
        1)、增加Web Server目录安全访问控制机制；
        2)、ghttp.Server增加自定义状态码回调函数注册处理；
    4、gdb
        1)、gdb包增加gdb.GetStruct/gdb.Model.Struct方法，获取查询结果记录自动转换为指定对象；
        2)、gdb增加Value/Record/Result类型，增加对Value类型的系列类型转换方法；
        3)、gdb包增加db.GetCount,tx.GetCount,model.Count数量查询方法；

## 功能改进
    1、改进gredis客户端功能封装；
    2、改进grand包随机数生成性能；
    3、grand/gdb/gredis包增加benchmark性能测试脚本；
    4、改进gjson/gparser包的ToStruct方法实现；
    5、gdb  ：改进gdb.New获取ORM操作对象性能；
    6、gcfg ：改进配置文件检索功能；
    7、gview：模板引擎增加多目录检索功能；
    8、gfile：增加源码main包目录获取方法MainPkgPath；
    9、ghttp
        1)、ghttp.Request增加请求进入和完成时间记录，并增加到默认日志内容中；
        2)、ghttp.Server事件回调之间支持通过ghttp.Request.Param自定义参数进行流程传参；
    10、gdb
        1)、改进gdb.Result与gdb.List, gdb.Record与gdb.Map之间的类型转换，便于业务层数据编码处理(如json/xml)；
        2)、改进gdb.Tx.GetValue返回值类型；
        3)、gdb.Model.Data参数支持更加灵活的map参数；

## 问题修复
    1、ghttp
        1)、修复ghttp包路由缓存问题；
        2)、修复服务注册时的控制器及执行对象方法丢失问题；
    2、gconv
        1)、修正gconv.Float64方法位大小设置问题；
        2)、修复gconv.Int64(float64(xxx))问题；
    2、gdb
        1)、修复gdb.GetAll针对返回数据列表的for..range...的返回结果slice相同指针问题；
        2)、修复gdb.Delete方法错误；
        3)、修复gdb.Model.And/Or方法；
        4)、修复gdb.Model.Where方法参数处理问题；
    3、garray：修复garray包Remove方法锁机制问题；
    4、gtype ：修复gtype.Float32/gtype.Float64对象类型的方法逻辑错误；
    5、gfsnotify：修复在windows下文件参数中不同文件分隔符引起的热更新机制失效问题；
    6、修复gvalid包验证问题：如果值为nil，并且不需要require*验证时，其他验证失效。并增加单元测试项，测试通过。





# `v0.97.399 beta` (2018-04-23)
    1、 增加gfsnotify文件监控模块；
    2、 配置管理模块增加配置文件自动检测更新机制；
    3、 模板引擎增加对模板文件的自动检测更新机制；
    4、 改进gconv包基本类型转换功能，提高转换性能；
    5、 增加gpage分页管理包，支持动态分页、静态分页以及自定义分页样式特性；
    6、 ghttp.Request增加Exit方法，用以标记服务退出，当在服务执行前调用后，服务将不再执行；
    7、 ghttp.Response去掉WriteString方法，统一使用Write方法返回数据流，是使用灵活的参数形式；
    8、 模板引擎增加模板变量暴露接口LockFunc/RLockFunc，以便支持开发者灵活处理模板变量；
    9、 ghttp.Server增加access & error log功能，并支持开发者自定义日志处理回调函数注册；
    10、增加gredis包，支持对redis的客户端操作封装，并将gredis.Redis对象加入到gins单例管理器中进行统一配置管理维护；
    11、gins单例管理器增加对单例对象配置文件的自动检测更新机制，当配置文件在外部发生变更时，自动刷新单例管理器中的单例对象；
    12、gdb数据库ORM包增加And/Or条件链式方法，并改进Where/Data方法参数灵活性；
    13、对于新增加的模块，同时也增加了对应的开发文档，并梳理完善了现有的其他模块开发文档；
    14、修复ISSUE:
        #IISWI github.com/gogf/gf/issues/IISWI,
        #IISMY github.com/gogf/gf/issues/IISMY,
        反馈并跟踪完成第三方依赖mxj包的ISSUE修复(github.com/clbanning/mxj/issues/48)；



