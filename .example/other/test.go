package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func AddAPKCmdTask11(assistantId int, cmd int32, cmdData []byte, FromClientId string, desc string, priority int, status int32, cmdkey string) (int64, error) {
	//var res, err = g.DB("test").Insert("assistant_tasks", g.Map{
	//	"assistant_id": assistantId,
	//	"cmd":          cmd,
	//	"cmdData":      cmdData,
	//	"status":       status,
	//	"FromClientId": FromClientId,
	//	"desc":         desc,
	//	"priority":     priority,
	//	"cmdkey":       cmdkey,
	//})
	var res, err = g.DB("test").Table("assistant_tasks").Data(g.Map{
		"assistant_id": assistantId,
		"cmd":          cmd,
		"cmdData":      cmdData,
		"status":       status,
		"FromClientId": FromClientId,
		"desc":         desc,
		"priority":     priority,
		"cmdkey":       cmdkey,
	}).Insert()

	if err != nil {
		glog.Error("插入手机任务队列报错", err.Error())
		return 0, err
	}
	taskId, err := res.LastInsertId()
	return taskId, err
}

func main() {
	g.DB().SetDebug(true)
	AddAPKCmdTask11(1, 2058, []byte(""), "", "", 60, 0, "")
}
