package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/os/gfpool"
)

func main() {
	for {
		time.Sleep(time.Second)
		if f, err := gfpool.Open("/home/john/temp/log.log", os.O_RDONLY, 0666, time.Hour); err == nil {
			fmt.Println(f.Name())
			f.Close()
		} else {
			fmt.Println(err)
		}
	}
}
