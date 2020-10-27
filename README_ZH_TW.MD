# GoFrame
[![Go Doc](https://godoc.org/github.com/gogf/gf?status.svg)](https://godoc.org/github.com/gogf/gf) 
[![Build Status](https://travis-ci.org/gogf/gf.svg?branch=master)](https://travis-ci.org/gogf/gf) 
[![Go Report](https://goreportcard.com/badge/github.com/gogf/gf?v=1)](https://goreportcard.com/report/github.com/gogf/gf)
[![Code Coverage](https://codecov.io/gh/gogf/gf/branch/master/graph/badge.svg)](https://codecov.io/gh/gogf/gf/branch/master)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/gogf/gf)
[![License](https://img.shields.io/github/license/gogf/gf.svg?style=flat)](https://github.com/gogf/gf)

[English](README.MD) | [简体中文](README_ZH.MD) | 繁體中文

`GF(Go Frame)`是一款模組化、高性能、生產級的Go基礎開發框架。
實現了比較完善的基礎設施建設以及開發工具鏈，提供了常用的基礎開發模組，
如：快取、日誌、隊列、陣列、集合、容器、定時器、命令行、記憶體鎖、物件池、
配置管理、資源管理、資料驗證、資料編碼、定時任務、資料庫ORM、TCP/UDP組件、行程管理/通訊等等。
並提供了Web服務開發的系列核心元件，如：Router、Cookie、Session、Middleware、服務註冊、模組引擎等等，
支持熱重啟、熱更新、域名綁定、TLS/HTTPS、Rewrite等特性。

> 如果您初識`Go`語言，您可以將`GoFrame`類似於`PHP`中的`Laravel`, `Java`中的`SpringBoot`或者`Python`中的`Django`。

# 特點
* 模組化、鬆耦合設計；
* 模組豐富、開箱即用；
* 簡便易用、易於維護；
* 高程式品質、高單元測試覆蓋率；
* 社群活躍，大牛謙遜低調脾氣好；
* 詳盡的開發文件及範例；
* 完善的本地中文化支持；
* 設計為團隊及企業使用；

# 連結
- **主庫**：https://github.com/gogf/gf 
- **碼雲**：https://gitee.com/johng/gf 

# 安裝
```html
go get -u -v github.com/gogf/gf
```
推薦使用 `go.mod`:
```
require github.com/gogf/gf latest
```

# 限制
```shell
golang版本 >= 1.11
```



# 架構
<div align=center>
<img src="https://goframe.org/images/arch.png?v=12"/>
</div>

# 模組

1. **核心模組**

    `GoFrame`提供了一些基礎的、常用的模組，簡單、易用和輕量級，並保持極少的外部依賴，這些模組由`gf`主倉庫細緻維護。

1. **社群模組**

    社群模組主要由社群貢獻並維護，大部分也是由`gf`主倉庫的貢獻者提供及維護，存放於`gogf`空間下，與`gf`主倉庫處於同一級別。有的社群模組是從`gf`主倉庫中剝離出來單獨維護的模組，這些模組並不是特別常用，或者對外部依賴較重。


# 性能

以下是目前最流行的`WEB Server` Golang框架/類庫性能測試結果。
性能測試用例源程式碼倉庫: https://github.com/gogf/gf-performance

## 環境

    OS   : Ubuntu 18.04 amd64
    CPU  : AMD A8-6600K x 4
    MEM  : 32GB
    GO   : v1.13.4

## 工具

`ab`: Apache HTTP server benchmarking tool.

測試指令:
```
ab -t 10 -c 100 http://127.0.0.1:3000/hello
ab -t 10 -c 100 http://127.0.0.1:3000/query?id=10000
ab -t 10 -c 100 http://127.0.0.1:3000/json
```
並發客戶端數量從 `100` 遞增到 `10000`。

> 每個項目的每個用例均運行`5`次，取最優的結果展示。

## 1. Hello World
<table>
<tr>
<th>Throughputs</th>
<th>Mean Latency</th>
<th>P99 Latency</th>
</tr>
<tr>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/throughputs1.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/meanlatency1.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/p99latency1.jpeg"></td>
</tr>
</table>

## 2. Json Response
<table>
<tr>
<th>Throughputs</th>
<th>Mean Latency</th>
<th>P99 Latency</th>
</tr>
<tr>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/throughputs3.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/meanlatency3.jpeg"></td>
<td width="30%"><img src="http://gfcdn.johng.cn/images/performance/p99latency3.jpeg"></td>
</tr>
</table>


# 文件

開發文件：https://goframe.org

API文件：https://godoc.org/github.com/gogf/gf

# 幫助
- QQ交流群：[116707870](//shang.qq.com/wpa/qunwpa?idkey=195f91eceeb5d7fa76009b7cd5a4641f70bf4897b7f5a520635eb26ff17adfe7)
- 微信交流群：微信添加`389961817`備註`GF`
- 主庫ISSUE：https://github.com/gogf/gf/issues

> 建議通過閱讀`GoFrame`的源碼以及API文件深度學習`GoFrame`，了解更多的精妙設計。

# 授權

`GF` 使用非常友好的 [MIT](LICENSE) 開源授權進行發布，永久`100%`開源免費。

# 用戶

- [騰訊科技](https://www.tencent.com/)
- [中興科技](https://www.zte.com.cn/china/)
- [螞蟻金服](https://www.antfin.com/)
- [醫聯科技](https://www.medlinker.com/)
- [庫幣科技](https://www.kucoin.io/)
- [樂有家](https://www.leyoujia.com/)
- [IGG](https://igg.com)
- [喜馬拉雅](https://www.ximalaya.com)

> 在這裡只列舉了部分知名的用戶，如果您的企業或者產品正在使用`GoFrame`，歡迎到 [這裡](https://github.com/gogf/gf/issues/168) 留言。

# 貢獻

感謝所有參與`GoFrame`開發的貢獻者。 [[貢獻者列表](https://github.com/gogf/gf/graphs/contributors)].
<a href="https://github.com/gogf/gf/graphs/contributors"><img src="https://opencollective.com/goframe/contributors.svg?width=890&button=false" />< /a>


# 捐贈

如果您喜歡`GF`，要不[給開發者來杯咖啡吧](DONATOR.MD)！
請在捐贈時備註您的`github`/`gitee`賬號名稱。

# 贊助

贊助支持`GF`框架的快速研發，如果您感興趣，請聯繫 微信 `389961817` / 郵件 `john@goframe.org`。

# 感謝
<a href="https://www.jetbrains.com/?from=GoFrame"><img src="https://goframe.org/images/jetbrains.png" width="100" alt="JetBrains" /></a>

