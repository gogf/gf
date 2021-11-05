package gfile_test

import (
	"fmt"

	"github.com/gogf/gf/v2/os/gfile"
)

func ExampleHome() {
	// user's home directory
	homePath, _ := gfile.Home()
	fmt.Println(homePath)

	// May Output:
	// C:\Users\hailaz
}
