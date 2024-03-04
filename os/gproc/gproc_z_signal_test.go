// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"github.com/gogf/gf/v2/test/gtest"
	"os"
	"syscall"
	"testing"
	"time"
)

func Test_Signal(t *testing.T) {
	go Listen()

	// non shutdown signal
	gtest.C(t, func(t *gtest.T) {
		sigRec := make(chan os.Signal, 1)
		AddSigHandler(func(sig os.Signal) {
			sigRec <- sig
		}, syscall.SIGUSR1, syscall.SIGUSR2)

		sendSignal(syscall.SIGUSR1)
		select {
		case s := <-sigRec:
			t.AssertEQ(s, syscall.SIGUSR1)
			t.AssertEQ(false, isWaitChClosed())
		case <-time.After(time.Second):
			t.Error("signal SIGUSR1 handler timeout")
		}

		sendSignal(syscall.SIGUSR2)
		select {
		case s := <-sigRec:
			t.AssertEQ(s, syscall.SIGUSR2)
			t.AssertEQ(false, isWaitChClosed())
		case <-time.After(time.Second):
			t.Error("signal SIGUSR2 handler timeout")
		}

		sendSignal(syscall.SIGHUP)
		select {
		case <-sigRec:
			t.Error("signal SIGHUP should not be listen")
		case <-time.After(time.Millisecond * 100):
		}

		// multiple listen
		go Listen()
		go Listen()
		sendSignal(syscall.SIGUSR1)
		cnt := 0
		timeout := time.After(time.Second)
		for {
			select {
			case <-sigRec:
				cnt++
			case <-timeout:
				if cnt == 0 {
					t.Error("signal SIGUSR2 handler timeout")
				}
				if cnt != 1 {
					t.Error("multi Listen() repetitive execution")
				}
				return
			}
		}
	})

	// test shutdown signal
	gtest.C(t, func(t *gtest.T) {
		sigRec := make(chan os.Signal, 1)
		AddSigHandlerShutdown(func(sig os.Signal) {
			sigRec <- sig
		})

		sendSignal(syscall.SIGTERM)
		select {
		case s := <-sigRec:
			t.AssertEQ(s, syscall.SIGTERM)
			t.AssertEQ(true, isWaitChClosed())
		case <-time.After(time.Second):
			t.Error("signal SIGUSR2 handler timeout")
		}
	})
}

func sendSignal(sig os.Signal) {
	signalChan <- sig
}

func isWaitChClosed() bool {
	select {
	case _, ok := <-waitChan:
		return !ok
	default:
		return false
	}
}
