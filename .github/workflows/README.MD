# 流水线说明

## gf.yml

### Setup Golang

安装 Golang 环境。

### Cache Go Modules

还原缓存，命中同分支以及父分支和默认分支的同 key 缓存。

restore-keys 可以部分匹配缓存（这是 Setup Golang 的缓存无法做到的，所以另外设置缓存配置）

### 单元测试

在普通提交和合并请求的时候，执行带有缓存的单元测试。

在 master 的 push 的时候（合并到 master），执行带有覆盖率报告的单元测试。
