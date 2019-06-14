package main

import (
<<<<<<< HEAD
	"fmt"
)

func main() {
	var i float64 = 0
	for index := 0; index < 10; index++ {
		i += 0.1
		fmt.Println(i)
	}
}
=======
	"github.com/gogf/gf/g/encoding/gjson"
)

func main() {
	j := gjson.New(`[1,2,3]`)
	j.Remove("1")
	j.Dump()
}
>>>>>>> master
