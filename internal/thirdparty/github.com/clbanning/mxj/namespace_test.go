package mxj

import (
	"fmt"
	"testing"
)

func TestNamespaceHeader(t *testing.T) {
	fmt.Println("\n---------------- namespace_test.go ...")
}

func TestBeautifyXml(t *testing.T) {
	fmt.Println("\n----------------  TestBeautifyXml ...")
	const flatxml = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ns="http://example.com/ns"><soapenv:Header/><soapenv:Body><ns:request><ns:customer><ns:id>123</ns:id><ns:name type="NCHZ">John Brown</ns:name></ns:customer></ns:request></soapenv:Body></soapenv:Envelope>`
	v, err := BeautifyXml([]byte(flatxml), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(flatxml)
	fmt.Println(string(v))
}
