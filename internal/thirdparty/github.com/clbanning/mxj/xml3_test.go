// xml3_test.go - patch tests

package mxj

import (
	"fmt"
	"testing"
)

func TestXml3(t *testing.T) {
	fmt.Println("\n------------ xml3_test.go")
}

// for: https://github.com/clbanning/mxj/pull/26
func TestOnlyAttributes(t *testing.T) {
	fmt.Println("========== TestOnlyAttributes")
	dom, err := NewMapXml([]byte(`
		<memballoon model="virtio">
			<address type="pci" domain="0x0000" bus="0x00" slot="0x05" function="0x0"/>
			<empty/>
		</memballoon>)`))
	if err != nil {
		t.Fatal(err)
	}
	xml, err := dom.XmlIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(xml))
}

func TestOnlyAttributesSeq(t *testing.T) {
	fmt.Println("========== TestOnlyAttributesSeq")
	dom, err := NewMapXmlSeq([]byte(`
		<memballoon model="virtio">
			<address type="pci" domain="0x0000" bus="0x00" slot="0x05" function="0x0"/>
			<empty/>
		</memballoon>)`))
	if err != nil {
		t.Fatal(err)
	}
	xml, err := dom.XmlSeqIndent("", "  ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(xml))
}
