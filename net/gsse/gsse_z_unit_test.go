package gsse_test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsse"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"testing"
	"time"
)

func Test_SendMessage(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/sse", gsse.Handle(func(client *gsse.Client) {
		client.SendMessage("send message")
	}))
	s.SetDumpRouterMap(false)
	_ = s.Start()
	defer func() { _ = s.Shutdown() }()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(gctx.New(), "/sse"),
			"event:message\ndata:send message\n\n")
	})
}

func Test_SendMessageWithId(t *testing.T) {
	ch := make(chan bool, 1)
	s := g.Server(guid.S())
	s.BindHandler("/sse", gsse.Handle(func(client *gsse.Client) {
		client.OnClose(func(client *gsse.Client) {
			ch <- client.Terminated()
		})
		client.SendMessageWithId("send message with id", "1")
	}))
	s.SetDumpRouterMap(false)
	_ = s.Start()
	defer func() { _ = s.Shutdown() }()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(gctx.New(), "/sse"),
			"event:message\ndata:send message with id\nid:1\n\n")

		select {
		case value := <-ch:
			t.AssertEQ(value, true)
		}
	})
}

func Test_SendEvent(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/sse", gsse.Handle(func(client *gsse.Client) {
		client.KeepAlive()
		go func(client *gsse.Client) {
			<-time.After(time.Second)
			client.SendEvent("test", "send event")
			client.Close()
		}(client)
	}))
	s.SetDumpRouterMap(false)
	_ = s.Start()
	defer func() { _ = s.Shutdown() }()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(gctx.New(), "/sse"),
			":\n\nevent:test\ndata:send event\n\n")
	})
}

func Test_SendEventWithId(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/sse", gsse.Handle(func(client *gsse.Client) {
		client.KeepAlive()
		go func(client *gsse.Client) {
			<-time.After(time.Second)
			client.SendEventWithId("test", "send event", "2")
			client.Close()
		}(client)
	}))
	s.SetDumpRouterMap(false)
	_ = s.Start()
	defer func() { _ = s.Shutdown() }()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)

		t.Assert(client.GetContent(gctx.New(), "/sse"),
			":\n\nevent:test\ndata:send event\nid:2\n\n")
	})
}
