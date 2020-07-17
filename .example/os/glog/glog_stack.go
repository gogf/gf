package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/glog"
)

func main() {

	glog.PrintStack()
	glog.New().PrintStack()

	fmt.Println(glog.GetStack())
	fmt.Println(glog.New().GetStack())
}
