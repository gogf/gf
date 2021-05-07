// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghtml_test

import (
	"github.com/gogf/gf/frame/g"
	"testing"

	"github.com/gogf/gf/encoding/ghtml"
	"github.com/gogf/gf/test/gtest"
)

func Test_StripTags(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := `<p>Test paragraph.</p><!-- Comment -->  <a href="#fragment">Other text</a>`
		dst := `Test paragraph.  Other text`
		t.Assert(ghtml.StripTags(src), dst)
	})
}

func Test_Entities(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := `A 'quote' "is" <b>bold</b>`
		dst := `A &#39;quote&#39; &#34;is&#34; &lt;b&gt;bold&lt;/b&gt;`
		t.Assert(ghtml.Entities(src), dst)
		t.Assert(ghtml.EntitiesDecode(dst), src)
	})
}

func Test_SpecialChars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := `A 'quote' "is" <b>bold</b>`
		dst := `A &#39;quote&#39; &#34;is&#34; &lt;b&gt;bold&lt;/b&gt;`
		t.Assert(ghtml.SpecialChars(src), dst)
		t.Assert(ghtml.SpecialCharsDecode(dst), src)
	})
}

func Test_SpecialCharsMapOrStruct_Map(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := g.Map{
			"Title":   "<h1>T</h1>",
			"Content": "<div>C</div>",
		}
		err := ghtml.SpecialCharsMapOrStruct(a)
		t.Assert(err, nil)
		t.Assert(a["Title"], `&lt;h1&gt;T&lt;/h1&gt;`)
		t.Assert(a["Content"], `&lt;div&gt;C&lt;/div&gt;`)
	})
	gtest.C(t, func(t *gtest.T) {
		a := g.MapStrStr{
			"Title":   "<h1>T</h1>",
			"Content": "<div>C</div>",
		}
		err := ghtml.SpecialCharsMapOrStruct(a)
		t.Assert(err, nil)
		t.Assert(a["Title"], `&lt;h1&gt;T&lt;/h1&gt;`)
		t.Assert(a["Content"], `&lt;div&gt;C&lt;/div&gt;`)
	})
}

func Test_SpecialCharsMapOrStruct_Struct(t *testing.T) {
	type A struct {
		Title   string
		Content string
	}
	gtest.C(t, func(t *gtest.T) {
		a := &A{
			Title:   "<h1>T</h1>",
			Content: "<div>C</div>",
		}
		err := ghtml.SpecialCharsMapOrStruct(a)
		t.Assert(err, nil)
		t.Assert(a.Title, `&lt;h1&gt;T&lt;/h1&gt;`)
		t.Assert(a.Content, `&lt;div&gt;C&lt;/div&gt;`)
	})
}
