package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gcron"
	"time"
)

func test() {

}

func main() {
	_, err := gcron.Add("*/10 * * * * ?", test)
	fmt.Println(err)
	fmt.Println(gcron.Entries())
	time.Sleep(10 * time.Second)
}
