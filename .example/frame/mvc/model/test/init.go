package test

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

func DB() gdb.DB {
	return g.DB()
}
