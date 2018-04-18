package x2j

import (
	"bytes"
	"fmt"
	"testing"
)

var doc = `<entry><vars><foo>bar</foo><foo2><hello>world</hello></foo2></vars></entry>`


func TestToMap(t *testing.T) {
	fmt.Println("\nToMap - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	m,err := ToMap(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println(WriteMap(m))
}

func TestToJson(t *testing.T) {
	fmt.Println("\nToJson - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	s,err := ToJson(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("json:",s)
}

func TestToJsonIndent(t *testing.T) {
	fmt.Println("\nToJsonIndent - Read doc:",doc)
	rdr := bytes.NewBufferString(doc)
	s,err := ToJsonIndent(rdr)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("json:",s)
}

func TestBulkParser(t *testing.T) {
	s := doc + `<this><is>an</err>`+ doc
	fmt.Println("\nBulkParser (with error) - Read doc:",s)
	rdr := bytes.NewBufferString(s)
	err := XmlMsgsFromReader(rdr,phandler,ehandler)
	if err != nil {
		fmt.Println("reader terminated:",err.Error())
	}
}

func phandler(m map[string]interface{}) bool {
	fmt.Println("phandler m:",m)
	return true
}

func ehandler(err error) bool {
	fmt.Println("ehandler err:",err.Error())
	return true
}

func TestBulkParserToJson(t *testing.T) {
	s := doc + `<this><is>an</err>`+ doc
	fmt.Println("\nBulkParser (with error) - Read doc:",s)
	rdr := bytes.NewBufferString(s)
	err := XmlMsgsFromReaderAsJson(rdr,phandlerj,ehandler)
	if err != nil {
		fmt.Println("reader terminated:",err.Error())
	}
}

func phandlerj(s string) bool {
	fmt.Println("phandlerj s:",s)
	return true
}

