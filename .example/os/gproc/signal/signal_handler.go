package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func signalHandlerForMQ() {
	var (
		sig          os.Signal
		receivedChan = make(chan os.Signal)
	)
	signal.Notify(
		receivedChan,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
	for {
		sig = <-receivedChan
		fmt.Println("MQ is shutting down due to signal:", sig.String())
		time.Sleep(time.Second)
		fmt.Println("MQ is shut down smoothly")
		return
	}
}

func main() {
	fmt.Println("Process start, pid:", os.Getpid())
	go signalHandlerForMQ()

	var (
		sig          os.Signal
		receivedChan = make(chan os.Signal)
	)
	signal.Notify(
		receivedChan,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
	)
	for {
		sig = <-receivedChan
		fmt.Println("MainProcess is shutting down due to signal:", sig.String())
		return
	}
}
