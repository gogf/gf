package main

import (
	"context"
	"fmt"
)

func main() {
	c1 := context.WithValue(context.Background(), "key1", "value1")
	c2 := context.WithValue(c1, "key2", "value2")
	fmt.Printf("%v", c2.Value("key2"))
}
