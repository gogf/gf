package main

import "github.com/gogf/gf/g/encoding/gparser"

func main() {
	xml := `<?xml version="1.0" encoding="GBK"?>

	<Output type="o">
	<itotalSize>0</itotalSize>
	<ipageSize>1</ipageSize>
	<ipageIndex>2</ipageIndex>
	<itotalRecords>3</itotalRecords>
	<nworkOrderDtos/>
	<nworkOrderFrontXML/>
	</Output>`
	p, err := gparser.LoadContent([]byte(xml))
	if err != nil {
		panic(err)
	}
	p.Dump()
}
