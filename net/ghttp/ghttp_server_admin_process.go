// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// Allow executing management command after server starts after this interval in milliseconds.
	adminActionIntervalLimit = 2000
	adminActionNone          = 0
	adminActionRestarting    = 1
	adminActionShuttingDown  = 2
	adminActionReloadEnvKey  = "GF_SERVER_RELOAD"
	adminActionRestartEnvKey = "GF_SERVER_RESTART"
	adminGProcCommGroup      = "GF_GPROC_HTTP_SERVER"
)

var (
	// serverActionLocker is the locker for server administration operations.
	serverActionLocker sync.Mutex

	// serverActionLastTime is timestamp in milliseconds of last administration operation.
	serverActionLastTime = gtype.NewInt64(gtime.TimestampMilli())

	// serverProcessStatus is the server status for operation of current process.
	serverProcessStatus = gtype.NewInt()
)

// RestartAllServer restarts all the servers of the process gracefully.
// The optional parameter `newExeFilePath` specifies the new binary file for creating process.
func RestartAllServer(ctx context.Context, newExeFilePath string) error {
	if !gracefulEnabled {
		return gerror.NewCode(gcode.CodeInvalidOperation, "graceful reload feature is disabled")
	}
	serverActionLocker.Lock()
	defer serverActionLocker.Unlock()
	if err := checkProcessStatus(); err != nil {
		return err
	}
	if err := checkActionFrequency(); err != nil {
		return err
	}
	return restartWebServers(ctx, nil, newExeFilePath)
}

// ShutdownAllServer shuts down all servers of current process gracefully.
func ShutdownAllServer(ctx context.Context) error {
	serverActionLocker.Lock()
	defer serverActionLocker.Unlock()
	if err := checkProcessStatus(); err != nil {
		return err
	}
	if err := checkActionFrequency(); err != nil {
		return err
	}
	shutdownWebServersGracefully(ctx, nil)
	return nil
}

// checkProcessStatus checks the server status of current process.
func checkProcessStatus() error {
	status := serverProcessStatus.Val()
	if status > 0 {
		switch status {
		case adminActionRestarting:
			return gerror.NewCode(gcode.CodeInvalidOperation, "server is restarting")

		case adminActionShuttingDown:
			return gerror.NewCode(gcode.CodeInvalidOperation, "server is shutting down")
		}
	}
	return nil
}

// checkActionFrequency checks the operation frequency.
// It returns error if it is too frequency.
func checkActionFrequency() error {
	interval := gtime.TimestampMilli() - serverActionLastTime.Val()
	if interval < adminActionIntervalLimit {
		return gerror.NewCodef(
			gcode.CodeInvalidOperation,
			"too frequent action, please retry in %d ms",
			adminActionIntervalLimit-interval,
		)
	}
	serverActionLastTime.Set(gtime.TimestampMilli())
	return nil
}

// forkReloadProcess creates a new child process and copies the fd to child process.
func forkReloadProcess(ctx context.Context, newExeFilePath ...string) error {
	var (
		path = os.Args[0]
	)
	if len(newExeFilePath) > 0 && newExeFilePath[0] != "" {
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
				for _, item := range gstr.SplitAndTrim(fdv, ",") {
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
	p.Env = append(p.Env, adminActionReloadEnvKey+"="+string(buffer))
	if _, err := p.Start(ctx); err != nil {
		glog.Errorf(
			ctx,
			"%d: fork process failed, error:%s, %s",
			gproc.Pid(), err.Error(), string(buffer),
		)
		return err
	}
	return nil
}

// forkRestartProcess creates a new server process.
func forkRestartProcess(ctx context.Context, newExeFilePath ...string) error {
	var (
		path = os.Args[0]
	)
	if len(newExeFilePath) > 0 && newExeFilePath[0] != "" {
		path = newExeFilePath[0]
	}
	if err := os.Unsetenv(adminActionReloadEnvKey); err != nil {
		intlog.Errorf(ctx, `%+v`, err)
	}
	env := os.Environ()
	env = append(env, adminActionRestartEnvKey+"=1")
	p := gproc.NewProcess(path, os.Args, env)
	if _, err := p.Start(ctx); err != nil {
		glog.Errorf(
			ctx,
			`%d: fork process failed, error:%s, are you running using "go run"?`,
			gproc.Pid(), err.Error(),
		)
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
		for k := range j.Var().Map() {
			m := make(map[string]string)
			for mapKey, mapValue := range j.Get(k).MapStrStr() {
				m[mapKey] = mapValue
			}
			sfm[k] = m
		}
	}
	return sfm
}

// restartWebServers restarts all servers.
func restartWebServers(ctx context.Context, signal os.Signal, newExeFilePath string) error {
	serverProcessStatus.Set(adminActionRestarting)
	if runtime.GOOS == "windows" {
		if signal != nil {
			// Controlled by signal.
			forceCloseWebServers(ctx)
			if err := forkRestartProcess(ctx, newExeFilePath); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
			return nil
		}
		// Controlled by web page.
		// It should ensure the response wrote to client and then close all servers gracefully.
		gtimer.SetTimeout(ctx, time.Second, func(ctx context.Context) {
			forceCloseWebServers(ctx)
			if err := forkRestartProcess(ctx, newExeFilePath); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		})
		return nil
	}
	if err := forkReloadProcess(ctx, newExeFilePath); err != nil {
		glog.Printf(ctx, "%d: server restarts failed", gproc.Pid())
		serverProcessStatus.Set(adminActionNone)
		return err
	} else {
		if signal != nil {
			glog.Printf(ctx, "%d: server restarting by signal: %s", gproc.Pid(), signal)
		} else {
			glog.Printf(ctx, "%d: server restarting by web admin", gproc.Pid())
		}
	}

	return nil
}

// shutdownWebServersGracefully gracefully shuts down all servers.
func shutdownWebServersGracefully(ctx context.Context, signal os.Signal) {
	serverProcessStatus.Set(adminActionShuttingDown)
	if signal != nil {
		glog.Printf(
			ctx,
			"%d: server gracefully shutting down by signal: %s",
			gproc.Pid(), signal.String(),
		)
	} else {
		glog.Printf(ctx, "%d: server gracefully shutting down by api", gproc.Pid())
	}
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			server := v.(*Server)
			server.doServiceDeregister()
			for _, s := range server.servers {
				s.shutdown(ctx)
			}
		}
	})
}

// forceCloseWebServers forced shuts down all servers.
func forceCloseWebServers(ctx context.Context) {
	serverMapping.RLockFunc(func(m map[string]interface{}) {
		for _, v := range m {
			for _, s := range v.(*Server).servers {
				s.close(ctx)
			}
		}
	})
}

// handleProcessMessage receives and handles the message from processes,
// which are commonly used for graceful reloading feature.
func handleProcessMessage() {
	var (
		ctx = context.TODO()
	)
	for {
		if msg := gproc.Receive(adminGProcCommGroup); msg != nil {
			if bytes.EqualFold(msg.Data, []byte("exit")) {
				intlog.Printf(ctx, "%d: process message: exit", gproc.Pid())
				shutdownWebServersGracefully(ctx, nil)
				allShutdownChan <- struct{}{}
				intlog.Printf(ctx, "%d: process message: exit done", gproc.Pid())
				return
			}
		}
	}
}
