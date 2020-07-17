package main

import (
	"encoding/json"
	"fmt"

	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	t := gtime.Now()
	b, err := json.Marshal(t)
	fmt.Println(err)
	fmt.Println(string(b))
}
