package gproc

import (
	"github.com/gogf/gf/v2/test/gtest"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_Signal(t *testing.T) {
	var (
		sigRec     os.Signal
		sigsRec    = make([]os.Signal, 0)
		sigHandler = func(sig os.Signal) {
			sigRec = sig
			sigsRec = append(sigsRec, sig)
		}
	)

	clearTesting := func() {
		for sig, _ := range signalHandlerMap {
			signalHandlerMap[sig] = make([]SigHandler, 0)
		}
		sigRec = nil
		sigsRec = make([]os.Signal, 0)
	}

	go Listen()

	// non shutdown signal
	gtest.C(t, func(t *gtest.T) {
		defer clearTesting()

		AddSigHandler(sigHandler, syscall.SIGUSR1, syscall.SIGUSR2)

		sendSignalWithSleep(syscall.SIGUSR1)
		t.AssertEQ(sigRec, syscall.SIGUSR1)
		t.AssertEQ(false, isWaitChClosed())

		sendSignalWithSleep(syscall.SIGUSR2)
		t.AssertEQ(sigRec, syscall.SIGUSR2)
		t.AssertEQ(false, isWaitChClosed())

		sendSignalWithSleep(syscall.SIGHUP)
		t.AssertNE(sigRec, syscall.SIGHUP)
		t.AssertEQ(false, isWaitChClosed())
	})

	// test multiple listen case
	gtest.C(t, func(t *gtest.T) {
		go Listen()
		defer clearTesting()

		AddSigHandler(sigHandler, syscall.SIGUSR1)
		sendSignalWithSleep(syscall.SIGUSR1)
		t.AssertEQ(sigRec, syscall.SIGUSR1)
		t.AssertEQ(len(sigsRec), 1)
	})

	// test shutdown signal
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			clearTesting()
			waitChan = make(chan struct{}) // channel will be closed when shutdown signal received, reset wait chan
		}()

		AddSigHandlerShutdown(sigHandler)

		sendSignalWithSleep(syscall.SIGTERM)
		t.AssertEQ(sigRec, syscall.SIGTERM)
		t.AssertEQ(true, isWaitChClosed())
	})
}

func sendSignalWithSleep(sig os.Signal) {
	signalChan <- sig
	time.Sleep(time.Millisecond * 100) // make sure handlers executed
}

func isWaitChClosed() bool {
	select {
	case _, ok := <-waitChan:
		return !ok
	default:
		return false
	}
}
