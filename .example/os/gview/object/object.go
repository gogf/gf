package main

import (
	"github.com/gogf/gf/v2/frame/g"
)

type T struct {
	Name string
}

func (t *T) Hello(name string) string {
	return "Hello " + name
}

func (t *T) Test() string {
	return "This is test"
}

func main() {
	t := &T{"John"}
	v := g.View()
	content := `{{.t.Hello "there"}}, my name's {{.t.Name}}. {{.t.Test}}.`
	if r, err := v.ParseContent(content, g.Map{"t": t}); err != nil {
		g.Dump(err)
	} else {
		g.Dump(r)
	}
}
