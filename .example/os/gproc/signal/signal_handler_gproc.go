package main

import (
	"fmt"
	"github.com/gogf/gf/os/gproc"
	"os"
	"time"
)

func signalHandlerForMQ(sig os.Signal) {
	fmt.Println("MQ is shutting down due to signal:", sig.String())
	time.Sleep(time.Second)
	fmt.Println("MQ is shut down smoothly")
}

func signalHandlerForMain(sig os.Signal) {
	fmt.Println("MainProcess is shutting down due to signal:", sig.String())
}

func main() {
	fmt.Println("Process start, pid:", os.Getpid())
	gproc.AddSigHandlerShutdown(
		signalHandlerForMQ,
		signalHandlerForMain,
	)
	gproc.Listen()
}
