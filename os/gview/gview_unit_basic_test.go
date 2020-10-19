// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gview_test

import (
	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
)

func init() {
	os.Setenv("GF_GVIEW_ERRORPRINT", "false")
}

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		str := `hello {{.name}},version:{{.version}};hello {{GetName}},version:{{GetVersion}};{{.other}}`
		pwd := gfile.Pwd()
		view := gview.New()
		view.SetDelimiters("{{", "}}")
		view.AddPath(pwd)
		view.SetPath(pwd)
		view.Assign("name", "gf")
		view.Assigns(g.Map{"version": "1.7.0"})
		view.BindFunc("GetName", func() string { return "gf" })
		view.BindFuncMap(gview.FuncMap{"GetVersion": func() string { return "1.7.0" }})
		result, err := view.ParseContent(str, g.Map{"other": "that's all"})
		t.Assert(err != nil, false)
		t.Assert(result, "hello gf,version:1.7.0;hello gf,version:1.7.0;that's all")

		//测试api方法
		str = `hello {{.name}}`
		result, err = gview.ParseContent(str, g.Map{"name": "gf"})
		t.Assert(err != nil, false)
		t.Assert(result, "hello gf")

		//测试instance方法
		result, err = gview.Instance().ParseContent(str, g.Map{"name": "gf"})
		t.Assert(err != nil, false)
		t.Assert(result, "hello gf")
	})
}

func Test_Func(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		str := `{{eq 1 1}};{{eq 1 2}};{{eq "A" "B"}}`
		result, err := gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;false;false`)

		str = `{{ne 1 2}};{{ne 1 1}};{{ne "A" "B"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;false;true`)

		str = `{{lt 1 2}};{{lt 1 1}};{{lt 1 0}};{{lt "A" "B"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;false;false;true`)

		str = `{{le 1 2}};{{le 1 1}};{{le 1 0}};{{le "A" "B"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `true;true;false;true`)

		str = `{{gt 1 2}};{{gt 1 1}};{{gt 1 0}};{{gt "A" "B"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `false;false;true;false`)

		str = `{{ge 1 2}};{{ge 1 1}};{{ge 1 0}};{{ge "A" "B"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `false;true;true;false`)

		str = `{{"<div>测试</div>"|text}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `测试`)

		str = `{{"<div>测试</div>"|html}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `&lt;div&gt;测试&lt;/div&gt;`)

		str = `{{"<div>测试</div>"|htmlencode}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `&lt;div&gt;测试&lt;/div&gt;`)

		str = `{{"&lt;div&gt;测试&lt;/div&gt;"|htmldecode}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `<div>测试</div>`)

		str = `{{"https://goframe.org"|url}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `https%3A%2F%2Fgoframe.org`)

		str = `{{"https://goframe.org"|urlencode}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `https%3A%2F%2Fgoframe.org`)

		str = `{{"https%3A%2F%2Fgoframe.org"|urldecode}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `https://goframe.org`)
		str = `{{"https%3NA%2F%2Fgoframe.org"|urldecode}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(gstr.Contains(result, "invalid URL escape"), true)

		str = `{{1540822968 | date "Y-m-d"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `2018-10-29`)
		str = `{{date "Y-m-d"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)

		str = `{{"我是中国人" | substr 2 -1}};{{"我是中国人" | substr 2  2}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `中国人;中国`)

		str = `{{"我是中国人" | strlimit 2  "..."}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `我是...`)

		str = `{{"I'm中国人" | replace "I'm" "我是"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `我是中国人`)

		str = `{{compare "A" "B"}};{{compare "1" "2"}};{{compare 2 1}};{{compare 1 1}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `-1;-1;1;0`)

		str = `{{"热爱GF热爱生活" | hidestr 20  "*"}};{{"热爱GF热爱生活" | hidestr 50  "*"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `热爱GF*爱生活;热爱****生活`)

		str = `{{"热爱GF热爱生活" | highlight "GF" "red"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `热爱<span style="color:red;">GF</span>热爱生活`)

		str = `{{"gf" | toupper}};{{"GF" | tolower}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err != nil, false)
		t.Assert(result, `GF;gf`)

		str = `{{concat "I" "Love" "GoFrame"}}`
		result, err = gview.ParseContent(str, nil)
		t.Assert(err, nil)
		t.Assert(result, `ILoveGoFrame`)
	})
}

func Test_FuncNl2Br(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		str := `{{"Go\nFrame" | nl2br}}`
		result, err := gview.ParseContent(str, nil)
		t.Assert(err, nil)
		t.Assert(result, `Go<br>Frame`)
	})
	gtest.C(t, func(t *gtest.T) {
		s := ""
		for i := 0; i < 3000; i++ {
			s += "Go\nFrame\n中文"
		}
		str := `{{.content | nl2br}}`
		result, err := gview.ParseContent(str, g.Map{
			"content": s,
		})
		t.Assert(err, nil)
		t.Assert(result, strings.Replace(s, "\n", "<br>", -1))
	})
}

func Test_FuncInclude(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		header := `<h1>HEADER</h1>`
		main := `<h1>hello gf</h1>`
		footer := `<h1>FOOTER</h1>`
		layout := `{{include "header.html" .}}
{{include "main.html" .}}
{{include "footer.html" .}}`
		templatePath := gfile.Pwd() + gfile.Separator + "template"
		gfile.Mkdir(templatePath)
		defer gfile.Remove(templatePath)
		//headerFile, _ := gfile.Create(templatePath + gfile.Separator + "header.html")
		err := ioutil.WriteFile(templatePath+gfile.Separator+"header.html", []byte(header), 0644)
		if err != nil {
			t.Error(err)
		}
		ioutil.WriteFile(templatePath+gfile.Separator+"main.html", []byte(main), 0644)
		ioutil.WriteFile(templatePath+gfile.Separator+"footer.html", []byte(footer), 0644)
		ioutil.WriteFile(templatePath+gfile.Separator+"layout.html", []byte(layout), 0644)
		view := gview.New(templatePath)
		result, err := view.Parse("notfound.html")
		t.Assert(err != nil, true)
		t.Assert(result, ``)
		result, err = view.Parse("layout.html")
		t.Assert(err != nil, false)
		t.Assert(result, `<h1>HEADER</h1>
<h1>hello gf</h1>
<h1>FOOTER</h1>`)
		notfoundPath := templatePath + gfile.Separator + "template" + gfile.Separator + "notfound.html"
		gfile.Mkdir(templatePath + gfile.Separator + "template")
		gfile.Create(notfoundPath)
		ioutil.WriteFile(notfoundPath, []byte("notfound"), 0644)
		result, err = view.Parse("notfound.html")
		t.Assert(err != nil, true)
		t.Assert(result, ``)
	})
}

func Test_SetPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		view := gview.Instance("addpath")
		err := view.AddPath("tmp")
		t.AssertNE(err, nil)

		err = view.AddPath("gview.go")
		t.AssertNE(err, nil)

		os.Setenv("GF_GVIEW_PATH", "tmp")
		view = gview.Instance("setpath")
		err = view.SetPath("tmp")
		t.AssertNE(err, nil)

		err = view.SetPath("gview.go")
		t.AssertNE(err, nil)

		view = gview.New(gfile.Pwd())
		err = view.SetPath("tmp")
		t.AssertNE(err, nil)

		err = view.SetPath("gview.go")
		t.AssertNE(err, nil)

		os.Setenv("GF_GVIEW_PATH", "template")
		gfile.Mkdir(gfile.Pwd() + gfile.Separator + "template")
		view = gview.New()
	})
}

func Test_ParseContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		str := `{{.name}}`
		view := gview.New()
		result, err := view.ParseContent(str, g.Map{"name": func() {}})
		t.Assert(err != nil, true)
		t.Assert(result, ``)
	})
}

func Test_HotReload(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		dirPath := gfile.Join(
			gfile.TempDir(),
			"testdata",
			"template-"+gconv.String(gtime.TimestampNano()),
		)
		defer gfile.Remove(dirPath)
		filePath := gfile.Join(dirPath, "test.html")

		// Initialize data.
		err := gfile.PutContents(filePath, "test:{{.var}}")
		t.Assert(err, nil)

		view := gview.New(dirPath)

		time.Sleep(100 * time.Millisecond)
		result, err := view.Parse("test.html", g.Map{
			"var": "1",
		})
		t.Assert(err, nil)
		t.Assert(result, `test:1`)

		// Update data.
		err = gfile.PutContents(filePath, "test2:{{.var}}")
		t.Assert(err, nil)

		time.Sleep(100 * time.Millisecond)
		result, err = view.Parse("test.html", g.Map{
			"var": "2",
		})
		t.Assert(err, nil)
		t.Assert(result, `test2:2`)
	})
}

func Test_XSS(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		s := "<br>"
		r, err := v.ParseContent("{{.v}}", g.Map{
			"v": s,
		})
		t.Assert(err, nil)
		t.Assert(r, s)
	})
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetAutoEncode(true)
		s := "<br>"
		r, err := v.ParseContent("{{.v}}", g.Map{
			"v": s,
		})
		t.Assert(err, nil)
		t.Assert(r, ghtml.Entities(s))
	})
	// Tag "if".
	gtest.C(t, func(t *gtest.T) {
		v := gview.New()
		v.SetAutoEncode(true)
		s := "<br>"
		r, err := v.ParseContent("{{if eq 1 1}}{{.v}}{{end}}", g.Map{
			"v": s,
		})
		t.Assert(err, nil)
		t.Assert(r, ghtml.Entities(s))
	})
}
