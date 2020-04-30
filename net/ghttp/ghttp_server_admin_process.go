// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
)

const (
	// Allow executing management command after server starts after this interval in milliseconds.
	gADMIN_ACTION_INTERVAL_LIMIT = 2000
	gADMIN_ACTION_NONE           = 0
	gADMIN_ACTION_RESTARTING     = 1
	gADMIN_ACTION_SHUTINGDOWN    = 2
	gADMIN_ACTION_RELOAD_ENVKEY  = "GF_SERVER_RELOAD"
	gADMIN_ACTION_RESTART_ENVKEY = "GF_SERVER_RESTART"
	gADMIN_GPROC_COMM_GROUP      = "GF_GPROC_HTTP_SERVER"
)

// serverActionLocker is the locker for server administration operations.
var serverActionLocker sync.Mutex

// serverActionLastTime is timestamp in milliseconds of last administration operation.
var serverActionLastTime = gtype.NewInt64(gtime.TimestampMilli())

// serverProcessStatus is the server status for operation of current process.
var serverProcessStatus = gtype.NewInt()

// RestartAllServer restarts all the servers of the process.
// The optional parameter <newExeFilePath> specifies the new binary file for creating process.
func RestartAllServer(newExeFilePath ...string) error {
	serverActionLocker.Lock()
	defer serverActionLocker.Unlock()
	if err := checkProcessStatus(); err != nil {
		return err
	}
	if err := checkActionFrequency(); err != nil {
		return err
	}
	return restartWebServers("", newExeFilePath...)
}

// ShutdownAllServer shuts down all servers of current process.
func ShutdownAllServer() error {
	serverActionLocker.Lock()
	defer serverActionLocker.Unlock()
	if err := checkProcessStatus(); err != nil {
		return err
	}
	if err := checkActionFrequency(); err != nil {
		return err
	}
	shutdownWebServers()
	return nil
}

// checkProcessStatus checks the server status of current process.
func checkProcessStatus() error {
	status := serverProcessStatus.Val()
	if status > 0 {
		switch status {
		case gADMIN_ACTION_RESTARTING:
			return errors.New("server is restarting")
		case gADMIN_ACTION_SHUTINGDOWN:
			return errors.New("server is shutting down")
		}
	}
	return nil
}

// checkActionFrequency checks the operation frequency.
// It returns error if it is too frequency.
func checkActionFrequency() error {
	interval := gtime.TimestampMilli() - serverActionLastTime.Val()
	if interval < gADMIN_ACTION_INTERVAL_LIMIT {
		return errors.New(fmt.Sprintf("too frequent action, please retry in %d ms", gADMIN_ACTION_INTERVAL_LIMIT-interval))
	}
	serverActionLastTime.Set(gtime.TimestampMilli())
	return nil
}

// forkReloadProcess creates a new child process and copies the fd to child process.
func forkReloadProcess(newExeFilePath ...string) error {
	path := os.Args[0]
	if len(newExeFilePath) > 0 {
		path = newExeFilePath[0]
	}
	var (
		p   = gproc.NewProcess(path, os.Args, os.Environ())
		sfm = getServerFdMap()
	)
	for name, m := range sfm {
		for fdk, fdv := range m {
			if len(fdv) > 0 {
				s := ""
				for _, item := range strings.Split(fdv, ",") {
					array := strings.Split(item, "#")
					fd := uintptr(gconv.Uint(array[1]))
					if fd > 0 {
						s += fmt.Sprintf("%s#%d,", array[0], 3+len(p.ExtraFiles))
						p.ExtraFiles = append(p.ExtraFiles, os.NewFile(fd, ""))
					} else {
						s += fmt.Sprintf("%s#%d,", array[0], 0)
					}
				}
				sfm[name][fdk] = strings.TrimRight(s, ",")
			}
		}
	}
	buffer, _ := gjson.Encode(sfm)
	p.Env = append(p.Env, gADMIN_ACTION_RELOAD_ENVKEY+"="+string(buffer))
	if _, err := p.Start(); err != nil {
		glog.Errorf("%d: fork process failed, error:%s, %s", gproc.Pid(), err.Error(), string(buffer))
		return err
	}
	return nil
}

// forkRestartProcess creates a new server process.
func forkRestartProcess(newExeFilePath ...string) error {
	path := os.Args[0]
	if len(newExeFilePath) > 0 {
		path = newExeFilePath[0]
	}
	os.Unsetenv(gADMIN_ACTION_RELOAD_ENVKEY)
	env := os.Environ()
	env = append(env, gADMIN_ACTION_RESTART_ENVKEY+"=1")
	p := gproc.NewProcess(path, os.Args, env)
	if _, err := p.Start(); err != nil {
		glog.Errorf("%d: fork process failed, error:%s", gproc.Pid(), err.Error())
		return err
	}
	return nil
}

// getServerFdMap returns all the servers name to file descriptor mapping as map.
func getServerFdMap() map[string]listenerFdMap {
	sfm := make(map[string]listenerFdMap)
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for k, v := range m {
			sfm[k] = v.(*Server).getListenerFdMap()
		}
	})
	return sfm
}

// bufferToServerFdMap converts binary content to fd map.
func bufferToServerFdMap(buffer []byte) map[string]listenerFdMap {
	sfm := make(map[string]listenerFdMap)
	if len(buffer) > 0 {
		j, _ := gjson.LoadContent(buffer)
		for k, _ := range j.ToMap() {
			m := make(map[string]string)
			for k, v := range j.GetMap(k) {
				m[k] = gconv.String(v)
			}
			sfm[k] = m
		}
	}
	return sfm
}

// restartWebServers restarts all servers.
func restartWebServers(signal string, newExeFilePath ...string) error {
	serverProcessStatus.Set(gADMIN_ACTION_RESTARTING)
	if runtime.GOOS == "windows" {
		if len(signal) > 0 {
			// Controlled by signal.
			forceCloseWebServers()
			forkRestartProcess(newExeFilePath...)
		} else {
			// Controlled by web page.
			// It should ensure the response wrote to client and then close all servers gracefully.
			gtimer.SetTimeout(time.Second, func() {
				forceCloseWebServers()
				forkRestartProcess(newExeFilePath...)
			})
		}
	} else {
		if err := forkReloadProcess(newExeFilePath...); err != nil {
			glog.Printf("%d: server restarts failed", gproc.Pid())
			serverProcessStatus.Set(gADMIN_ACTION_NONE)
			return err
		} else {
			if len(signal) > 0 {
				glog.Printf("%d: server restarting by signal: %s", gproc.Pid(), signal)
			} else {
				glog.Printf("%d: server restarting by web admin", gproc.Pid())
			}

		}
	}
	return nil
}

// shutdownWebServers shuts down all servers.
func shutdownWebServers(signal ...string) {
	serverProcessStatus.Set(gADMIN_ACTION_SHUTINGDOWN)
	if len(signal) > 0 {
		glog.Printf("%d: server shutting down by signal: %s", gproc.Pid(), signal[0])
		forceCloseWebServers()
		allDoneChan <- struct{}{}
	} else {
		glog.Printf("%d: server shutting down by api", gproc.Pid())
		gtimer.SetTimeout(time.Second, func() {
			forceCloseWebServers()
			allDoneChan <- struct{}{}
		})
	}
}

// gracefulShutdownWebServers gracefully shuts down all servers.
func gracefulShutdownWebServers() {
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			for _, s := range v.(*Server).servers {
				s.shutdown()
			}
		}
	})
}

// forceCloseWebServers forced shuts down all servers.
func forceCloseWebServers() {
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			for _, s := range v.(*Server).servers {
				s.close()
			}
		}
	})
}

// handleProcessMessage listens the signal from system.
func handleProcessMessage() {
	for {
		if msg := gproc.Receive(gADMIN_GPROC_COMM_GROUP); msg != nil {
			if bytes.EqualFold(msg.Data, []byte("exit")) {
				gracefulShutdownWebServers()
				allDoneChan <- struct{}{}
				return
			}
		}
	}
}
