package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/container/gtype"
	"github.com/jin502437344/gf/os/gtimer"
)

func main() {
	v := gtype.NewInt()
	//w := gtimer.New(10, 10*time.Millisecond)
	fmt.Println("start:", time.Now())
	for i := 0; i < 1000000; i++ {
		gtimer.AddTimes(time.Second, 1, func() {
			v.Add(1)
		})
	}
	fmt.Println("end  :", time.Now())
	time.Sleep(1000 * time.Millisecond)
	fmt.Println(v.Val(), time.Now())

	//gtimer.AddSingleton(time.Second, func() {
	//    fmt.Println(time.Now().String())
	//})
	//select { }
}
