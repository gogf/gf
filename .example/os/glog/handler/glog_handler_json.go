package main

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"os"
)

// JsonOutputsForLogger is for JSON marshaling in sequence.
type JsonOutputsForLogger struct {
	Time    string `json:"time"`
	Level   string `json:"level"`
	Content string `json:"content"`
}

// LoggingJsonHandler is an example handler for logging JSON format content.
var LoggingJsonHandler glog.Handler = func(ctx context.Context, in *glog.HandlerInput) {
	jsonForLogger := JsonOutputsForLogger{
		Time:    in.TimeFormat,
		Level:   in.LevelFormat,
		Content: gstr.Trim(in.String()),
	}
	jsonBytes, err := json.Marshal(jsonForLogger)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		return
	}
	in.Buffer.Write(jsonBytes)
	in.Buffer.WriteString("\n")
	in.Next()
}

func main() {
	g.Log().SetHandlers(LoggingJsonHandler)

	g.Log().Debug("Debugging...")
	g.Log().Warning("It is warning info")
	g.Log().Error("Error occurs, please have a check")
}
