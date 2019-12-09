package main

import (
	"github.com/gogf/gf/os/gproc"
	"log"
)

func main() {
	process := gproc.NewProcessCmd("go run test.go")
	if pid, err := process.Start(); err != nil {
		log.Print("Build failed:", err)
		return
	} else {
		log.Printf("Build running: pid: %d", pid)
	}

}
