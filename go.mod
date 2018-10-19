module gitee.com/johng/gf

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/sarama v1.19.0
	github.com/Shopify/toxiproxy v2.1.3+incompatible // indirect
	github.com/axgle/mahonia v0.0.0-20180208002826-3358181d7394
	github.com/clbanning/mxj v1.8.2
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.1.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fatih/structs v1.1.0
	github.com/fsnotify/fsnotify v1.4.7
	github.com/ghodss/yaml v1.0.0
	github.com/go-sql-driver/mysql v1.4.0
	github.com/gofrs/flock v0.7.0 // indirect
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/gorilla/websocket v1.4.0
	github.com/grokify/html-strip-tags-go v0.0.0-20180907063347-e9e44961e26f
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-runewidth v0.0.3 // indirect
	github.com/olekukonko/tablewriter v0.0.0-20180912035003-be2c049b30cc
	github.com/pierrec/lz4 v2.0.5+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20181016184325-3113b8401b8a // indirect
	github.com/theckman/go-flock v0.7.0
	golang.org/x/sys v0.0.0
	gopkg.in/check.v1 v1.0.0
	gopkg.in/yaml.v2 v2.2.1
)

replace (
	golang.org/x/sys => ./vendor/golang.org/x/sys
	gopkg.in/check.v1 => ./vendor/gopkg.in/check.v1
	gopkg.in/yaml.v2 => ./vendor/gopkg.in/yaml.v2
)
