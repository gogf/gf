package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	// 去年今日
	fmt.Println(gtime.Now().AddDate(-1, 0, 0).Format("Y-m-d"))

	// 去年今日，UTC时间
	fmt.Println(gtime.Now().AddDate(-1, 0, 0).Format("Y-m-d H:i:s T"))
	fmt.Println(gtime.Now().AddDate(-1, 0, 0).UTC().Format("Y-m-d H:i:s T"))

	// 下个月1号凌晨0点整
	fmt.Println(gtime.Now().AddDate(0, 1, 0).Format("Y-m-d 00:00:00"))

	// 2个小时前
	fmt.Println(gtime.Now().Add(-time.Hour).Format("Y-m-d H:i:s"))
}
