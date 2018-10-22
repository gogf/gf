package x2j

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGoofy(t *testing.T) {
	var doc = `<xml><tag one="1" pi="3.1415962535" bool="true"/><tagJR key="value"/></xml>`
	type goofy struct {
		S string
		Sp *string
	}
	g := new(goofy)
	g.S = "Now is the time for"
	tmp := "all good men to come to"
	g.Sp = &tmp

	m, err := DocToMap(doc)
	if err != nil {
		fmt.Println("err:",err.Error())
		return
	}

	m["goofyVal"] = interface{}(g)
	m["byteVal"] = interface{}([]byte(`the aid of their country`))
	m["nilVal"] = interface{}(nil)

	fmt.Println("\nTestGoofy ... MapToDoc:",m)
	var v []byte
	v,err = json.Marshal(m)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))

	type goofier struct {
		G *goofy
		B []byte
		N *string
	}
	gg := new(goofier)
	gg.G = g
	gg.B = []byte(`the tree of freedom must periodically be`)
	gg.N = nil
	m["goofierVal"] = interface{}(gg)

	fmt.Println("\nTestGoofier ... MapToDoc:",m)
	v,err = json.Marshal(m)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("v:",string(v))
}

