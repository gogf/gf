package main

import (
	"github.com/gogf/gf/internal/intlog"
)

func main() {
	intlog.Print(1, 2, 3)
	intlog.Printf("%d", 1)
	intlog.Error(1)
	intlog.Errorf("%d", 1)
}
