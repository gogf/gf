package mxj

import (
	"fmt"
	"testing"
)

var s = `"'<>&`

func TestEscapeChars(t *testing.T) {
	fmt.Println("\n================== TestEscapeChars")

	ss := escapeChars(s)

	if ss != `&quot;&apos;&lt;&gt;&amp;` {
		t.Fatal(s, ":", ss)
	}

	fmt.Println(" s:", s)
	fmt.Println("ss:", ss)
}

func TestXMLEscapeChars(t *testing.T) {
	fmt.Println("================== TestXMLEscapeChars")

	XMLEscapeChars(true)
	defer XMLEscapeChars(false)

	m := map[string]interface{}{"mychars": s}

	x, err := AnyXmlIndent(s, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("s:", string(x))

	x, err = AnyXmlIndent(m, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("m:", string(x))
}

func TestXMLSeqEscapeChars(t *testing.T) {
	fmt.Println("================== TestXMLSeqEscapeChars")
	data := []byte(`
		<doc>
			<shortDescription>&gt;0-2y</shortDescription>
		</doc>`)
	fmt.Println("data:", string(data))

	m, err := NewMapXmlSeq(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("m: %v\n", m)

	XMLEscapeChars(true)
	defer XMLEscapeChars(false)

	x, err := m.XmlSeqIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("m:", string(x))
}

func TestXMLSeqEscapeChars2(t *testing.T) {
	fmt.Println("================== TestXMLSeqEscapeChars2")
	data := []byte(`
		<doc>
			<shortDescription test="&amp;something here">&gt;0-2y</shortDescription>
			<shortDescription test="something there" quote="&quot;">&lt;10-15</shortDescription>
		</doc>`)
	fmt.Println("data:", string(data))

	m, err := NewMapXmlSeq(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("m: %v\n", m)

	XMLEscapeChars(true)
	defer XMLEscapeChars(false)

	x, err := m.XmlSeqIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("m:", string(x))
}
