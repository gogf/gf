// +build linux darwin dragonfly freebsd netbsd openbsd

package reuseport

import (
    "fmt"
    "gitee.com/johng/gf/third/golang.org/x/sys/unix"
    "syscall"
)

// See net.RawConn.Control
func Control(network, address string, c syscall.RawConn) (err error) {
	c.Control(func(fd uintptr) {
	    fmt.Println("addr", fd, int(fd))
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1); err != nil {
            panic(err)
		    return
		}
        fmt.Println("port", fd, int(fd))
		if err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
			panic(err)
		    return
		}
	})
	return
}
