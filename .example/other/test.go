package main

import (
	"os"
	"time"
)

func main() {
	file, err := os.Create("/tmp/testfile")
	if err != nil {
		panic(err)
	}
	for {
		_, err = file.Write([]byte("test\n"))
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}
