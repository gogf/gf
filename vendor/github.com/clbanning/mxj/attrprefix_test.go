// attrprefix_test.go - change attrPrefix var

package mxj

import (
	"fmt"
	"testing"
)

var data = []byte(`
<doc>
	<elem1 attr1="this" attr2="is">a test</elem1>
	<elem2 attr1="this" attr2="is not">a test</elem2>
</doc>
`)

func TestPrefixDefault(t *testing.T) {
	fmt.Println("----------------- TestPrefixDefault ...")
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	vals, err := m.ValuesForKey("-attr1")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 -attr1 vals", len(vals))
	}
	vals, err = m.ValuesForKey("-attr2")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 -attr2 vals", len(vals))
	}
}

func TestPrefixNoHyphen(t *testing.T) {
	fmt.Println("----------------- TestPrefixNoHyphen ...")
	PrependAttrWithHyphen(false)
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	vals, err := m.ValuesForKey("attr1")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 attr1 vals", len(vals))
	}
	vals, err = m.ValuesForKey("attr2")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 attr2 vals", len(vals))
	}
}

func TestPrefixUnderscore(t *testing.T) {
	fmt.Println("----------------- TestPrefixUnderscore ...")
	SetAttrPrefix("_")
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	vals, err := m.ValuesForKey("_attr1")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 _attr1 vals", len(vals))
	}
	vals, err = m.ValuesForKey("_attr2")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 _attr2 vals", len(vals))
	}
}

func TestPrefixAt(t *testing.T) {
	fmt.Println("----------------- TestPrefixAt ...")
	SetAttrPrefix("@")
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	vals, err := m.ValuesForKey("@attr1")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 @attr1 vals", len(vals))
	}
	vals, err = m.ValuesForKey("@attr2")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 2 {
		t.Fatal("didn't get 2 @attr2 vals", len(vals))
	}
}

func TestMarshalPrefixDefault(t *testing.T) {
	fmt.Println("----------------- TestMarshalPrefixDefault ...")
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	x, err := m.XmlIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(x))
}

func TestMarshalPrefixNoHyphen(t *testing.T) {
	fmt.Println("----------------- TestMarshalPrefixNoHyphen ...")
	PrependAttrWithHyphen(false)
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	_, err = m.XmlIndent("", "  ")
	if err == nil {
		t.Fatal("error not reported for invalid key label")
	}
	fmt.Println("err ok:", err)
}

func TestMarshalPrefixUnderscore(t *testing.T) {
	fmt.Println("----------------- TestMarshalPrefixUnderscore ...")
	SetAttrPrefix("_")
	m, err := NewMapXml(data)
	if err != nil {
		t.Fatal(err)
	}
	x, err := m.XmlIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(x))
}

