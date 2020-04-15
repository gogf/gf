package main

import "fmt"

// apiString is the type assert api for String.
type apiString interface {
	String() string
}

func main() {
	for i := 0; i < 10; i++ {
		switch 1 {
		case 1:
			continue
		}
		fmt.Println(i)
	}

}
