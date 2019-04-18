package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gcron"
	"time"
)

func test() {

}

func main() {
	_, err := gcron.Add("*/10 * * * * ?", test)
	if err != nil {
		panic(err)
	}
	g.Dump(gcron.Entries())
	time.Sleep(10 * time.Second)
}
