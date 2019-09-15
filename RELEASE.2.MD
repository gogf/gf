# `v2.0.0`

## 新特性

1. 增加`gf`命令行开发辅助工具：
    - 支持`GF`框架下载更新；
    - 支持初始化新建项目命令；
    - 支持跨平台交叉编译命令；
    - 命令行工具支持自动更新命令；
    - 支持二进制文件打包，生成二进制文件或者Go程序文件；
    - 支持指定数据库生成数据表模型，支持本地配置文件读取数据库配置；
1. 新增`gi18n`国际化管理模块；
1. 新增`gres`资源管理器模块：
    - 资源管理器支持虚拟的文件/目录操作方法；
    - 默认整合支持到了WebServer、配置管理、模板引擎中；
    - 可将任意的文件打包为`Go`内容，支持开发者自定义加解密；
    - 任意文件如网站静态文件、配置文件等可编译到二进制文件中，也可编译到发布的可执行文件中；
    - 开发者可只需编译发布一个可执行文件，除了方便了软件分发，也为保护软件知识产权内容提供了可能；
1. 新增`gini`模块：
    - 支持`ini`文件的读取/生成；
    - 同时配置管理模块也增加了对`ini`文件的支持；
    - 配置管理模块目前自动识别支持`ini/xml/json/toml/yaml`五种数据格式；
1. `Session`功能重构，新增`gsession`模块，`WebServer`默认使用文件存储`Session`；
1. `WebServer`新增中间件特性，并保留原有的HOOK设计，两者都可实现请求拦截、预处理等等特性；
1. `WebServer`新增更便捷的层级路由注册方式；
1. `gcmd`命令行参数解析模块重构，增加`Parser`解析对象；
1. 新增`gdebug`模块，用于堆栈获取/打印；


## 重大调整
1. 去掉`1.x`版本中已经被标记为`deprecated`的方法；
1. 调整`container`分类的容器模块，将默认并发安全参数调整为默认非并发安全；
1. 目录调整:
    - 去掉`third`目录，统一使用`go module`管理包依赖；
    - 将原有`g`目录中的模块移出到框架主目录，原有的`g`模块移动到`frame/g`目录；
    - 将原有`geg`示例代码目录名称修改为`.example`；



## 功能改进

1. `ghttp`
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
    - `debug`模式支持完整带参数整合的SQL语句调试打印；
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