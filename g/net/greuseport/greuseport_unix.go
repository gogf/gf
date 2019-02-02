// +build linux darwin dragonfly freebsd netbsd openbsd

package greuseport

import (
    "github.com/gogf/gf/third/golang.org/x/sys/unix"
    "syscall"
)

func init() {
    Enabled = true
}

// See net.RawConn.Control
func Control(network, address string, c syscall.RawConn) (err error) {
	c.Control(func(fd uintptr) {
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
            panic(err)
		    return
		}
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
			panic(err)
		    return
		}

	})
	return
}
