package main

//import (
//	"github.com/gogf/gf/v2/os/glog"
//	"github.com/robertkowalski/graylog-golang"
//)
//
//type MyGrayLogWriter struct {
//	gelf    *gelf.Gelf
//	logger  *glog.Logger
//}
//
//func (w *MyGrayLogWriter) Write(p []byte) (n int, err error) {
//	w.gelf.Send(p)
//	return w.logger.Write(p)
//}
//
//func main() {
//	glog.SetWriter(&MyGrayLogWriter{
//		logger : glog.New(),
//		gelf   : gelf.New(gelf.Config{
//			GraylogPort     : 80,
//			GraylogHostname : "graylog-host.com",
//			Connection      : "wan",
//			MaxChunkSizeWan : 42,
//			MaxChunkSizeLan : 1337,
//		}),
//	})
//	glog.Print("test log")
//}
