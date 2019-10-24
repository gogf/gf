package main

import (
	"context"
	"fmt"
)

func main() {
	c := context.WithValue(context.Background(), "key", "value")
	fmt.Printf("%v", c)
}
