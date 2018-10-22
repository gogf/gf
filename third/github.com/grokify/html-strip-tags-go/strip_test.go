package strip

import (
	"testing"
)

func TestStripTags(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"", ""},
		{"Hello, World!", "Hello, World!"},
		{"foo&amp;bar", "foo&amp;bar"},
		{`Hello <a href="www.example.com/">World</a>!`, "Hello World!"},
		{"Foo <textarea>Bar</textarea> Baz", "Foo Bar Baz"},
		{"Foo <!-- Bar --> Baz", "Foo  Baz"},
		{"<", "<"},
		{"foo < bar", "foo < bar"},
		{`Foo<script type="text/javascript">alert(1337)</script>Bar`, "FooBar"},
		{`Foo<div title="1>2">Bar`, "FooBar"},
		{`I <3 Ponies!`, `I <3 Ponies!`},
		{`<script>foo()</script>`, ``},
	}

	for _, test := range tests {
		if got := StripTags(test.input); got != test.want {
			t.Errorf("%q: want %q, got %q", test.input, test.want, got)
		}
	}
}
