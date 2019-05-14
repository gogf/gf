// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

//func MarshalOrdered(value interface{}) ([]byte, error) {
//	buffer := bytes.NewBuffer(nil)
//	rv     := reflect.ValueOf(value)
//	kind   := rv.Kind()
//	if kind == reflect.Ptr {
//		rv   = rv.Elem()
//		kind = rv.Kind()
//	}
//	switch kind {
//		case reflect.Slice: fallthrough
//		case reflect.Array:
//			buffer.WriteByte('[')
//			length := rv.Len()
//			for i := 0; i < length; i++ {
//				if p, err := MarshalOrdered(rv.Index(i).Interface()); err != nil {
//					return nil, err
//				} else {
//					buffer.Write(p)
//					if i < length - 1 {
//						buffer.WriteByte(',')
//					}
//				}
//			}
//			buffer.WriteByte(']')
//		case reflect.Map: fallthrough
//		case reflect.Struct:
//			m     := gconv.Map(value, "json")
//			keys  := make([]string, len(m))
//			index := 0
//			for key := range m {
//				keys[index] = key
//				index++
//			}
//			sort.Strings(keys)
//			buffer.WriteByte('{')
//			for i, key := range keys {
//				if p, err := MarshalOrdered(m[key]); err != nil {
//					return nil, err
//				} else {
//					buffer.WriteString(fmt.Sprintf(`"%s":%s`, key, string(p)))
//					if i < index - 1 {
//						buffer.WriteByte(',')
//					}
//				}
//			}
//			buffer.WriteByte('}')
//		default:
//			if p, err := json.Marshal(value); err != nil {
//				return nil, err
//			} else {
//				buffer.Write(p)
//			}
//	}
//	return buffer.Bytes(), nil
//}