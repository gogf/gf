package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println(json.Marshal(nil))
}
