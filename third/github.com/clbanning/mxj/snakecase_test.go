package mxj

import (
	"fmt"
	"testing"
)

func TestStakeCase(t *testing.T) {
	PrependAttrWithHyphen(true)
	fmt.Println("\n----------- TestSnakeCase")
	CoerceKeysToSnakeCase()
	defer CoerceKeysToSnakeCase()

	data1 := `<xml-rpc><element-one attr-1="an attribute">something</element-one></xml-rpc>`
	data2 := `<xml_rpc><element_one attr_1="an attribute">something</element_one></xml_rpc>`

	m, err := NewMapXml([]byte(data1))
	if err != nil {
		t.Fatal(err)
	}

	x, err := m.Xml()
	if err != nil {
		t.Fatal(err)
	}
	if string(x) != data2 {
		t.Fatal(string(x), "!=", data2)
	}

	// Use-case from: https://github.com/clbanning/mxj/pull/33#issuecomment-273724506
	data1 = `<rpc-reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" xmlns:junos="http://xml.juniper.net/junos/11.2R4/junos" message-id="97741fa3-99e8-46ba-b103-bab6b459d884">
<software-information>
<host-name>srx100</host-name>
<product-model>srx100b</product-model>
<product-name>srx100b</product-name>
<jsr/>
<package-information>
<name>junos</name>
<comment>JUNOS Software Release [11.2R4.3]</comment>
</package-information>
</software-information>
</rpc-reply>`
	data2 = `<rpc_reply xmlns="urn:ietf:params:xml:ns:netconf:base:1.0" xmlns:junos="http://xml.juniper.net/junos/11.2R4/junos" message_id="97741fa3-99e8-46ba-b103-bab6b459d884">
<software_information>
<host_name>srx100</host_name>
<product_model>srx100b</product_model>
<product_name>srx100b</product_name>
<jsr/>
<package_information>
<name>junos</name>
<comment>JUNOS Software Release [11.2R4.3]</comment>
</package_information>
</software_information>
</rpc_reply>`

	m, err = NewMapXmlSeq([]byte(data1))
	if err != nil {
		t.Fatal(err)
	}

	x, err = m.XmlSeqIndent("", "")
	if err != nil {
		t.Fatal(err)
	}
	if string(x) != data2 {
		t.Fatal(string(x), "!=", data2)
	}
}

