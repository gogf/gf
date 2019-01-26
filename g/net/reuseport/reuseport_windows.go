// +build windows

package reuseport

// See net.RawConn.Control
//func Control(network, address string, c syscall.RawConn) (err error) {
//	return c.Control(func(fd uintptr) {
//        if err = windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1); err != nil {
//            return
//        }
//	})
//}
