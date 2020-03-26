package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
)

func main() {
	fmt.Println(gfile.Basename("/tmp/1585227151172826000/access.20200326205231173924.log"))
	fmt.Println(
		gregex.IsMatchString(`.+\.\d{20}\.log`,
			gfile.Basename("/tmp/1585227151172826000/access.20200326205231173924.log")))
}
