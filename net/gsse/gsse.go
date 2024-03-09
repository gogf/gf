package gsse

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
)

// Client wraps the SSE(Server-Sent Event) ghttp.Request and provides SSE APIs
type Client struct {
	request   *ghttp.Request
	cancel    context.CancelFunc
	onClose   func(*Client)
	keepAlive bool
}

const (
	noEvent = ""
	message = "message"

	noId = ""

	emptyComment = ""
)
