package gsse

import (
	"context"
	"github.com/gogf/gf/v2/net/ghttp"
	"time"
)

// Handle wraps func(*gsse.Client) to func(*ghttp.Request), use by ghttp.Server.BindHandler().
func Handle(fn func(*Client)) func(*ghttp.Request) {
	return func(request *ghttp.Request) {
		client := newClient(request)
		if fn != nil {
			fn(client)
		}

		var (
			keepAliveCtx    context.Context
			keepAliveCancel context.CancelFunc = func() {
				// empty func if not keep alive
			}
		)
		if client.keepAlive {
			keepAliveCtx, keepAliveCancel =
				context.WithCancel(context.Background())
		}
		go func() {
			<-client.Context().Done()
			if client.onClose != nil {
				go client.onClose(client)
			}
			keepAliveCancel()
		}()
		if client.keepAlive {
			for {
				select {
				case <-keepAliveCtx.Done():
					return
				default:
					client.heartbeat()
					time.Sleep(5 * time.Second)
				}
			}
		}
	}
}
