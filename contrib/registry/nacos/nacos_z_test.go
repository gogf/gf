package nacos

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	registry := New([]string{"127.0.0.1:8848"},
		//you can create a namespace in nacos, then copy the namespaceId to here
		//TODO when testing,please replace the namespace id with yours
		WithNameSpaceId("bc32382f-ada0-484b-87bb-3d5bfa0bc3ce"),
		WithContextPath("/nacos"))
	ctx := context.Background()
	svc := &gsvc.LocalService{
		Name:      "goframe-provider-0-tcp",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:8080,127.0.0.1:8088"),
	}
	s, err := registry.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s.GetEndpoints())
	t.Log(s.GetName())
	t.Log(s.GetVersion())
	t.Log(s.GetMetadata())
	t.Log(s.GetKey())
	t.Log(s.GetValue())
	// if success, you can see a name of goframe-provider-0-tcp in nacos web page
	err = registry.Deregister(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSearch(t *testing.T) {
	registry := New([]string{"127.0.0.1:8848"})
	ctx := context.Background()
	svc := &gsvc.LocalService{
		Name:      "goframe-provider-0-tcp",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:8080"),
	}
	s, err := registry.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s.GetEndpoints())
	t.Log(s.GetName())
	t.Log(s.GetVersion())
	t.Log(s.GetMetadata())
	t.Log(s.GetKey())
	t.Log(s.GetValue())
	time.Sleep(time.Second * 1)
	result, err := registry.Search(ctx, gsvc.SearchInput{
		Prefix: s.GetPrefix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(result))
	t.Log(s.GetEndpoints())
	t.Log(s.GetName())
	t.Log(s.GetVersion())
	t.Log(s.GetMetadata())
	t.Log(s.GetKey())
	t.Log(s.GetValue())
}

func TestWatch(t *testing.T) {
	r := New([]string{"127.0.0.1:8848"})

	ctx := gctx.New()

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}
	t.Log("watch")
	watch, err := r.Watch(ctx, svc.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	// watch svc
	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		g.Log().Info(ctx, "Register Proceed service: ", instance)
	}

	err = r.Deregister(ctx, s1)
	if err != nil {
		t.Log(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		g.Log().Info(ctx, "Deregister Proceed service: ", instance)
	}

	err = watch.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = watch.Proceed()
	if err == nil {
		// if nil, stop failed
		t.Fatal()
	}
}
