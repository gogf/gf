# 流水线说明

## gf.yml

### Setup Golang

安装 Golang 环境，并根据 `**/go.mod` 的 hash 加载缓存。

### 单元测试

在普通提交和合并请求的时候，执行带有缓存的单元测试。

在 master 的 push 的时候（合并到 master），执行带有覆盖率报告的单元测试。
