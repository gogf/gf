package x2j

import (
	"fmt"
	"testing"
	"encoding/xml"
)

func TestUnmarshal(t *testing.T) {
	var doc = []byte(`<doc> <name>Mayer Hawthorne</name> <song> <title>A Long Time</title> <verse no="1"> <line no="1">Henry was a renegade</line> <line no="2">Didn't like to play it safe</line> </verse> </song> </doc>`)

	fmt.Println("\nUnmarshal test ... *map[string]interface{}, *string")
	m := make(map[string]interface{},0)
	err := Unmarshal(doc,&m)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("m:",m)

	var s string
	err = Unmarshal(doc,&s)
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("s:",s)
}

func TestStructValue(t *testing.T) {
	var doc = []byte(`<info><name>clbanning</name><address>unknown</address></info>`)
	type Info struct {
		XMLName xml.Name `xml:"info"`
		Name string      `xml:"name"`
		Address string   `xml:"address"`
	}

	var myInfo Info

	fmt.Println("\nUnmarshal test ... struct:",string(doc))
	err := Unmarshal(doc,&myInfo)
	if err != nil {
		fmt.Println("err:",err.Error())
	} else {
		fmt.Printf("myInfo: %+v\n",myInfo)
	}
}

func TestMapValue(t *testing.T) {
	var doc = `<doc> <name>Mayer Hawthorne</name> <song> <title>A Long Time</title> <verse no="1"> <line no="1">Henry was a renegade</line> <line no="2">Didn't like to play it safe</line> </verse> </song> </doc>`

	fmt.Println("\nTestMapValue of doc.song.verse w/ len(attrs) == 0.")
	fmt.Println("doc:",doc)
	v,err := DocValue(doc,"doc.song.verse")
	if err != nil {
		fmt.Println("err:",err.Error())
	}
	fmt.Println("result:",v)
}
