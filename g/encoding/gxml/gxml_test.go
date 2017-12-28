package gxml_test

import (
    "testing"
    "fmt"

    "encoding/xml"
)

var content = `<?xml version="1.0" encoding="utf-8"?>
<servers version="1">
    <server>
        <serverName>Shanghai_VPN</serverName>
        <serverIP>127.0.0.1</serverIP>
    </server>
    <server>
        <serverName>Beijing_VPN</serverName>
        <serverIP>127.0.0.2</serverIP>
    </server>
</servers>`
func Test_Xml(t *testing.T) {
    //xml, err := gxml.Decode(bytes.TrimSpace([]byte(content)))
    //if err != nil {
    //    glog.Error(err)
    //}

    //v := make(map[string]interface{})
    v := make([]interface{}, 0)
    e := xml.Unmarshal([]byte(content), &v)
    fmt.Println(e)
    fmt.Println(v)
}
