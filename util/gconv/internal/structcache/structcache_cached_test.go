package structcache

import (
	"reflect"
	"testing"
)

func Test_parseStructToCachedStructInfo(t *testing.T) {
	// can't use gtest.C, because it's import cycle
	csi := &CachedStructInfo{
		tagOrFiledNameToFieldInfoMap: make(map[string]*CachedFieldInfo),
	}
	type UserInner struct {
		Inner int
	}

	type User struct {
		Username1 string  `orm:"user_name"`
		Username2 *string `orm:"user_name"`
		OtherInfo string  `orm:"other_info"`
		Inner     int

		*UserInner
	}

	var u User
	userType := reflect.TypeOf(u)
	parseStructToCachedStructInfo(userType, []int{}, csi, []string{"orm"})
	// OtherInfo\Username1\Username2\user_nam\other_info\Inner
	if len(csi.tagOrFiledNameToFieldInfoMap) != 6 {
		t.Errorf("tagOrFiledNameToFieldInfoMap failed, got: %v, want: 5", len(csi.tagOrFiledNameToFieldInfoMap))
	}

	t.Run("user_name must has OtherSameNameField", func(t *testing.T) {
		fieldInfo := csi.GetFieldInfo("user_name")
		if len(fieldInfo.OtherSameNameField) != 1 {
			t.Errorf("OtherSameNameField failed, got: %v, want: 1", len(fieldInfo.OtherSameNameField))
		}
		if fieldInfo.StructField.Name != "Username1" {
			t.Errorf("StructField.Name failed, got: %v, want: Username1", fieldInfo.StructField.Name)
		}
		if fieldInfo.OtherSameNameField[0].StructField.Name != "Username2" {
			t.Errorf("OtherSameNameField[0].StructField.Name failed, got: %v, want: Username2", fieldInfo.OtherSameNameField[0].StructField.Name)
		}
		if fieldInfo.StructField.Index[0] != 0 {
			t.Errorf("StructField.Index[0] failed, got: %v, want: 0", fieldInfo.StructField.Index[0])
		}
		if fieldInfo.OtherSameNameField[0].FieldIndexes[0] != 1 {
			t.Errorf("OtherSameNameField[0].FieldIndexes[0] failed, got: %v, want: 1", fieldInfo.OtherSameNameField[0].FieldIndexes[0])
		}
		if fieldInfo.StructField.Type.String() != reflect.String.String() {
			t.Errorf("StructField.Type failed, got: %v, want: string", fieldInfo.StructField.Type.String())
		}
		if fieldInfo.OtherSameNameField[0].StructField.Type.String() != "*string" {
			t.Errorf("OtherSameNameField[0].StructField.Type failed, got: %v, want: *string", fieldInfo.OtherSameNameField[0].StructField.Type.String())
		}
	})

	t.Run("Username1 must has OtherSameNameField", func(t *testing.T) {
		fieldInfo := csi.GetFieldInfo("Username1")
		if len(fieldInfo.OtherSameNameField) != 1 {
			t.Errorf("OtherSameNameField failed, got: %v, want: 1", len(fieldInfo.OtherSameNameField))
		}
		if fieldInfo.OtherSameNameField[0].StructField.Name != "Username2" {
			t.Errorf("OtherSameNameField[0].StructField.Name failed, got: %v, want: Username2", fieldInfo.OtherSameNameField[0].StructField.Name)
		}
		if fieldInfo.StructField.Name != "Username1" {
			t.Errorf("StructField.Name failed, got: %v, want: Username1", fieldInfo.StructField.Name)
		}
		if fieldInfo.StructField.Type.String() != reflect.String.String() {
			t.Errorf("StructField.Type failed, got: %v, want: string", fieldInfo.StructField.Type.String())
		}
		if fieldInfo.OtherSameNameField[0].StructField.Type.String() != "*string" {
			t.Errorf("OtherSameNameField[0].StructField.Type failed, got: %v, want: *string", fieldInfo.OtherSameNameField[0].StructField.Type.String())
		}

	})

	t.Run("Username2 must has OtherSameNameField", func(t *testing.T) {
		// Username2 has no OtherSameNameField
		fieldInfo := csi.GetFieldInfo("Username2")
		if len(fieldInfo.OtherSameNameField) != 0 {
			t.Errorf("OtherSameNameField failed, got: %v, want: 0", len(fieldInfo.OtherSameNameField))
		}
		if fieldInfo.StructField.Name != "Username2" {
			t.Errorf("StructField.Name failed, got: %v, want: Username2", fieldInfo.StructField.Name)
		}
	})

	t.Run("OtherInfo must has OtherSameNameField", func(t *testing.T) {
		// other_info has no OtherSameNameField
		fieldInfo := csi.GetFieldInfo("other_info")
		if len(fieldInfo.OtherSameNameField) != 0 {
			t.Errorf("OtherSameNameField failed, got: %v, want: 0", len(fieldInfo.OtherSameNameField))
		}
		if fieldInfo.StructField.Name != "OtherInfo" {
			t.Errorf("StructField.Name failed, got: %v, want: OtherInfo", fieldInfo.StructField.Name)
		}
		fieldInfo = csi.GetFieldInfo("OtherInfo")
		if len(fieldInfo.OtherSameNameField) != 0 {
			t.Errorf("OtherSameNameField failed, got: %v, want: 0", len(fieldInfo.OtherSameNameField))
		}
		if fieldInfo.StructField.Name != "OtherInfo" {
			t.Errorf("StructField.Name failed, got: %v, want: OtherInfo", fieldInfo.StructField.Name)
		}
	})

	t.Run("Inner must has OtherSameNameField", func(t *testing.T) {
		// Inner has no OtherSameNameField
		fieldInfo := csi.GetFieldInfo("Inner")
		if len(fieldInfo.OtherSameNameField) != 1 {
			t.Errorf("OtherSameNameField failed, got: %v, want: 1", len(fieldInfo.OtherSameNameField))
		}
		if fieldInfo.StructField.Name != "Inner" {
			t.Errorf("StructField.Name failed, got: %v, want: Inner", fieldInfo.StructField.Name)
		}
		if fieldInfo.StructField.Index[0] != 3 {
			t.Errorf("StructField.Index[0] failed, got: %v, want: 3", fieldInfo.StructField.Index[0])
		}
		if fieldInfo.OtherSameNameField[0].FieldIndexes[0] != 4 {
			t.Errorf("FieldIndexes[0] failed, got: %v, want: 4", fieldInfo.FieldIndexes[0])
		}
	})
}
