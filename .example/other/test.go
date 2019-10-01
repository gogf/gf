package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	b, _ := json.Marshal([]interface{}{1, 2, 3, 4, 5, 123.456, "a"})
	fmt.Println(gconv.String(b))
}
