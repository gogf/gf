package gxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"reflect"

	"github.com/gogf/gf/v2/util/gconv"
)

type gxmlMap map[string]interface{}

func (m *gxmlMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = gxmlMap{}
	var a map[string]interface{}
	for {
		// var e xmlMapEntry
		err := d.Decode(&a)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	fmt.Println(a)
	return nil
}

func (m gxmlMap) MarshalXML(enc *xml.Encoder, _ xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	if err := marshal(enc, m); err != nil {
		return err
	}

	return nil
}

func Decode2(content []byte) (map[string]interface{}, error) {
	var enc = xml.NewDecoder(bytes.NewReader(content))

	var rm map[string]interface{}
	if err := enc.Decode((*gxmlMap)(&rm)); err != nil {
		return nil, err
	}
	return nil, nil
}

func Encode(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	var (
		b    bytes.Buffer
		enc  = xml.NewEncoder(&b)
		data = mergeRootTag(m, rootTag...)
	)

	if err := enc.Encode(gxmlMap(data)); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func EncodeWithIndent(m map[string]interface{}, rootTag ...string) ([]byte, error) {
	var (
		b    bytes.Buffer
		enc  = xml.NewEncoder(&b)
		data = mergeRootTag(m, rootTag...)
	)
	enc.Indent("", "\t")

	if err := enc.Encode(gxmlMap(data)); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func mergeRootTag(m map[string]interface{}, rootTag ...string) map[string]interface{} {
	if len(rootTag) == 0 && len(m) > 1 {
		rootTag = []string{"doc"}
	}
	for i := len(rootTag) - 1; i >= 0; i-- {
		m = map[string]interface{}{rootTag[i]: m}
	}
	return m
}

func marshal(enc *xml.Encoder, m map[string]interface{}) error {
	for key, value := range m {
		if reflect.TypeOf(value).Kind() == reflect.Map {
			v, t := value.(map[string]interface{})
			if !t {
				v = gconv.Map(value)
			}
			start := xml.StartElement{Name: xml.Name{Local: key}}

			if err := enc.EncodeToken(start); err != nil {
				return err
			}

			if err := marshal(enc, v); err != nil {
				return err
			}

			if err := enc.EncodeToken(start.End()); err != nil {
				return err
			}
			return nil
		} else {
			err := enc.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: key}})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
