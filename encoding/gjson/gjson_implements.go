// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (j Json) MarshalJSON() ([]byte, error) {
	return j.ToJson()
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func (j *Json) UnmarshalJSON(b []byte) error {
	r, err := loadContentWithOptions(b, Options{
		Type:      ContentTypeJSON,
		StrNumber: true,
	})
	if r != nil {
		// Value copy.
		*j = *r
	}
	return err
}

// UnmarshalValue is an interface implement which sets any type of value for Json.
func (j *Json) UnmarshalValue(value interface{}) error {
	if r := NewWithOptions(value, Options{
		StrNumber: true,
	}); r != nil {
		// Value copy.
		*j = *r
	}
	return nil
}

// MapStrAny implements interface function MapStrAny().
func (j *Json) MapStrAny() map[string]interface{} {
	if j == nil {
		return nil
	}
	return j.Map()
}

// Interfaces implements interface function Interfaces().
func (j *Json) Interfaces() []interface{} {
	if j == nil {
		return nil
	}
	return j.Array()
}

// String returns current Json object as string.
func (j *Json) String() string {
	if j.IsNil() {
		return ""
	}
	return j.MustToJsonString()
}
