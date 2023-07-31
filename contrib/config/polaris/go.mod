module github.com/gogf/gf/contrib/config/polaris/v2

go 1.15

require (
	github.com/gogf/gf/v2 v2.5.1
	github.com/polarismesh/polaris-go v1.5.1
)

replace github.com/gogf/gf/v2 => ../../../

replace (
	golang.org/x/net v0.2.0 => golang.org/x/net v0.0.0-20221019024206-cb67ada4b0ad
	golang.org/x/sys v0.2.0 => golang.org/x/sys v0.0.0-20220906165534-d0df966e6959
)
