package main

import "fmt"

func main() {
	Path := "////////"
	if Path != "/" {
		for len(Path) > 1 && Path[len(Path)-1] == '/' {
			Path = Path[:len(Path)-1]
		}
	}
	fmt.Println(Path)
}
