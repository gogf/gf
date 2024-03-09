package gsse

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gmutex"
)

// Client wraps the SSE(Server-Sent Event) ghttp.Request and provides SSE APIs
type Client struct {
	request   *ghttp.Request
	cancel    context.CancelFunc
	onClose   func(*Client)
	keepAlive bool
	mutex     *gmutex.Mutex
}

const (
	noEvent      = ""
	noId         = ""
	emptyComment = ""
)
