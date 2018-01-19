package mxj

import (
	"fmt"
	"testing"
)

func TestSetSubkeyFieldSeparator(t *testing.T) {
	PrependAttrWithHyphen(true)

	fmt.Println("----------- TestSetSubkeyFieldSeparator")
	data := `
		<doc>
			<elem attr="1">value 1</elem>
			<elem attr="2">value 2</elem>
			<elem attr="3">value 3</elem>
		</doc>`

	m, err := NewMapXml([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	vals, err := m.ValuesForKey("elem", "-attr:2:text")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 1 {
		t.Fatal(":len(vals);", len(vals), vals)
	}
	if vals[0].(map[string]interface{})["#text"].(string) != "value 2" {
		t.Fatal(":expecting: value 2; got:", vals[0].(map[string]interface{})["#text"])
	}

	SetFieldSeparator("|")
	defer SetFieldSeparator()
	vals, err = m.ValuesForKey("elem", "-attr|2|text")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 1 {
		t.Fatal("|len(vals);", len(vals), vals)
	}
	if vals[0].(map[string]interface{})["#text"].(string) != "value 2" {
		t.Fatal("|expecting: value 2; got:", vals[0].(map[string]interface{})["#text"])
	}
}

