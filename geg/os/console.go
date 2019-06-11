package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g/os/gcmd"
)

func doEcho() {
    fmt.Println("do echo")
}

func main() {
    fmt.Println(gcmd.Value.GetAll())

    fmt.Println(gcmd.Value.Get(1))

    gcmd.BindHandle("echo", doEcho)
    gcmd.RunHandle("echo")

    gcmd.AutoRun()
=======
	"fmt"
	"github.com/gogf/gf/g/os/gcmd"
)

func doEcho() {
	fmt.Println("do echo")
}

func main() {
	fmt.Println(gcmd.Value.GetAll())

	fmt.Println(gcmd.Value.Get(1))

	gcmd.BindHandle("echo", doEcho)
	gcmd.RunHandle("echo")

	gcmd.AutoRun()
>>>>>>> upstream/master
}
