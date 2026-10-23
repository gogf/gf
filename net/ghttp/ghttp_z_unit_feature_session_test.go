// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Session_Cookie(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set(r.Get("k").String(), r.Get("v").String())
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Get(r.Get("k").String()))
	})
	s.BindHandler("/remove", func(r *ghttp.Request) {
		r.Session.Remove(r.Get("k").String())
	})
	s.BindHandler("/clear", func(r *ghttp.Request) {
		r.Session.RemoveAll()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetBrowserMode(true)
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		r1, e1 := client.Get(ctx, "/set?k=key1&v=100")
		if r1 != nil {
			defer r1.Close()
		}
		t.Assert(e1, nil)
		t.Assert(r1.ReadAllString(), "")

		t.Assert(client.GetContent(ctx, "/set?k=key2&v=200"), "")

		t.Assert(client.GetContent(ctx, "/get?k=key1"), "100")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/get?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key4"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/clear"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "")
	})
}

func Test_Session_Header(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set(r.Get("k").String(), r.Get("v").String())
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Get(r.Get("k").String()))
	})
	s.BindHandler("/remove", func(r *ghttp.Request) {
		r.Session.Remove(r.Get("k").String())
	})
	s.BindHandler("/clear", func(r *ghttp.Request) {
		r.Session.RemoveAll()
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		response, e1 := client.Get(ctx, "/set?k=key1&v=100")
		if response != nil {
			defer response.Close()
		}
		sessionId := response.GetCookie(s.GetSessionIdName())
		t.Assert(e1, nil)
		t.AssertNE(sessionId, nil)
		t.Assert(response.ReadAllString(), "")

		client.SetHeader(s.GetSessionIdName(), sessionId)

		t.Assert(client.GetContent(ctx, "/set?k=key2&v=200"), "")

		t.Assert(client.GetContent(ctx, "/get?k=key1"), "100")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/get?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key3"), "")
		t.Assert(client.GetContent(ctx, "/remove?k=key4"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key1"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "200")
		t.Assert(client.GetContent(ctx, "/clear"), "")
		t.Assert(client.GetContent(ctx, "/get?k=key2"), "")
	})
}

func Test_Session_StorageFile(t *testing.T) {
	sessionId := ""
	s := g.Server(guid.S())
	s.BindHandler("/set", func(r *ghttp.Request) {
		r.Session.Set(r.Get("k").String(), r.Get("v").String())
		r.Response.Write(r.Get("k").String(), "=", r.Get("v").String())
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		r.Response.Write(r.Session.Get(r.Get("k").String()))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		response, e1 := client.Get(ctx, "/set?k=key&v=100")
		if response != nil {
			defer response.Close()
		}
		sessionId = response.GetCookie(s.GetSessionIdName())
		t.Assert(e1, nil)
		t.AssertNE(sessionId, nil)
		t.Assert(response.ReadAllString(), "key=100")
	})
	time.Sleep(time.Second)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		client.SetHeader(s.GetSessionIdName(), sessionId)
		t.Assert(client.GetContent(ctx, "/get?k=key"), "100")
		t.Assert(client.GetContent(ctx, "/get?k=key1"), "")
	})
}

func Test_Session_Custom_Id(t *testing.T) {
	var (
		sessionId = "1234567890"
		key       = "key"
		value     = "value"
		s         = g.Server(guid.S())
	)
	s.BindHandler("/id", func(r *ghttp.Request) {
		if err := r.Session.SetId(sessionId); err != nil {
			r.Response.WriteExit(err.Error())
		}
		if err := r.Session.Set(key, value); err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(r.Session.Id())
	})
	s.BindHandler("/value", func(r *ghttp.Request) {
		r.Response.WriteExit(r.Session.Get(key))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		r, err := client.Get(ctx, "/id")
		t.AssertNil(err)
		defer r.Close()
		t.Assert(r.ReadAllString(), sessionId)
		t.Assert(r.GetCookie(s.GetSessionIdName()), sessionId)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		client.SetHeader(s.GetSessionIdName(), sessionId)
		t.Assert(client.GetContent(ctx, "/value"), value)
	})
}

func Test_Session_New_Id(t *testing.T) {
	var (
		sessionId     = "1234567890"
		newSessionId  = "0987654321"
		newSessionId2 = "abcdefghij"
		key           = "key"
		value         = "value"
		s             = g.Server(guid.S())
	)
	s.BindHandler("/id", func(r *ghttp.Request) {
		if err := r.Session.SetId(sessionId); err != nil {
			r.Response.WriteExit(err.Error())
		}
		if err := r.Session.Set(key, value); err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(r.Session.Id())
	})

	s.BindHandler("/newIdBySession", func(r *ghttp.Request) {
		// Use before session init
		if err := r.Session.SetId(newSessionId); err != nil {
			r.Response.WriteExit(err.Error())
		}
		if err := r.Session.Set(key, value); err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(r.Session.Id())
	})

	s.BindHandler("/newIdByCookie", func(r *ghttp.Request) {
		if err := r.Session.Remove("someKey"); err != nil {
			r.Response.WriteExit(err.Error())
		}

		r.Cookie.SetSessionId(newSessionId2)
		//r.Response.WriteExit(r.Session.Id())    // only change in cookie

		r.Response.WriteExit(newSessionId2)
	})

	s.BindHandler("/value", func(r *ghttp.Request) {
		r.Response.WriteExit(r.Session.Get(key))
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		r, err := client.Get(ctx, "/id")
		t.AssertNil(err)
		defer r.Close()
		t.Assert(r.ReadAllString(), sessionId)
		t.Assert(r.GetCookie(s.GetSessionIdName()), sessionId)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		client.SetHeader(s.GetSessionIdName(), sessionId)
		t.Assert(client.GetContent(ctx, "/value"), value)
	})

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		client.SetHeader(s.GetSessionIdName(), sessionId)
		r, err := client.Get(ctx, "/newIdBySession")
		t.AssertNil(err)
		defer r.Close()
		t.Assert(r.ReadAllString(), newSessionId)
		t.Assert(r.GetCookie(s.GetSessionIdName()), newSessionId)
	})

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		r, err := client.Get(ctx, "/newIdByCookie")
		client.SetHeader(s.GetSessionIdName(), sessionId)
		t.AssertNil(err)
		defer r.Close()
		t.Assert(r.ReadAllString(), newSessionId2)
		t.Assert(r.GetCookie(s.GetSessionIdName()), newSessionId2)
	})
}

func Test_Session_SetSessionIdFunc(t *testing.T) {
	const (
		customSessionId      = "custom-session-id"
		nextSessionId        = "next-session-id"
		headerNewSessionId   = "X-New-Session-Id"
		headerEmptySessionId = "X-Empty-Session-Id"
		sessionValueKey      = "k"
		sessionValue         = "v"
	)
	s := g.Server(guid.S())
	s.SetSessionStorage(gsession.NewStorageMemory())
	s.SetSessionIdFunc(func(r *ghttp.Request) string {
		if id := r.GetHeader(headerNewSessionId); id != "" {
			return id
		}
		if r.GetHeader(headerEmptySessionId) != "" {
			return ""
		}
		return customSessionId
	})
	s.BindHandler("/id", func(r *ghttp.Request) {
		if err := r.Session.Set(sessionValueKey, sessionValue); err != nil {
			r.Response.WriteExit(err.Error())
		}
		id, err := r.Session.Id()
		if err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(id)
	})
	s.BindHandler("/get", func(r *ghttp.Request) {
		value, err := r.Session.Get(sessionValueKey)
		if err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(value.String())
	})
	s.BindHandler("/regen", func(r *ghttp.Request) {
		if err := r.Session.Set(sessionValueKey, sessionValue); err != nil {
			r.Response.WriteExit(err.Error())
		}
		id, err := r.Session.RegenerateId(true)
		if err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(id)
	})
	s.SetDumpRouterMap(false)
	defer func() {
		if err := s.Shutdown(); err != nil {
			t.Error(err)
		}
	}()
	if err := s.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)
	prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())

	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(prefix)
		response, err := client.Get(ctx, "/id")
		if response != nil {
			defer response.Close()
		}
		t.AssertNil(err)
		t.Assert(response.ReadAllString(), customSessionId)
		t.Assert(response.GetCookie(s.GetSessionIdName()), customSessionId)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(prefix)
		client.SetHeader(s.GetSessionIdName(), customSessionId)
		t.Assert(client.GetContent(ctx, "/get"), sessionValue)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(prefix)
		client.SetHeader(s.GetSessionIdName(), customSessionId)
		client.SetHeader(headerNewSessionId, nextSessionId)
		response, err := client.Get(ctx, "/regen")
		if response != nil {
			defer response.Close()
		}
		t.AssertNil(err)
		t.Assert(response.ReadAllString(), nextSessionId)
		t.Assert(response.GetCookie(s.GetSessionIdName()), nextSessionId)
	})
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(prefix)
		client.SetHeader(headerEmptySessionId, "1")
		response, err := client.Get(ctx, "/id")
		if response != nil {
			defer response.Close()
		}
		t.AssertNil(err)
		t.Assert(response.StatusCode, http.StatusInternalServerError)
	})

	s2 := g.Server(guid.S())
	s2.SetSessionStorage(gsession.NewStorageMemory())
	s2.SetSessionIdFunc(func(r *ghttp.Request) string {
		return customSessionId
	})
	s2.SetSessionIdFunc(nil)
	s2.BindHandler("/id", func(r *ghttp.Request) {
		if err := r.Session.Set(sessionValueKey, sessionValue); err != nil {
			r.Response.WriteExit(err.Error())
		}
		id, err := r.Session.Id()
		if err != nil {
			r.Response.WriteExit(err.Error())
		}
		r.Response.WriteExit(id)
	})
	s2.SetDumpRouterMap(false)
	defer func() {
		if err := s2.Shutdown(); err != nil {
			t.Error(err)
		}
	}()
	if err := s2.Start(); err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		client := g.Client()
		client.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s2.GetListenedPort()))
		content := client.GetContent(ctx, "/id")
		t.Assert(len(content), 32)
		t.AssertNE(content, customSessionId)
	})
}
