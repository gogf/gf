module github.com/gogf/gf/example

go 1.15

require (
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.0.0-rc2
	github.com/gogf/gf/v2 v2.0.0-rc2
)

replace (
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/v2 => ../
)
