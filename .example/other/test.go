package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.DB().SetDebug(true)
	tx, err := g.DB().Begin()
	if err != nil {
		panic(err)
	}
	smsTaskInfo := "`sms_sys`.sms_task_info"
	m := tx.Table(smsTaskInfo)
	_, err = m.Where("`delete`=0 AND is_review_temp=1 AND UNIX_TIMESTAMP(NOW())-UNIX_TIMESTAMP(create_time)>=86400").Delete()
}
