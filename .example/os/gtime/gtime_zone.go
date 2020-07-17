package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	// 先使用标准库打印当前时间
	fmt.Println(time.Now().String())
	// 设置进程时区，全局有效
	err := gtime.SetTimeZone("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	// 使用gtime获取当前时间
	fmt.Println(gtime.Now().String())
	// 使用标准库获取当前时间
	fmt.Println(time.Now().String())
}
