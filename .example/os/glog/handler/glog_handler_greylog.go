package main

//import (
//	"context"
//	"github.com/gogf/gf/v2/frame/g"
//	"github.com/gogf/gf/v2/os/glog"
//	"github.com/robertkowalski/graylog-golang"
//)
//
//var greyLogClient = gelf.New(gelf.Config{
//	GraylogPort:     80,
//	GraylogHostname: "graylog-host.com",
//	Connection:      "wan",
//	MaxChunkSizeWan: 42,
//	MaxChunkSizeLan: 1337,
//})
//
//// LoggingGreyLogHandler is an example handler for logging content to remote GreyLog service.
//var LoggingGreyLogHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
//	in.Next()
//	greyLogClient.Log(in.Buffer.String())
//}
//
//func main() {
//	g.Log().SetHandlers(LoggingGreyLogHandler)
//
//	g.Log().Debug("Debugging...")
//	g.Log().Warning("It is warning info")
//	g.Log().Error("Error occurs, please have a check")
//	glog.Print("test log")
//}
