package main

import (
	"encoding/json"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gcoime"
)

func main() {
	g.Sliceg.Sli
	for _, line := range lines {
		gl := gdb.Map{}
		r := NewMysqlTableLogger()
		if err := json.Unmarshal([]byte(line), &r); err != nil || r == nil { //写入日志表数据错误
			glog.Errorf("%s[%s] SaveToDB json.DecodeTo : %v, log=%s", c.LogHead, c.LogStatus(), err, GetSubString(line, 1024))
			continue
		}
		if tableName == "" {
			tableName = r.GetTableName()
		}
		gls = append(gls, r)
	}
}
