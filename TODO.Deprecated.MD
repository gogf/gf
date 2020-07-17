> This markdown is deprecated and will be removed in future. All TODO features are staged in the issue: https://github.com/gogf/gf/issues

# ON THE WAY
1. 增加图形验证码支持，至少支持数字和英文字母；
1. Cookie&Session数据池化处理；
1. ghttp.Client增加proxy特性；
1. gtime增加对时区转换的封装，并简化失去转换时对类似+80500时区的支持；
1. orm增加sqlite对Save方法的支持(去掉触发器语句);
1. ghttp.Server增加Ip访问控制功能(DenyIps&AllowIps)；
1. ghttp增加返回数据压缩机制；
1. ghttp.Server增加proxy功能特性，本地proxy和远程proxy，本地即将路由规则映射；远程即反向代理；
1. gjson对大json数据的解析效率问题；
1. ghttp增加route name特性，并同时支持backend和template(提供内置函数)引用，可以通过RedirectRoute方法给定route name和路由参数跳转到指定的路由地址上；
1. gvalid校验支持当第一个规则失败后便不再校验后续的规则，最好做成链式操作；
1. ghttp.Request增加对输入参数的自动HtmlEncode机制；
1. 常量命名风格根据golint进行修改；
1. 开放rwmutex包，并将gjson的互斥锁使用自定义的mutex替换；
1. 文档完善：
    - gconv struct tag、
    - 控制器及执行对象注册的Init&Shut方法、
    - ghttp.Response&ServeFile、gfcache、gproc shell执行、
    - ghttp Server&Client basic auth、
    - glog分类&日志等级&链式操作、gdb debug自动输出调试信息、gmlock内存锁、
1. 服务注册域名增加对泛域名的支持；
1. 项目参考：
    - https://github.com/namreg/godown
    - https://github.com/Masterminds/sprig
1. gform参考 https://gohouse.github.io/gorose/dist/index.html 进行改进
1. gtcp提供简便的包发送/接收方法(SendPkg/RecvPkg)以解决常见的TCP通信粘包问题，并完善文档（参考：https://www.cnblogs.com/kex1n/p/6502002.html）；
1. 路由增加不区分大小写得匹配方式；
1. 改进WebServer获取POST参数处理逻辑，当提交非form数据时，例如json数据，针对某些方法可以直接解析；
1. WebServer增加可选择的路由覆盖配置，默认情况下不覆盖；
1. 增加jumplist的数据结构容器；
1. DelayQueue/PriorityQueue；
1. 权限管理模块；
1. 从ghttp中剥离SESSION功能构成单独的模块gsession；
1. 改进gproc进程间通信处理逻辑，提高稳定性，以应对进程间大批量的数据发送/接收；
1. 添加Save/Replace/BatchSave/BatchReplace方法对sqlite数据库的支持；
1. 添加sqlite数据库的单元测试用例；
1. gredis增加cluster支持；
1. gset.Add/Remove/Contains方法增加批量操作支持；
1. gmlock增加手动清理机制：当内存锁不再使用时，由调用端决定是否清理内存锁；
1. gtimer增加DelayAdd*方法返回Entry对象，以便DelayAdd*的定时任务也能进行状态控制；gcron同理需要改进；
1. grpool增加支持阻塞添加任务接口；
1. gdb.Model在链式安全的对象创建中增加sync.Pool的使用；
1. 增加g.Table快捷方法以方便操作数据表，但是得考虑后续模型操作设计，特别是脚手架的模型管理；






# DONE
1. Cookie设置中文失效问题；
1. gvalid增加支持对[]rune的长度校验(一个中文占3个字节)；
1. grpool性能压测结果变慢的问题；
1. ghttp的热重启的本地进程端口监听，在不使用该特性时默认关闭掉；
1. gtcp增加对TLS加密通信的支持；
1. 改进gdb对pgsql/mssql/oracle的支持，使用方法覆盖的方式改进操作，而不是完全依靠正则替换的方式；
1. gdb的Cache缓存功能增加可自定义缓存接口，以便支持外部缓存功能，缓存接口可以通过io.ReadWriter接口实现；
1. 改进ghttp分组路由中对hook的支持方式，以便格式与BindHookHandler统一；
1. 使用gconv将slice映射到struct属性上，例如redis hscan的结果集；
1.  gconv完善针对不同类型的判断，例如：尽量减少sprintf("%v", xxx)来执行string类型的转换；
2.  ghttp.Server请求执行中增加服务退出的方法，不再执行后续操作；
3.  ghttp.Response对象完善并改进数据返回方法(Write/WriteString)；
4.  ghttp.Server请求执行中增加服务退出的方法，不再执行后续操作；
5.  增加fsnotify包支持；
6.  改进gcfg和gview的文件自动更新机制；
7.  将模板变量进行暴露，以便应用端可以进行灵活控制；
8.  跟踪第三方mxj包的issue问题：https://github.com/clbanning/mxj/issues/48；
9.  gdb Where方法参数的改进，研究是否可以将string参数类型修改为interface{}；
10. gpage分页控制功能；
11. https支持；
12. ghttp.Server日志中增加请求时间和返回时间，以便计算执行时间差；
13. 由于去掉了gdb的单例模式，并且将gins的部分对象封装迁移到了g包中，需要同时梳理文档，完善修改；
14. 在代码中增加https与http同时开启使用的示例代码，这块大家问得比较多；
15. ghttp.Server多个事件之间通过ghttp.Request.Param自定义参数传参；
16. 研究是否增加配置文件目录检索功能，特别是如何友好改进开发环境的配置文件默认目录问题；
17. 增加ghttp.Server不同状态码的自定义处理方法；
18. ghttp.Server平滑重启方案；
19. 完善gconv类型转换功能，增加time.Time/time.Duration类型转换，并增加benchmark测试脚本
20. 当二进制参数为nil时，gjson.LoadContent并将gjson.Json对象ToMap时会报错；
21. 改进控制器及执行对象注册，更友好地支持动态路由注册，例如：注册规则为 /channel/:name，现有的控制器及执行对象注册很难友好支持这种动态形式；
22. 当前gpage分页包的输出标签不支持li，大多数CSS框架都是li+a标签模式，需要提供可更加灵活的定制化功能实现；
23. 平滑重启机制改进，以便于开发阶段调试；
24. 对grpool进行优化改进，包括属性原子操作封装采用gtype实现，修正设计BUG：https://github.com/gogf/gf/issues/6；
25. gredis增加redis密码支持；
26. 改进ghttp.Server平滑重启机制，当新进程接管服务后，再使用进程间通信方式通知父进程销毁；
27. gproc进程间通信增加分组特性，不同的进程间可以通过进程ID以及分组名称发送/获取进程消息；
28. ORM增加获取被执行的sql语句的方法；
29. gdb增加查询缓存特性；
30. gpage分页增加对自定义后缀的支持，如:2.html, 2.php等等；
31. gvalid包增加struct tag的校验规则、自定义错误提示信息绑定的支持特性；
32. 增加文件缓存包，可根据fsnotify机制进行缓存更新；
33. *any/:name路由匹配路由改进支持不带名字的*/:路由规则；
34. ghttp静态文件服务改进(特别是403返回状态的修改)；
35. map转struct增加对tag的支持；
36. gcache检查在i386下的int64->int转换问题；
37. ghttp获取参数支持直接转struct功能；
38. gfsnotify增加对于目录的监控；
39. 检查windows下的平滑重启失效问题；
40. ghttp.Server的Cookie及Session锁机制优化(去掉map锁机制);
41. 解决glog串日志情况；
42. glog增加对日志文件名称的生成规则设定，支持时间格式规则；
43. ghttp日志增加客户端IP信息；
44. 完善gform配置管理说明，g.DB/Database和gdb.New的区别；
1. 完善配置管理章节，说明默认的配置文件更改方式；
1. 服务注册时判断方法定义满足规范时才执行绑定，否则提示WARN信息；
1. `gfsnotify`增加添加监听文件时的监听ID返回，以便调用端删除监听时只删除自己添加的监听，而不影响其他对该同一文件的监听回调；
1. `gfsnotify`针对添加目录监听时无法使用多个`Watcher`，考虑改进，并考虑动态扩容全局`Watcher`方案；
1. 由于系统对inotify实例数量(`fs.inotify.max_user_instances`)以及队列大小(`fs.inotify.max_user_watches`)有限制，需要改进`gfsnotify`；
1. WebServer事件回调允许对同一个路由规则绑定多个事件回调；
1. gcfg/gview/ghttp等模块加上对临时文件目录的自动添加监听判断（基本是开发环境下，特别是windows环境），去掉临时文件的监听，避免临时文件过大引起的运行缓慢占用内存问题；
1. 改进gfpool在文件指针变化时的更新；
1. ghttp hook回调使用方式在注册路由比较多的时候，优先级可能使得开发者混乱，考虑方式便于管理；
1. gform对于MySQL字段类型为datetime类型的时区问题分析；
1. 改进证书打开失败时的WebServer错误提示，前置HOOK校验后关闭后续的HOOK逻辑执行；
1. 目前WebServer的HOOK是按照优先级执行的，需要增加覆盖特性；
1. 更新跨域请求CORS相关功能文档；
1. ghttp.Response增加输出内容后自动退出当前请求机制，不需要用户手动return，参考beego如何实现；
1. gcfg包目前允许添加重复的目录路径，需要在SetPath/AddPath时判断重复性，不能添加重复的路径；
1. gdb执行数据写入时，如果参数为struct/[]struct，自动映射与表字段对应关系，不再使用gconv标签标识；
1. gdb的Data方法支持struct参数传入；
1. gfcache依旧使用gcache作为缓存控制对象，不要使用gmap；
1. 增加对ghttp路由注册的{.struct}/{.method}单元测试；
1. gconv针对struct的转换增加json tag支持，gconv.Map默认也支持json tag, 完善开发文档；
1. 增加SO_REUSEPORT的支持；
1. gkafka这个包比较重，未来从框架中剥离出来；
1. str_ireplace: http://php.net/manual/en/function.str-ireplace.php
1. strpos/stripos/strrpos/strripos: http://php.net/manual/en/function.stripos.php
1. gfile对于文件的读写强行使用了gfpool，在某些场景下不合适，需要考虑剥离开，并为开发者提供单独的指针池文件操作特性；
1. ghttp.Client自动Close机制；
1. ghttp路由功能增加分组路由特性；
1. 增加可选择性的orm tag特性，用以数据表记录与struct对象转换的键名属性映射；
1. gview中的template标签失效问题；